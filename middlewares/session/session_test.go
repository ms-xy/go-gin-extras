package session

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

func TestSession(t *testing.T) {
	withTmpFile(t, func(cachePath string) {
		cache, err := NewSessionCache(cachePath)
		require.Nil(t, err, "NewSessionCache should not return an error: %s", err)

		s := NewSession(cache)
		err = s.Start()
		require.Nil(t, err, "session.Start should not return an error: %s", err)
		require.NotNil(t, s.GetID(), "session.Start should set a valid session id")
		require.True(t, s.IsValid(), "session.Start should yield a valid session")

		sid := s.GetID()
		err = s.Start()
		require.Nil(t, err, "session.Start should not return an error: %s", err)
		require.Equal(t, sid, s.GetID(), "calling session.Start with existing sid should not modify sid")
		require.True(t, s.IsValid(), "session.Start should yield a valid session")

		err = s.Delete()
		require.Nil(t, err, "session.Delete should not return an error: %s", err)
		require.Equal(t, uuid.Nil, s.GetID(), "session.Delete should set session id to nil")
		require.False(t, s.IsValid(), "session.Delete should invalidate session")
		cache.db.View(func(tx *bbolt.Tx) error {
			bucket := tx.Bucket(_bucket_sessions_).Bucket(sid[:])
			require.Nil(t, bucket, "session.Delete must result in cache delete of session bucket")
			return nil
		})
	})
}
