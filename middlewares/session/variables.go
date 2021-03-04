package session

import (
	"encoding/hex"
	"os"
	"strconv"
	"strings"

	env "github.com/ms-xy/go-common/environment"
)

const (
	SESSION_ID = "session_id"
)

var (
	// CachePath = path for bbolt key-val store
	// SessionDomain = domain for session cookies
	// SessionMaxAge = max age for session cookies, default 0
	// SessionSecret = secret for session cookie signing
	// SessionAesKey = aes key for session cookie encryption
	// SessionSecure = if true session cookies are only transmitted via https, default is false
	// SessionHttpOnly = if false session cookies are exposed via document.cookie, default is true
	CachePath       string = env.GetOrDefault("CACHE_PATH", "session.cache")
	SessionDomain   string = env.GetOrDefault("SESSION_DOMAIN", "127.0.0.1")
	SessionMaxAge   int    = mustParseInt(env.GetOrDefault("SESSION_MAX_AGE", "0"))
	SessionSecret   []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_SECRET", "688b9fdb0a43bf50c93efe6c06890a0ba9462c4662390b3a078901ff01841b23"))
	SessionAesKey   []byte = mustDecodeHex(env.GetOrDefault("SESSION_COOKIE_AES_KEY", "a6a07fbdf38d88fc92f794d570fd4b4d8c7b712e734e1adbba4690b981e28d5b"))
	SessionSecure   bool   = strings.ToLower(os.Getenv("SESSION_SECURE")) == "true"
	SessionHttpOnly bool   = (strings.ToLower(env.GetOrDefault("SESSION_HTTP_ONLY", "true")) == "true")
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
