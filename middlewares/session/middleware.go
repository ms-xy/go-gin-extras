package session

import (
	"database/sql"
	"time"

	"github.com/ms-xy/go-gin-extras/middlewares/common"

	"github.com/alexedwards/scs"
	"github.com/alexedwards/scs/stores/mysqlstore"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const (
	ContextKey = "github.com/ms-xy/go-gin-extras/middlewares/session-manager"
)

var (
	// HandlePanic indicates whether or not the session middleware will handle
	// its own panics or not. Default is false. Setting it to true will result
	// in the same behavior as using gin.Use(Logging(), Recovery()) from this
	// repository.
	HandlePanic = false
)

func DefaultSessionMiddleware() gin.HandlerFunc {
	db, err := sql.Open("mysql", MySqlDataSource)
	if err != nil {
		panic(err)
	}
	store := mysqlstore.New(db, 10*time.Minute)
	return SessionMiddleware(store)
}

func SessionMiddleware(store scs.Store) gin.HandlerFunc {
	sm := scs.NewManager(store)
	lifetime := time.Duration(SessionMaxAge) * time.Second
	if lifetime > 0 {
		sm.Lifetime(lifetime)
	}
	idletime := time.Duration(SessionIdleTimeout) * time.Second
	if idletime <= 0 {
		idletime = 30 * time.Minute
	}
	sm.IdleTimeout(idletime)
	sm.Name(SessionCookie)
	sm.HttpOnly(SessionHttpOnly)
	sm.Secure(SessionSecure)
	sm.Domain(SessionDomain)
	sm.Persist(true)

	return func(c *gin.Context) {
		// only if HandlePanic is set, register recovery function
		if HandlePanic {
			start := time.Now()
			defer func() {
				if r := recover(); r != nil {
					common.ResponseWriteError(c, r)
					common.WriteLogEntry(c, "SessionMiddleware", time.Since(start))
				}
			}()
		}

		// gin does not set request context on its own, ensure it's set
		r := c.Request.WithContext(c)

		// manager.Load calls load, which in turn checks if session exists
		// and initializes a new one if not
		session := sm.Load(r)

		// new sessions are not persisted yet, touch should take care of that
		err := session.Touch(c.Writer)
		if err != nil {
			panic(err)
		}

		// save session to context
		c.Set(ContextKey, session)

		// if handle panic is true, must call c.Next() for the handler to catch
		if HandlePanic {
			c.Next()
		}
	}
}

func GetSession(c *gin.Context) *scs.Session {
	return c.MustGet(ContextKey).(*scs.Session)
}
