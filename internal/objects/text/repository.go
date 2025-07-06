package text

import (
	"context"
	"database/sql"
	"github.com/Fuonder/datakeeper.git/internal/models"
)

type Service interface {
	AddNewTextRecord(ctx context.Context, text models.TextData) (models.TextData, error)
	UpdateTextRecord(ctx context.Context, tx *sql.Tx, text models.TextData) error
	GetTextRecord(ctx context.Context, textID int, userID int) (models.TextData, error)
	GetUserTextRecords(ctx context.Context, userID int) ([]models.TextData, error)
	DeleteTextRecord(ctx context.Context, textID int) error
}
