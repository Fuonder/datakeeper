package auth

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

const (
	GetUserPasswordQuery = `SELECT pwd_hash FROM users WHERE login = $1;`
)

type DatabaseAuth interface {
	ValidateUserCredentials(ctx context.Context, user models.User) error
}

type DBAuth struct {
	db *sql.DB
	mu *sync.RWMutex
}

func NewDBAuth(db *sql.DB, mu *sync.RWMutex) (*DBAuth, error) {
	return &DBAuth{db: db, mu: mu}, nil
}

func (a *DBAuth) ValidateUserCredentials(ctx context.Context, user models.User) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var hashPassword string
	err := a.db.QueryRowContext(ctx, GetUserPasswordQuery, user.Login).Scan(&hashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrWrongCredentials
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(user.PwdHash))
	if err != nil {
		return models.ErrWrongCredentials
	}
	return nil
}
