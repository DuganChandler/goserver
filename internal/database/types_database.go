package database

import (
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type RefreshToken struct {
	Token    string        `json:"token"`
	Duration time.Duration `json:"duration"`
}

type User struct {
	Id           int          `json:"id"`
	Email        string       `json:"email"`
	Password     string       `json:"password"`
	Token        string       `json:"token"`
	RefreshToken RefreshToken `json:"refresh_token"`
}
