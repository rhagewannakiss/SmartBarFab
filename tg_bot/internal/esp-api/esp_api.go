package esp_api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/a-palonskaa/SmartBar/tg_bot/internal/coctail"
	"github.com/a-palonskaa/SmartBar/tg_bot/internal/session"
)

func NextDrinkHandler(userSessions *session.UserSessions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSessions.Mu.Lock()
		defer userSessions.Mu.Unlock()

		for id, userSession := range userSessions.Sessions {
			if userSession.State == session.DONE && time.Now().After(userSession.ScheduledTime) {
				drink := userSession.ScheduledDrink
				userSession.ScheduledDrink = coctail.UNDEFINED
				userSession.State = session.AUTHORIZED
				log.Info().Msgf("ESP got '%s' from %d", drink, id)
				w.Write([]byte(coctail.CoctailToESPNames[drink]))
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func StatusHandler(userSessions *session.UserSessions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userSessions.Mu.Lock()
		defer userSessions.Mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userSessions)
	}
}
