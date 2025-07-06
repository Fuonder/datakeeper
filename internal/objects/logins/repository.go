package logins

import (
	"context"
	"database/sql"
	"github.com/Fuonder/datakeeper.git/internal/models"
)

type Service interface {
	AddNewLoginRecord(ctx context.Context, login models.LoginData) (models.LoginData, error)
	UpdateLoginRecord(ctx context.Context, tx *sql.Tx, login models.LoginData) error
	GetLoginRecord(ctx context.Context, loginID int, userID int) (models.LoginData, error)
	GetUserLoginRecords(ctx context.Context, userID int) ([]models.LoginData, error)
	DeleteLoginRecord(ctx context.Context, loginID int) error
}
