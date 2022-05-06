package session

import (
	"encoding/hex"
	"strconv"

	cfg "github.com/ms-xy/go-common/configuration"
)

type SessionConfig struct {
	// ContextKey defines the key within the session where the middleware stores its data
	ContextKey string
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
	Config SessionConfig = SessionConfig{}
)

func init() {
	loader, err := cfg.LoadConfiguration(
		func() (cfg.ConfigurationLoader, error) { return cfg.LoadEnvConfiguration() },
		func() (cfg.ConfigurationLoader, error) {
			loader, _ := cfg.LoadYamlConfiguration("session-manager-settings.yaml")
			return loader, nil
		},
		func() (cfg.ConfigurationLoader, error) {
			loader, _ := cfg.LoadJsonConfiguration("session-manager-settings.json")
			return loader, nil
		},
	)
	if err != nil {
		panic(err)
	}
	loader.GetTypeSafeOrDefault("session.middleware.handlepanic", &Config.HandlePanic, false)

	loader.GetTypeSafeOrDefault("session.mysql.dsn", &Config.MySqlDataSource, "test-user:test-password@tcp(localhost)/go_om")
	loader.GetTypeSafeOrDefault("session.mysql.table", &Config.CookiePath, "sessions")

	loader.GetTypeSafeOrDefault("session.cookie.name", &Config.CookieName, "session")
	loader.GetTypeSafeOrDefault("session.cookie.domain", &Config.CookieDomain, "127.0.0.1")
	loader.GetTypeSafeOrDefault("session.cookie.path", &Config.CookiePath, "")
	loader.GetTypeSafeOrDefault("session.cookie.secure", &Config.CookieSecure, false)
	loader.GetTypeSafeOrDefault("session.cookie.httpOnly", &Config.CookieHttpOnly, true)
	loader.GetTypeSafeOrDefault("session.cookie.persist", &Config.CookiePersist, true)

	loader.GetTypeSafeOrDefault("session.maxAge", &Config.SessionLifetime, 86400)
	loader.GetTypeSafeOrDefault("session.idleTime", &Config.SessionIdletime, 1800)
	loader.GetTypeSafeOrDefault("session.contextKey", &Config.ContextKey, "session")

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
