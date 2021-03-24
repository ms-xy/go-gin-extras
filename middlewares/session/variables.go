package session

import (
	"encoding/hex"
	"os"
	"strconv"
	"strings"

	env "github.com/ms-xy/go-common/environment"
)

var (
	// MySqlDataSource is the parameter used for creation of the store when calling DefaultSessionMiddleware().
	// See mysqlstore in scs for info on the required table schema.
	MySqlDataSource string = env.GetOrDefault("MYSQL_DATASOURCE", "test-user:test-password@tcp(localhost)/go_om")

	// SessionCookie is the name of the session cookie used, defaults to 'session'
	SessionCookie string = env.GetOrDefault("SESSION_COOKIE", "session")
	// SessionDomain is the name of the domain associated with the session cookie
	SessionDomain string = env.GetOrDefault("SESSION_DOMAIN", "127.0.0.1")
	// SessionMaxAge is the maximum session lifetime in seconds
	SessionMaxAge int = mustParseInt(env.GetOrDefault("SESSION_MAX_AGE", "86400")) // seconds, 24 hour default
	// SessionIdleTimeout is the maximum idle time before a non-active session is discarded
	SessionIdleTimeout int = mustParseInt(env.GetOrDefault("SESSION_IDLE_TIMEOUT", "1800")) // seconds, 30 mins default
	// SessionSecure sets wether or not the cookie should be https only
	SessionSecure bool = strings.ToLower(os.Getenv("SESSION_SECURE")) == "true"
	// SessionHttpOnly sets wether the cookie be accessible via javascript
	SessionHttpOnly bool = (strings.ToLower(env.GetOrDefault("SESSION_HTTP_ONLY", "true")) == "true")

	// SessionSecret []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_SECRET", "688b9fdb0a43bf50c93efe6c06890a0ba9462c4662390b3a078901ff01841b23"))
	// SessionAesKey []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_AES_KEY", "a6a07fbdf38d88fc92f794d570fd4b4d8c7b712e734e1adbba4690b981e28d5b"))
)

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
