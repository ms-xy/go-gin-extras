package session

import (
	"encoding/hex"
	"strconv"

	cfg "github.com/ms-xy/go-common/configuration"
	"github.com/ms-xy/go-common/log"
)

// SESSION is the key which this middleware uses to store it's data in a
// gin.Context map
const (
	SESSION                = "__session__"
	HEADER_XSESSION        = "X-Session"
	HEADER_XSESSION_EXPIRY = "X-Session-Expiry"
)

type SessionConfig struct {
	// HandlePanic indicates whether or not the session middleware will handle
	// its own panics or not. Default is false. Setting it to true will result
	// in the same behavior as using gin.Use(Logging(), Recovery()) from this
	// repository.
	HandlePanic bool

	// MySqlDataSource is the parameter used for creation of the store when calling DefaultSessionMiddleware().
	// See mysqlstore in scs for info on the required table schema.
	MySqlDataSource string
	MySqlTableName  string

	// CookieName is the name of the session cookie used, defaults to 'session'
	CookieName string
	// CookieDomain is the name of the domain associated with the session cookie
	CookieDomain string
	// SessionDomain is the name of the domain associated with the session cookie
	CookiePath string
	// {true,false}: Is the cookie HTTPS only? [default=false]
	CookieSecure bool
	// {true,false}: Shall the cookie be restricted to http or is it also accessible via Javascript? [default=true]
	CookieHttpOnly bool
	// {true,false}: Should the cookie persist across browser sessions? [default=true]
	CookiePersist bool

	// SessionLifetime is the maximum session lifetime in seconds
	SessionLifetime int
	// SessionIdletime is the maximum idle time before a non-active session is discarded
	SessionIdletime int

	// SessionSecret []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_SECRET", "688b9fdb0a43bf50c93efe6c06890a0ba9462c4662390b3a078901ff01841b23"))
	// SessionAesKey []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_AES_KEY", "a6a07fbdf38d88fc92f794d570fd4b4d8c7b712e734e1adbba4690b981e28d5b"))
}

var (
	Configuration SessionConfig = SessionConfig{}
)

func LoadConfig(filepaths ...string) {
	loader := cfg.NewCombinedLoader()
	loader.CanLoadEnv()
	for _, path := range filepaths {
		log.Infof("loading config file '%s'", path)
		if err := loader.LoadYaml(path); err != nil {
			if err := loader.LoadJSON(path); err != nil {
				log.WithField("filepath", path).Error("invalid file type, need yaml or json")
			}
		}
	}
	loader.DumpConfig()

	loader.GetTypeSafeOrDefault("gin.middleware.session.handlepanic", &Configuration.HandlePanic, false)

	loader.GetTypeSafeOrDefault("gin.middleware.session.mysql.dsn", &Configuration.MySqlDataSource, "test-user:test-password@tcp(localhost)/go_om")
	loader.GetTypeSafeOrDefault("gin.middleware.session.mysql.table", &Configuration.CookiePath, "sessions")

	loader.GetTypeSafeOrDefault("gin.middleware.session.cookie.name", &Configuration.CookieName, "session")
	loader.GetTypeSafeOrDefault("gin.middleware.session.cookie.domain", &Configuration.CookieDomain, "127.0.0.1")
	loader.GetTypeSafeOrDefault("gin.middleware.session.cookie.path", &Configuration.CookiePath, "")
	loader.GetTypeSafeOrDefault("gin.middleware.session.cookie.secure", &Configuration.CookieSecure, false)
	loader.GetTypeSafeOrDefault("gin.middleware.session.cookie.httpOnly", &Configuration.CookieHttpOnly, true)
	loader.GetTypeSafeOrDefault("gin.middleware.session.cookie.persist", &Configuration.CookiePersist, true)

	loader.GetTypeSafeOrDefault("gin.middleware.session.maxAge", &Configuration.SessionLifetime, 86400)
	loader.GetTypeSafeOrDefault("gin.middleware.session.idleTime", &Configuration.SessionIdletime, 1800)

	// SessionSecret []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_SECRET", "688b9fdb0a43bf50c93efe6c06890a0ba9462c4662390b3a078901ff01841b23"))
	// SessionAesKey []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_AES_KEY", "a6a07fbdf38d88fc92f794d570fd4b4d8c7b712e734e1adbba4690b981e28d5b"))
}

func mustParseInt(s string) int {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(v)
}

func mustDecodeHex(s string) []byte {
	v, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return v
}
