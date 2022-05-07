package session

import (
	// "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/ms-xy/go-common/log"
	"github.com/ms-xy/go-gin-extras/middlewares/common"
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
	log.Debugf("token '%s' invalid", token)
	return nil, false, nil
}

func (ts *TestStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	ts.sessions[token] = Entry{b, expiry}
	return nil
}

func (ts *TestStore) PrintDump() {
	for token, entry := range ts.sessions {
		log.Debug("\t%s=%s\n", token, entry.data)
	}
}

func NewTestStore() *TestStore {
	return &TestStore{sessions: make(map[string]Entry)}
}

func TestMiddleware(t *testing.T) {
	log.SetLevel(log.LevelDebug)
	r := gin.New()
	mgr := scs.New()
	handler := SessionMiddleware(mgr, true)
	store := NewTestStore()
	mgr.Store = store
	r.Use(common.Logger())
	r.Use(handler)

	r.GET("/", func(c *gin.Context) {
		s := From(c)
		c.JSON(200, s)
	})

	// running a request without session_id should yield a session_id cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	sid := w.Result().Header.Get(HEADER_XSESSION)
	require.NotEmpty(t, sid)

	// running a request with session_id should retain the session_id
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set(HEADER_XSESSION, sid)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	sid2 := w.Result().Header.Get(HEADER_XSESSION)
	require.Equal(t, sid, sid2)
}
