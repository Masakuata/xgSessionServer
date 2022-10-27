package session

import (
	"time"
	"xgss/model"
)

const LIFETIME = 600

var sessionHolder = model.SyncSession{
	Sessions: make(map[string]model.Session),
}

var isCleaning = false

func sessionCleaner() {
	isCleaning = true

	for len(sessionHolder.Sessions) > 0 {
		now := time.Now()
		var closerCleaning = now.Add(1 * time.Hour)

		sessionHolder.Mutex.Lock()
		for token, session := range sessionHolder.Sessions {
			sessionTimestamp := time.Unix(session.Timestamp, 0)
			auxTimestamp := sessionTimestamp.Add(LIFETIME * time.Second)
			if auxTimestamp.Before(now) || auxTimestamp.Equal(now) {
				delete(sessionHolder.Sessions, token)
			} else {
				if sessionTimestamp.Before(closerCleaning) {
					closerCleaning = sessionTimestamp
				}
			}
		}

		sessionHolder.Mutex.Unlock()
		if len(sessionHolder.Sessions) > 0 {
			now = time.Now()
			time.Sleep(closerCleaning.Sub(now))
		}
	}

	isCleaning = false
}

func NewSession(token, email, password, role string, timestamp int64) {
	sessionHolder.Mutex.Lock()
	sessionHolder.Sessions[token] = model.Session{
		Email:     email,
		Password:  password,
		Role:      role,
		Timestamp: timestamp,
		Data:      map[string]any{},
	}
	sessionHolder.Mutex.Unlock()

	if !isCleaning {
		go sessionCleaner()
	}
}

func Exists(token string) bool {
	sessionHolder.Mutex.Lock()

	_, exists := sessionHolder.Sessions[token]
	defer sessionHolder.Mutex.Unlock()
	return exists
}

func AddData(token string, data map[string]any) {
	sessionHolder.Mutex.Lock()

	currentSession, exists := sessionHolder.Sessions[token]
	if exists {
		for key, value := range data {
			currentSession.Data[key] = value
		}
	}
	currentSession.Timestamp = time.Now().Unix()
	sessionHolder.Mutex.Unlock()
}

func GetData(token string) map[string]any {
	sessionHolder.Mutex.Lock()

	currentSession, _ := sessionHolder.Sessions[token]
	currentSession.Timestamp = time.Now().Unix()
	defer sessionHolder.Mutex.Unlock()
	return currentSession.Data
}

func UpdateLifetime(token string) {
	sessionHolder.Mutex.Lock()

	currentSession, exists := sessionHolder.Sessions[token]
	if exists {
		currentSession.Timestamp = time.Now().Unix()
	}
	sessionHolder.Mutex.Unlock()
}

func IsRole(token, role string) bool {
	sessionHolder.Mutex.Lock()

	currentSession, exists := sessionHolder.Sessions[token]

	defer sessionHolder.Mutex.Unlock()
	return exists && currentSession.Role == role
}

func Delete(token string) {
	sessionHolder.Mutex.Lock()

	delete(sessionHolder.Sessions, token)
	sessionHolder.Mutex.Unlock()
}
