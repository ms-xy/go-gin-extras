package session

import (
	"github.com/ms-xy/go-gin-extras/middlewares/common"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ContextKey = "github.com/ms-xy/go-gin-extras/middlewares/session"
)

var (
	// HandlePanic indicates whether or not the session middleware will handle
	// its own panics or not. Default is false. Setting it to true will result
	// in the same behavior as using gin.Use(Logging(), Recovery()) from this
	// repository.
	HandlePanic = false
)

func GetSessionMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		getGorillaSessionMiddleware(),
		getSessionMiddleware(CachePath),
	}
}

func getGorillaSessionMiddleware() gin.HandlerFunc {
	store := cookie.NewStore(SessionSecret, SessionAesKey)
	store.Options(sessions.Options{
		Domain:   SessionDomain,
		MaxAge:   SessionMaxAge,
		HttpOnly: SessionHttpOnly,
		Secure:   SessionSecure,
	})
	return sessions.Sessions("session", store)
}

func getSessionMiddleware(cachePath string) gin.HandlerFunc {
	cache, err := NewSessionCache(cachePath)
	if err != nil {
		panic(err) // should be handed up (?)
	}
	return func(c *gin.Context) {
		if HandlePanic {
			defer func() {
				start := time.Now()
				if r := recover(); r != nil {
					common.ResponseWriteError(c, r)
					common.WriteLogEntry(c, "SessionMiddleware", time.Since(start))
				}
			}()
		}
		s := NewSession(cache)
		s_sid := s.session.Get(SESSION_ID)
		if s_sid != nil {
			if s_id, ok := s_sid.(string); ok {
				s.SetID(uuid.MustParse(s_id))
			}
		}
		if err := s.Start(); err != nil {
			panic(err)
		}
		gsession := sessions.Default(c)
		gsession.Set(SESSION_ID, s.GetID().String())
		c.Set(ContextKey, s)
		gsession.Save()
	}
}

func GetSession(c *gin.Context) *Session {
	return c.MustGet(ContextKey).(*Session)
}
