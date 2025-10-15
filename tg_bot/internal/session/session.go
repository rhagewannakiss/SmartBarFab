package session

import (
	"sync"
	"time"

	"github.com/a-palonskaa/SmartBar/tg_bot/internal/coctail"
)

type SessionState int

const (
	WAIT_PASSWORD SessionState = iota
	AUTHORIZED
	WAIT_DRINK
	DONE
)

type UserSession struct {
	State          SessionState
	ScheduledTime  time.Time
	ScheduledDrink coctail.Coctail
}

type UserSessions struct {
	Sessions map[int64]*UserSession
	Mu       sync.Mutex
}

func NewUserSessions() *UserSessions {
	return &UserSessions{
		Sessions: make(map[int64]*UserSession, 0),
	}
}

func (us *UserSessions) GetSession(chatID int64) *UserSession {
	us.Mu.Lock()
	defer us.Mu.Unlock()

	if session, ok := us.Sessions[chatID]; ok {
		return session
	}

	s := &UserSession{State: WAIT_PASSWORD}
	us.Sessions[chatID] = s
	return s
}
