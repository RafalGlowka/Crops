package crops

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"sync"
)

type SessionData struct {
	UserId int64
}

type SessionManager struct {
	lock         sync.Mutex
	tokens       map[string]SessionData
	passwordSalt string
}

type Session interface {
	Set(key, value interface{}) error //set session value
	Get(key interface{}) interface{}  //get session value
	Delete(key interface{}) error     //delete session value
	SessionID() string                //back current sessionID
}

func (manager *SessionManager) sessionId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (manager *SessionManager) startSession(user *User) (userId int64, token string) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	sid := manager.sessionId()
	for {
		_, exist := manager.tokens[sid]
		if exist == false {
			break
		}
		sid = manager.sessionId()
	}
	sessionData := SessionData{UserId: user.id}
	manager.tokens[sid] = sessionData
	userId = user.id

	return userId, sid
}

func (manager *SessionManager) getSessionData(token string) *SessionData {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	data, exist := manager.tokens[token]
	if exist {
		return &data
	}
	return nil
}

func (manager *SessionManager) getPasswordHash(password string) string {
	str := manager.passwordSalt + password
	data := []byte(str)
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))

}
