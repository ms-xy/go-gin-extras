package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Of(c *gin.Context) *Session {
	return &Session{
		x: sessions.Default(c),
	}
}

type Session struct {
	x sessions.Session
}

func (this *Session) SessionID() uuid.UUID {
	return uuid.Must(uuid.Parse(this.x.Get(SESSION_ID).(string)))
	// return this.x.Get(SESSION_ID).(uuid.UUID)
}
