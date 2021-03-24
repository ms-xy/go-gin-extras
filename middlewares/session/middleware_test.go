package session

import (
	// "encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ms-xy/logtools"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	r := gin.New()
	HandlePanic = true
	r.Use(SessionMiddleware())

	r.GET("/", func(c *gin.Context) {
		s := GetSession(c)
		c.JSON(200, s)
	})

	// running a request without session data should yield a session_id cookie
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	// checking the session object is just for completeness sake
	// s2 := &Session{}
	// err := json.Unmarshal(w.Body.Bytes(), s2)
	// require.Nil(t, err)

	logtools.WithFields(logtools.Fields{"header": w.HeaderMap}).Println("result")
	t.Fail()
}
