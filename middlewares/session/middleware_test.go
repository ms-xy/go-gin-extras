package session

import (
	// "encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type Entry struct {
	data   []byte
	expiry time.Time
}

type TestStore struct {
	sessions map[string]Entry
}

func (ts *TestStore) Delete(token string) error {
	delete(ts.sessions, token)
	return nil
}

func (ts *TestStore) Find(token string) ([]byte, bool, error) {
	if session, exists := ts.sessions[token]; exists {
		if session.expiry.After(time.Now()) {
			return session.data, true, nil
		}
	}
	log.Printf("== %s invalid\n", token)
	return nil, false, nil
}

func (ts *TestStore) Save(token string, b []byte, expiry time.Time) (err error) {
	ts.sessions[token] = Entry{b, expiry}
	return nil
}

func (ts *TestStore) PrintDump() {
	for token, entry := range ts.sessions {
		log.Printf("\t%s=%s\n", token, entry.data)
	}
}

func NewTestStore() *TestStore {
	return &TestStore{sessions: make(map[string]Entry)}
}

func TestMiddleware(t *testing.T) {
	r := gin.New()
	HandlePanic = true
	store := NewTestStore()
	_, handler := SessionMiddleware(store)
	r.Use(handler)

	r.GET("/", func(c *gin.Context) {
		s := GetSession(c)
		c.JSON(200, s)
	})

	// running a request without session_id should yield a session_id cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	var sid *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == SessionCookie {
			sid = cookie
		}
	}
	require.NotNil(t, sid)
	require.NotEmpty(t, sid.Value)

	// running a request with session_id should retain the session_id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.AddCookie(sid)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	var sid2 *http.Cookie
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == SessionCookie {
			sid2 = cookie
		}
	}
	require.Equal(t, sid.Value, sid2.Value)
}
