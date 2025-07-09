package auth

import (
	"context"
	"github.com/Fuonder/datakeeper.git/internal/models"
)

type Service interface {
	Register(ctx context.Context, newUser models.User) (token string, err error)
	Login(ctx context.Context, user models.User) (token string, err error)
	GetJWT(ctx context.Context, login string) (tokenString string, err error)
	ValidateJWT(ctx context.Context, tokenString string) error
	GetUIDFromJWT(ctx context.Context, tokenString string) (int, error)
}
