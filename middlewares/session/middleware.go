package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetSessionMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		getGorillaSessionMiddleware(),
		getAutoSessionMiddleware(CachePath),
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

func getAutoSessionMiddleware(cachePath string) gin.HandlerFunc {
	cache := NewSessionCache(cachePath)
	return func(c *gin.Context) {
		session := sessions.Default(c)
		var sessionID uuid.UUID = uuid.Nil
		if v := session.Get(SESSION_ID); v != nil {
			_sessionID := uuid.Must(uuid.Parse(v.(string)))
			if cache.Exists(_sessionID) { // not working
				sessionID = _sessionID
			}
		}
		if sessionID == uuid.Nil {
			session.Set(SESSION_ID, cache.StartSession().String())
			if err := session.Save(); err != nil {
				panic(err)
			}
		}
		// trigger pending handlers (necessary?)
		// c.Next()
	}
}
