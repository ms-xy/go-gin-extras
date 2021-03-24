package session

import (
	"errors"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
	color "gopkg.in/gookit/color.v1"
)

var (
	colorWarning = color.New(color.Yellow)

	ERR_SID_NOT_FOUND = errors.New("session ID not found in session cache")

	_bucket_sessions_ = []byte("sessions")
)

type SessionCache struct {
	path string
	db   *bbolt.DB
}

func NewSessionCache(path string) (*SessionCache, error) {
	if db, err := bbolt.Open(path, 0666, nil); err != nil {
		return nil, err
	} else {
		cache := &SessionCache{
			path: path,
			db:   db,
		}
		err := db.Update(func(tx *bbolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(_bucket_sessions_)
			return err
		})
		return cache, err
	}
}

func (this *SessionCache) New(session *Session) error {
	err := this.db.Update(func(tx *bbolt.Tx) error {
		sid := uuid.Must(uuid.NewRandom())
		_, err := tx.Bucket(_bucket_sessions_).CreateBucket(sid[:])
		if err == nil {
			session.SetID(sid)
		}
		return err
		// todo sync session variables to bucket
	})
	return err
}

func (this *SessionCache) Continue(session *Session) error {
	err := this.db.View(func(tx *bbolt.Tx) error {
		sid := session.GetID()
		bucket := tx.Bucket(_bucket_sessions_).Bucket(sid[:])
		if bucket == nil {
			session.SetID(uuid.Nil)
			return ERR_SID_NOT_FOUND
		}
		// todo sync bucket variables to session
		return nil
	})
	return err
}

func (this *SessionCache) Delete(session *Session) error {
	err := this.db.Update(func(tx *bbolt.Tx) error {
		sid := session.GetID()
		session.SetID(uuid.Nil)
		allBuckets := tx.Bucket(_bucket_sessions_)
		bucket := allBuckets.Bucket(sid[:])
		if bucket == nil {
			return ERR_SID_NOT_FOUND
		}
		err := allBuckets.DeleteBucket(sid[:])
		return err
	})
	return err
}
