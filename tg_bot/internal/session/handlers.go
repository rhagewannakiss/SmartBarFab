package session

import (
	"fmt"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/a-palonskaa/SmartBar/tg_bot/internal/coctail"
	"github.com/a-palonskaa/SmartBar/tg_bot/internal/config"
)

var MessageConfig = map[string]string{
	"welcome":      "sosite-suki пиздатый проект вас приветствует🙌\nавторизуйтесь через /password [пароль]\nYours aliffka ❤️",
	"authorize":    "ах, котик, у тебя слишком большой... ввод...",
	"authOK":       "welcome to the club buddy",
	"authFail":     "unluck, my friend",
	"wrongcommand": "по клавишам научись попадать, наш дорогой пользователь💋",
	"commands":     "у вас следующие возможности:\n/make_drink\n/delay_drink [минуты]",
	"choose":       "доступные коктейли:\n" + coctail.BuildCoctailList(),
	"prepare":      "совсем скоро вы насладитесь необыкновенным ",
}

func (us *UserSessions) HandleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID
	text := msg.Text
	session := us.GetSession(chatID)

	switch session.State {
	case WAIT_PASSWORD:
		us.HandleWaitPassword(bot, chatID, text)
	case AUTHORIZED:
		us.HandleAuthorized(bot, chatID, text)
	case WAIT_DRINK:
		us.HandleWaitDrink(bot, chatID, text)
	}
}

func (us *UserSessions) HandleWaitPassword(bot *tgbotapi.BotAPI, chatID int64, text string) {
	if text == "/start" {
		bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["welcome"]))
		return
	}

	if len(text) > 10 && text[:10] == "/password " {
		password := text[10:]
		if password == config.BotPassword {
			us.GetSession(chatID).State = AUTHORIZED
			bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["authOK"]))
			bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["commands"]))
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["authFail"]))
		}
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["authorize"]))
}

func (us *UserSessions) HandleAuthorized(bot *tgbotapi.BotAPI, chatID int64, text string) {
	if text == "/make_drink" {
		us.GetSession(chatID).State = WAIT_DRINK
		bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["choose"]))
		return
	}

	if len(text) > 13 && text[:13] == "/delay_drink " {
		delayMin, _ := strconv.Atoi(text[13:])
		us.scheduleDrink(chatID, delayMin)
		us.GetSession(chatID).State = WAIT_DRINK
		bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["choose"]))
		return
	}

	bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["wrongcommand"]))
	bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["commands"]))
}

func (us *UserSessions) HandleWaitDrink(bot *tgbotapi.BotAPI, chatID int64, text string) {
	drink, ok := coctail.TgCoctailNamesToIR[text]
	if !ok {
		bot.Send(tgbotapi.NewMessage(chatID, MessageConfig["choose"]))
		return
	}

	session := us.GetSession(chatID)
	session.ScheduledDrink = drink
	session.State = DONE
	bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("%s '%s'", MessageConfig["prepare"], coctail.CoctailToNames[drink])))
}

func (us *UserSessions) scheduleDrink(chatID int64, delayMin int) {
	session := us.GetSession(chatID)
	if delayMin > 0 {
		session.ScheduledTime = time.Now().Add(time.Duration(delayMin) * time.Minute)
	} else {
		session.ScheduledTime = time.Now()
	}
}
