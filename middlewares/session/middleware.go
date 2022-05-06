package session

import (
	"database/sql"
	"time"

	"github.com/ms-xy/go-gin-extras/middlewares/common"
	"github.com/ms-xy/logtools"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// DefaultSessionMiddleware is a convenience wrapper around
// SessionMiddleware that creates a mysqlstore based version of it with
// parameters drawn from env, json or yaml config - refer to variables.go
// for details
func DefaultSessionMiddleware() (*scs.SessionManager, gin.HandlerFunc, GetterFunc) {
	db, err := sql.Open("mysql", Config.MySqlDataSource)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			token CHAR(43) PRIMARY KEY,
			data BLOB NOT NULL,
			expiry TIMESTAMP(6) NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
	`)
	if err != nil {
		panic(err)
	}
	mgr, handlerFunc, getterFunc := SessionMiddleware(Config.ContextKey, Config.HandlePanic)
	mgr.Store = mysqlstore.New(db)

	mgr.Cookie.Domain = Config.CookieDomain
	mgr.Cookie.HttpOnly = Config.CookieHttpOnly
	mgr.Cookie.Name = Config.CookieName
	mgr.Cookie.Path = Config.CookiePath
	mgr.Cookie.Persist = Config.CookiePersist
	mgr.Cookie.Secure = Config.CookieSecure
	// TODO mgr.Cookie.SameSite

	lifetime := time.Duration(Config.SessionLifetime) * time.Second
	if lifetime > 0 {
		mgr.Lifetime = lifetime
	}
	idletime := time.Duration(Config.SessionIdletime) * time.Second
	if idletime <= 0 {
		idletime = 30 * time.Minute
	}
	mgr.IdleTimeout = idletime

	return mgr, handlerFunc, getterFunc
}

// SessionMiddleware creates a session middleware wrapper using scs.Manager
// with the provided scs.Store as the backing session storage.
// The return values include the manager, useful for further customization and
// the produced handler function for inclusion in your handler chain via
// gin.Use(handler)
func SessionMiddleware(ctxKey string, handlePanic bool) (*scs.SessionManager, gin.HandlerFunc, GetterFunc) {
	mgr := scs.New()
	headerKey := "X-Session"

	return mgr, func(c *gin.Context) {
			// only if HandlePanic is set, register recovery function
			if handlePanic {
				start := time.Now()
				defer func() {
					if r := recover(); r != nil {
						common.ResponseWriteError(c, r)
						common.WriteLogEntry(c, "SessionMiddleware", time.Since(start))
					}
				}()
			}

			logtools.WithFields(logtools.Fields{
				"c.Request": c.Request,
			}).Warn("debug logging request")

			ctx, err := mgr.Load(c.Request.Context(), c.Request.Header.Get(headerKey))
			if err != nil {
				logtools.Error("An error occured:", err)
			} else {
				logtools.WithFields(logtools.Fields{
					"ctx": ctx,
				}).Warn("debug logging ctx")
			}

			// gin does not set request context on its own, ensure it's set
			// r := c.Request.WithContext(c)

			// // manager.Load calls load, which in turn checks if session exists
			// // and initializes a new one if not
			// session := mgr.Load(r.Context())

			// // new sessions are not persisted yet, touch should take care of that
			// err := session.Touch(c.Writer)
			// if err != nil {
			// 	panic(err)
			// }

			// // save session to context
			// c.Set(ctxKey, session)

			// // if handle panic is true, must call c.Next() for the handler to catch
			// if HandlePanic {
			// 	c.Next()
			// }

			/*
				headerKey := "X-Session"
				headerKeyExpiry := "X-Session-Expiry"

				ctx, err := mgr.Load(r.Context(), r.Header.Get(headerKey))
				if err != nil {
					log.Output(2, err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				bw := &bufferedResponseWriter{ResponseWriter: w}
				sr := r.WithContext(ctx)
				next.ServeHTTP(bw, sr)

				if s.Status(ctx) == scs.Modified {
					token, expiry, err := s.Commit(ctx)
					if err != nil {
						log.Output(2, err.Error())
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}

					w.Header().Set(headerKey, token)
					w.Header().Set(headerKeyExpiry, expiry.Format(http.TimeFormat))
				}

				if bw.code != 0 {
					w.WriteHeader(bw.code)
				}
				w.Write(bw.buf.Bytes())
			*/
		}, func(c *gin.Context) *scs.SessionManager {
			// GetSession returns the session object associated with the current
			// request context.
			return c.MustGet(ctxKey).(*scs.SessionManager)
		}
}

type GetterFunc func(c *gin.Context) *scs.SessionManager
