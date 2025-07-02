package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// User представляет таблицу users
type User struct {
	ID         int       `json:"id"`
	Login      string    `json:"login"`
	PwdHash    string    `json:"pwd_hash"`
	CreatedAt  time.Time `json:"created_at"`
	LastUpdate time.Time `json:"last_update,omitempty"`
}

// LoginData представляет таблицу logins_data
type LoginData struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	ServiceName  string    `json:"service_name"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	LastUpdate   time.Time `json:"last_update,omitempty"`
	Metadata     string    `json:"metadata,omitempty"`
}

// TextData представляет таблицу text_data
type TextData struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Data       string    `json:"data"`
	CreatedAt  time.Time `json:"created_at"`
	LastUpdate time.Time `json:"last_update,omitempty"`
	Metadata   string    `json:"metadata,omitempty"`
}

// CreditCardData представляет таблицу credit_cards_data
type CreditCardData struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	CardID     string    `json:"card_id"`
	OwnerName  string    `json:"owner_name"`
	CreatedAt  time.Time `json:"created_at"`
	LastUpdate time.Time `json:"last_update,omitempty"`
	Metadata   string    `json:"metadata,omitempty"`
}

// FileData представляет таблицу files_data
type FileData struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Path       string    `json:"path"`
	FileType   string    `json:"file_type"`
	CreatedAt  time.Time `json:"created_at"`
	LastUpdate time.Time `json:"last_update,omitempty"`
	Metadata   string    `json:"metadata,omitempty"`
}
