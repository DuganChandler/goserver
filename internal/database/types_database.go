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
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Token       string `json:"token"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}
