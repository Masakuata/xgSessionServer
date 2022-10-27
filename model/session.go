package model

import "sync"

type Session struct {
	Email     string
	Password  string
	Role      string
	Timestamp int64
	Data      map[string]any
}

type SyncSession struct {
	Mutex    sync.Mutex
	Sessions map[string]Session
}
