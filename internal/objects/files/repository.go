package files

import (
	"context"
	"github.com/Fuonder/datakeeper.git/internal/models"
)

type Service interface {
	AddOrUpdateFileRecord(ctx context.Context, file models.FileData) (models.FileData, error)
	GetFileRecord(ctx context.Context, fileID int, userID int) (models.FileData, error)
	GetUserFileRecords(ctx context.Context, userID int) ([]models.FileData, error)
}
