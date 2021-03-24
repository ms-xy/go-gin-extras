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

func SessionMiddleware() gin.HandlerFunc {
	db, err := sql.Open("mysql", MySqlDataSource)
	if err != nil {
		panic(err)
	}
	store := mysqlstore.New(db, 10*time.Minute)
	sm := scs.NewManager(store)
	sm.Lifetime(time.Duration(SessionMaxAge) * time.Second)
	sm.IdleTimeout(30 * time.Minute)
	sm.Name(SessionCookie)
	sm.HttpOnly(SessionHttpOnly)
	sm.Secure(SessionSecure)
	sm.Domain(SessionDomain)
	sm.Persist(true)
	return func(c *gin.Context) {
		if HandlePanic {
			start := time.Now()
			defer func() {
				if r := recover(); r != nil {
					common.ResponseWriteError(c, r)
					common.WriteLogEntry(c, "SessionMiddleware", time.Since(start))
				}
			}()
		}
		c.Set(ContextKey, sm)
		session := sm.Load(c.Request)
		err := session.Touch(c.Writer)
		if err != nil {
			panic(err)
		}
		ctx := sm.AddToContext(c.Request.Context(), session)
		r := c.Request.WithContext(ctx)
		c.Request = r
		c.Next()
	}
}

func GetSession(c *gin.Context) *scs.Session {
	sm := c.MustGet(ContextKey).(*scs.Manager)
	s := sm.LoadFromContext(c)
	return s
}
