package session

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func withTmpFile(t *testing.T, fn func(string)) {
	file, err := ioutil.TempFile("", "*")
	if err != nil {
		t.Fatal(err)
	}
	path := file.Name()
	defer func() {
		if err := os.Remove(path); err != nil {
			t.Fatal(err)
		}
	}()
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	fn(path)
}

func TestCache(t *testing.T) {
	withTmpFile(t, func(path string) {
		// test NewSessionCache
		cache, err := NewSessionCache(path)
		require.Equal(t, nil, err, "NewSessionCache should yield nil error")
		require.NotEqual(t, nil, cache, "NewSessionCache should produce non-nil cache")

		// create dummy session object for testing
		session := &Session{}
		session.Valid = false
		session.ID = uuid.Nil

		// test opening session in cache
		err = cache.New(session)
		require.Equal(t, nil, err, "cache.New should yield nil error")
		require.NotEqual(t, uuid.Nil, session.GetID(), "cache.New should yield non-nil session id")
		sid := session.GetID()

		// test continuing session
		err = cache.Continue(session)
		require.Equal(t, nil, err, "cache.Continue should yield nil error")
		require.Equal(t, sid, session.GetID(), "cache.Continue should not change valid session id")

		//test deleting a session
		err = cache.Delete(session)
		require.Equal(t, nil, err, "cache.Delete should yield nil error")
		require.Equal(t, uuid.Nil, session.GetID(), "cache.Delete should set session id to nil")
		session2 := &Session{}
		session2.ID = sid
		err = cache.Continue(session2)
		require.Equal(t, ERR_SID_NOT_FOUND, err, "cache.Delete should yield ERR_SID_NOT_FOUND if attempted with invalid session id")
	})
}
