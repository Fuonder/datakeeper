package text

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"sync"
)

type DBText struct {
	db *sql.DB
	mu *sync.RWMutex
}

const (
	CheckTextExistsQuery = `SELECT id FROM text_data WHERE id = $1 AND user_id = $2;`
	InsertNewTextQuery   = `
		INSERT INTO text_data (user_id, data, last_update, metadata) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id;
	`
	UpdateTextQuery = `
		UPDATE text_data 
		SET data = $1, last_update = $2, metadata = $3 
		WHERE id = $4 AND user_id = $5;
	`
	GetTextQuery = `
		SELECT id, user_id, data, created_at, last_update, metadata 
		FROM text_data WHERE id = $1 AND user_id = $2;
	`
	GetTextsQueryByUID = `
		SELECT id, user_id, data, created_at, last_update, metadata 
		FROM text_data WHERE user_id = $1;
	`
)

func NewDBText(db *sql.DB, mu *sync.RWMutex) (*DBText, error) {
	return &DBText{db: db, mu: mu}, nil
}

func (d *DBText) AddNewTextRecord(ctx context.Context, text models.TextData) (models.TextData, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.TextData{}, err
	}
	defer tx.Rollback()

	var existingID int
	err = tx.QueryRowContext(ctx, CheckTextExistsQuery, text.ID, text.UserID).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.TextData{}, err
	}
	if err == nil {
		text.ID = existingID
		if err := d.UpdateTextRecord(ctx, tx, text); err != nil {
			return models.TextData{}, err
		}
		if err := tx.Commit(); err != nil {
			return models.TextData{}, err
		}
		return text, nil
	}

	err = tx.QueryRowContext(ctx, InsertNewTextQuery, text.UserID, text.Data, text.LastUpdate, text.Metadata).Scan(&text.ID)
	if err != nil {
		return models.TextData{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.TextData{}, err
	}
	return text, nil
}

func (d *DBText) UpdateTextRecord(ctx context.Context, tx *sql.Tx, text models.TextData) error {
	res, err := tx.ExecContext(ctx, UpdateTextQuery, text.Data, text.LastUpdate, text.Metadata, text.ID, text.UserID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *DBText) GetTextRecord(ctx context.Context, textID int, userID int) (models.TextData, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var text models.TextData
	err := d.db.QueryRowContext(ctx, GetTextQuery, textID, userID).
		Scan(&text.ID, &text.UserID, &text.Data, &text.CreatedAt, &text.LastUpdate, &text.Metadata)
	if err != nil {
		return models.TextData{}, err
	}
	return text, nil
}

func (d *DBText) GetUserTextRecords(ctx context.Context, userID int) ([]models.TextData, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.QueryContext(ctx, GetTextsQueryByUID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var texts []models.TextData
	for rows.Next() {
		var text models.TextData
		if err := rows.Scan(&text.ID, &text.UserID, &text.Data, &text.CreatedAt, &text.LastUpdate, &text.Metadata); err != nil {
			return nil, err
		}
		texts = append(texts, text)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(texts) == 0 {
		return nil, models.ErrNoData
	}
	return texts, nil
}

func (d *DBText) DeleteTextRecord(ctx context.Context, textID int) error {
	return fmt.Errorf("not implemented")
}
