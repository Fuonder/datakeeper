package cards

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"sync"
)

const (
	InsertNewCardQuery = `INSERT INTO credit_cards_data 
		(user_id, card_id, owner_name, last_update, metadata) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
		`
	UpdateCardQuery = `UPDATE credit_cards_data
		SET card_id = $1,
			owner_name = $2,
			last_update = $3,
			metadata = $4
		WHERE id = $5 AND user_id = $6;
		`
	GetCardQuery = `
		SELECT id, user_id, card_id, owner_name, created_at, last_update, metadata
		FROM credit_cards_data
		WHERE id = $1 AND user_id = $2;
		`
	GetCardsQueryByUID = `
		SELECT id, user_id, card_id, owner_name, created_at, last_update, metadata
		FROM credit_cards_data
		WHERE user_id = $1;
		`
	DeleteCardQuery = `
		DELETE FROM credit_cards_data
		WHERE id = $1 AND user_id = $2;
		`
	CheckCardExistsQuery = `
		SELECT id FROM credit_cards_data
		WHERE id = $1 AND user_id = $2;
	`
)

//type DatabaseCards interface {
//	AddNewCardRecord(ctx context.Context, cardObject models.CreditCardData) error
//	UpdateCardRecord(ctx context.Context, cardObject models.CreditCardData) error
//	GetCardRecord(ctx context.Context, cardID int, userID int) (cardRecord models.CreditCardData, err error)
//	GetUserCardRecords(ctx context.Context, userID int) (cardRecords []models.CreditCardData, err error)
//	DeleteCardRecord(ctx context.Context, cardID int) error
//}

type DBCards struct {
	db *sql.DB
	mu *sync.RWMutex
}

func NewDBCards(db *sql.DB, mu *sync.RWMutex) (*DBCards, error) {
	return &DBCards{db: db, mu: mu}, nil
}

func (c *DBCards) AddNewCardRecord(
	ctx context.Context,
	card models.CreditCardData) (models.CreditCardData, error) {

	c.mu.Lock()
	defer c.mu.Unlock()

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return models.CreditCardData{}, fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback()

	var existingID int
	err = tx.QueryRowContext(ctx, CheckCardExistsQuery, card.ID, card.UserID).Scan(&existingID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.CreditCardData{}, fmt.Errorf("failed to check card existence: %w", err)
	}

	if err == nil {
		// Карта уже есть — делаем обновление
		card.ID = existingID
		if err := c.UpdateCardRecord(ctx, tx, card); err != nil {
			return models.CreditCardData{}, fmt.Errorf("update via add failed: %w", err)
		}
		if err = tx.Commit(); err != nil {
			return models.CreditCardData{}, fmt.Errorf("commit failed after update: %w", err)
		}
		return card, nil
	}

	// Вставка новой карты
	err = tx.QueryRowContext(
		ctx, InsertNewCardQuery,
		card.UserID, card.CardID, card.OwnerName, card.LastUpdate, card.Metadata,
	).Scan(&card.ID)
	if err != nil {
		return models.CreditCardData{}, fmt.Errorf("insert failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.CreditCardData{}, fmt.Errorf("commit failed: %w", err)
	}

	return card, nil
}

func (c *DBCards) UpdateCardRecord(ctx context.Context, tx *sql.Tx, card models.CreditCardData) error {
	res, err := tx.ExecContext(
		ctx, UpdateCardQuery,
		card.CardID, card.OwnerName, card.LastUpdate, card.Metadata,
		card.ID, card.UserID,
	)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check rows affected: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (c *DBCards) GetCardRecord(ctx context.Context, cardID int, userID int) (models.CreditCardData, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var card models.CreditCardData
	err := c.db.QueryRowContext(ctx, GetCardQuery, cardID, userID).
		Scan(&card.ID, &card.UserID, &card.CardID, &card.OwnerName, &card.CreatedAt, &card.LastUpdate, &card.Metadata)
	if err != nil {
		return models.CreditCardData{}, fmt.Errorf("get failed: %w", err)
	}
	return card, nil
}

// GetUserCardRecords TODO: CHECK WHEN ALL OBJECTS IS DONE
func (c *DBCards) GetUserCardRecords(ctx context.Context, userID int) ([]models.CreditCardData, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rows, err := c.db.QueryContext(ctx, GetCardsQueryByUID, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var cards []models.CreditCardData
	for rows.Next() {
		var card models.CreditCardData
		if err := rows.Scan(&card.ID, &card.UserID, &card.CardID, &card.OwnerName, &card.CreatedAt, &card.LastUpdate, &card.Metadata); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	if len(cards) == 0 {
		return nil, models.ErrNoData
	}
	return cards, nil
}

func (c *DBCards) DeleteCardRecord(ctx context.Context, cardID int) error {
	return fmt.Errorf("not implemented")
}
