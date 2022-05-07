package session

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/ms-xy/go-common/log"
	"github.com/ms-xy/go-gin-extras/middlewares/common"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// DefaultSessionMiddleware is a convenience wrapper around
// SessionMiddleware that creates a mysqlstore based version of it with
// parameters drawn from env, json or yaml config - refer to variables.go
// for details
func DefaultSessionMiddleware() gin.HandlerFunc {
	db, err := sql.Open("mysql", Configuration.MySqlDataSource)
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)

	sessionManager.Cookie.Domain = Configuration.CookieDomain
	sessionManager.Cookie.HttpOnly = Configuration.CookieHttpOnly
	sessionManager.Cookie.Name = Configuration.CookieName
	sessionManager.Cookie.Path = Configuration.CookiePath
	sessionManager.Cookie.Persist = Configuration.CookiePersist
	sessionManager.Cookie.Secure = Configuration.CookieSecure
	// TODO mgr.Cookie.SameSite

	lifetime := time.Duration(Configuration.SessionLifetime) * time.Second
	if lifetime > 0 {
		sessionManager.Lifetime = lifetime
	}
	idletime := time.Duration(Configuration.SessionIdletime) * time.Second
	if idletime <= 0 {
		idletime = 30 * time.Minute
	}
	sessionManager.IdleTimeout = idletime

	ginHandler := SessionMiddleware(sessionManager, Configuration.HandlePanic)
	return ginHandler
}

// SessionMiddleware creates a session middleware wrapper using scs.Manager
// with the provided scs.Store as the backing session storage.
// The return values include the manager, useful for further customization and
// the produced handler function for inclusion in your handler chain via
// gin.Use(handler)
func SessionMiddleware(sessionManager *scs.SessionManager, handlePanic bool) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		log.WithField("c.Request", c.Request).Warn("debug logging request")

		session, err := sessionManager.Load(c.Request.Context(), c.Request.Header.Get(HEADER_XSESSION))
		if err != nil {
			log.Panic(err)
		} else {
			log.WithField("ctx", session).Warn("debug logging ctx")
		}

		c.Set(SESSION, session)

		sessionToken, expiryTime, err := sessionManager.Commit(session)
		if err != nil {
			panic(err)
		}
		c.Writer.Header().Set(HEADER_XSESSION, sessionToken)
		c.Writer.Header().Set(HEADER_XSESSION_EXPIRY, expiryTime.Format(http.TimeFormat))

		c.Next()

		log.WithFields(log.F{
			"token":  sessionToken,
			"expiry": expiryTime,
		}).Warn("SessionData")

		/*
			headerKey := "X-Session"
			headerKeyExpiry := "X-Session-Expiry"

			ctx, err := mgr.Load(r.Context(), r.Header.Get(headerKey))
			if err != nil {
				log.Output(2, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// // replace "inner" ctx with session (wrap around)
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
	}
}

func From(c *gin.Context) context.Context {
	// GetSession returns the session object associated with the current
	// request context.
	return c.MustGet(SESSION).(context.Context)
}
