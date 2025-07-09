package cards

import (
	"context"
	"database/sql"
	"github.com/Fuonder/datakeeper.git/internal/models"
)

type Service interface {
	AddNewCardRecord(ctx context.Context, cardObject models.CreditCardData) (cardRecord models.CreditCardData, err error)
	UpdateCardRecord(ctx context.Context, tx *sql.Tx, card models.CreditCardData) error
	GetCardRecord(ctx context.Context, cardID int, userID int) (cardRecord models.CreditCardData, err error)
	GetUserCardRecords(ctx context.Context, userID int) (cardRecords []models.CreditCardData, err error)
	DeleteCardRecord(ctx context.Context, cardID int) error
}
