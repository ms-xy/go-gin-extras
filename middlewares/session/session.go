package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/google/uuid"
	"github.com/ms-xy/logtools"
)

const (
	SESSION_ID = "session_id"
)

type Session struct {
	session sessions.Session
	cache   *SessionCache
	ID      uuid.UUID         `json:"id"`
	Data    map[string]string `json:"data"`
	Valid   bool              `json:"valid"`
}

func NewSession(cache *SessionCache) *Session {
	session := new(Session)
	session.cache = cache
	session.Data = make(map[string]string)
	return session
}

func (this *Session) Start() error {
	if this.GetID() != uuid.Nil {
		if err := this.cache.Continue(this); err == nil {
			this.SetValid(true)
			return nil
		} else {
			logtools.Errorf("Cannot resume session(%s): %s", this.GetID().String(), err.Error())
			return err
		}
	}
	logtools.Debug("Creating new session")
	if err := this.cache.New(this); err == nil {
		this.SetValid(true)
		return nil
	} else {
		logtools.Errorf("Error creating new session: %s", err.Error())
		return err
	}
}

func (this *Session) Delete() error {
	err := this.cache.Delete(this)
	this.SetValid(false)
	return err
}

func (this *Session) GetID() uuid.UUID {
	return this.ID
}

func (this *Session) SetID(sid uuid.UUID) {
	this.ID = sid
}

func (this *Session) IsValid() bool {
	return this.Valid
}

func (this *Session) SetValid(v bool) {
	this.Valid = v
}
