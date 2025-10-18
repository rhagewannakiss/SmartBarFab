package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/a-palonskaa/SmartBar/tg_bot/internal/config"
	esp "github.com/a-palonskaa/SmartBar/tg_bot/internal/esp-api"
	"github.com/a-palonskaa/SmartBar/tg_bot/internal/session"
	"github.com/a-palonskaa/SmartBar/tg_bot/pkg/logger"
)

func main() {
	logger.InitLogger("logs/info.log")
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create a new BotAPI instance")
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userSessions := session.NewUserSessions()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/next_drink", esp.NextDrinkHandler(userSessions))
	mux.HandleFunc("/api/status", esp.StatusHandler(userSessions))

	srv := &http.Server{
		Addr:    ":" + config.ServerPort,
		Handler: mux,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Info().Msg("API started on :" + config.ServerPort)
		if err := srv.ListenAndServe(); err != nil {
			serverErr <- err
		}

	}()

	updates := getUpdatesChan(bot)
	go func() {
		for {
			select {
			case <-ctx.Done():
				bot.StopReceivingUpdates()
				return
			case update := <-updates:
				if update.Message == nil {
					continue
				}
				userSessions.HandleMessage(bot, update.Message)
			case err := <-serverErr:
				log.Error().Err(err).Msg("server error")
			}
		}
	}()

	<-stop
	cancel()
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error().Err(err).Msg("error during server shutdown")
	}
}

func getUpdatesChan(bot *tgbotapi.BotAPI) tgbotapi.UpdatesChannel {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	return bot.GetUpdatesChan(updateConfig)
}
