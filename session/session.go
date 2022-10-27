package session

import (
	"sync"
	"time"
)

type Session struct {
	email     string
	password  string
	timestamp int64
	data      map[string]any
}

type SyncSession struct {
	mutex    sync.Mutex
	sessions map[string]Session
}

const LIFETIME = 600

var sessionHolder = SyncSession{
	sessions: make(map[string]Session),
}

var isCleaning = false

func sessionCleaner() {
	isCleaning = true

	for len(sessionHolder.sessions) > 0 {
		now := time.Now()
		var closerCleaning = now.Add(1 * time.Hour)

		sessionHolder.mutex.Lock()
		for token, session := range sessionHolder.sessions {
			sessionTimestamp := time.Unix(session.timestamp, 0)
			auxTimestamp := sessionTimestamp.Add(LIFETIME * time.Second)
			if auxTimestamp.Before(now) || auxTimestamp.Equal(now) {
				delete(sessionHolder.sessions, token)
			} else {
				if sessionTimestamp.Before(closerCleaning) {
					closerCleaning = sessionTimestamp
				}
			}
		}

		sessionHolder.mutex.Unlock()
		if len(sessionHolder.sessions) > 0 {
			now = time.Now()
			time.Sleep(closerCleaning.Sub(now))
		}
	}

	isCleaning = false
}

func NewSession(token, email, password string, timestamp int64) {
	sessionHolder.mutex.Lock()
	sessionHolder.sessions[token] = Session{
		email:     email,
		password:  password,
		timestamp: timestamp,
		data:      map[string]any{},
	}
	sessionHolder.mutex.Unlock()

	if !isCleaning {
		go sessionCleaner()
	}
}

func Exists(token string) bool {
	sessionHolder.mutex.Lock()

	_, exists := sessionHolder.sessions[token]
	defer sessionHolder.mutex.Unlock()
	return exists
}

func AddData(token string, data map[string]any) {
	sessionHolder.mutex.Lock()

	currentSession, exists := sessionHolder.sessions[token]
	if exists {
		for key, value := range data {
			currentSession.data[key] = value
		}
	}
	currentSession.timestamp = time.Now().Unix()
	sessionHolder.mutex.Unlock()
}

func GetData(token string) map[string]any {
	sessionHolder.mutex.Lock()

	currentSession, _ := sessionHolder.sessions[token]
	currentSession.timestamp = time.Now().Unix()
	defer sessionHolder.mutex.Unlock()
	return currentSession.data
}

func UpdateLifetime(token string) {
	sessionHolder.mutex.Lock()

	currentSession, exists := sessionHolder.sessions[token]
	if exists {
		currentSession.timestamp = time.Now().Unix()
	}
	sessionHolder.mutex.Unlock()
}

func Delete(token string) {
	sessionHolder.mutex.Lock()

	delete(sessionHolder.sessions, token)
	sessionHolder.mutex.Unlock()
}
