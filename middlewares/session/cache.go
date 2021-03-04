package session

import (
	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

type SessionCache struct {
	path string
	db   *bbolt.DB
}

const (
	BUCKET_SESSIONS = "sessions"
)

func NewSessionCache(path string) *SessionCache {
	if db, err := bbolt.Open(path, 0666, nil); err != nil {
		panic(err)
	} else {
		return &SessionCache{
			path: path,
			db:   db,
		}
	}
}

func (this *SessionCache) Exists(id uuid.UUID) (exists bool) {
	exists = false
	this.db.View(func(tx *bbolt.Tx) error {
		if sessions := tx.Bucket([]byte(BUCKET_SESSIONS)); sessions != nil {
			exists = sessions.Bucket(id[:]) != nil
		}
		return nil
	})
	return
}

func (this *SessionCache) StartSession() uuid.UUID {
	for {
		if sessionID, err := uuid.NewRandom(); err == nil {
			if err := this.SessionPersist(sessionID); err == nil {
				return sessionID
			} else {
				// already exists -> new attempt
			}
		} else {
			panic(err)
		}
	}
}

func (this *SessionCache) SessionPersist(id uuid.UUID) error {
	return this.db.Update(func(tx *bbolt.Tx) error {
		if sessions, err := tx.CreateBucketIfNotExists([]byte(BUCKET_SESSIONS)); err != nil {
			if _, err := sessions.CreateBucket(id[:]); err != nil {
				return err
			}
		}
		return nil
	})
}
