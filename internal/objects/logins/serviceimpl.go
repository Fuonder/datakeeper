package logins

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"sync"
)

type DBLogins struct {
	db *sql.DB
	mu *sync.RWMutex
}

const (
	CheckLoginExistsQuery = `SELECT id FROM logins_data WHERE id = $1 AND user_id = $2;`
	InsertNewLoginQuery   = `
		INSERT INTO logins_data (user_id, service_name, login, password_hash, last_update, metadata)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`
	UpdateLoginQuery = `
		UPDATE logins_data 
		SET service_name = $1, login = $2, password_hash = $3, last_update = $4, metadata = $5
		WHERE id = $6 AND user_id = $7;
	`
	GetLoginQuery = `
		SELECT id, user_id, service_name, login, password_hash, created_at, last_update, metadata
		FROM logins_data WHERE id = $1 AND user_id = $2;
	`
	GetLoginsQueryByUID = `
		SELECT id, user_id, service_name, login, password_hash, created_at, last_update, metadata
		FROM logins_data WHERE user_id = $1;
	`
)

func NewDBLogins(db *sql.DB, mu *sync.RWMutex) (*DBLogins, error) {
	return &DBLogins{db: db, mu: mu}, nil
}

func (l *DBLogins) AddNewLoginRecord(ctx context.Context, login models.LoginData) (models.LoginData, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return models.LoginData{}, err
	}
	defer tx.Rollback()

	var existingID int
	err = tx.QueryRowContext(ctx, CheckLoginExistsQuery, login.ID, login.UserID).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.LoginData{}, err
	}
	if err == nil {
		login.ID = existingID
		if err := l.UpdateLoginRecord(ctx, tx, login); err != nil {
			return models.LoginData{}, err
		}
		if err := tx.Commit(); err != nil {
			return models.LoginData{}, err
		}
		return login, nil
	}

	err = tx.QueryRowContext(ctx, InsertNewLoginQuery,
		login.UserID, login.ServiceName, login.Login, login.PasswordHash, login.LastUpdate, login.Metadata,
	).Scan(&login.ID)
	if err != nil {
		return models.LoginData{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.LoginData{}, err
	}
	return login, nil
}

func (l *DBLogins) UpdateLoginRecord(ctx context.Context, tx *sql.Tx, login models.LoginData) error {
	res, err := tx.ExecContext(ctx, UpdateLoginQuery,
		login.ServiceName, login.Login, login.PasswordHash, login.LastUpdate, login.Metadata,
		login.ID, login.UserID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (l *DBLogins) GetLoginRecord(ctx context.Context, loginID int, userID int) (models.LoginData, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var login models.LoginData
	err := l.db.QueryRowContext(ctx, GetLoginQuery, loginID, userID).
		Scan(&login.ID, &login.UserID, &login.ServiceName, &login.Login, &login.PasswordHash,
			&login.CreatedAt, &login.LastUpdate, &login.Metadata)
	if err != nil {
		return models.LoginData{}, err
	}
	return login, nil
}

func (l *DBLogins) GetUserLoginRecords(ctx context.Context, userID int) ([]models.LoginData, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	rows, err := l.db.QueryContext(ctx, GetLoginsQueryByUID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logins []models.LoginData
	for rows.Next() {
		var login models.LoginData
		if err := rows.Scan(&login.ID, &login.UserID, &login.ServiceName, &login.Login, &login.PasswordHash,
			&login.CreatedAt, &login.LastUpdate, &login.Metadata); err != nil {
			return nil, err
		}
		logins = append(logins, login)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(logins) == 0 {
		return nil, models.ErrNoData
	}
	return logins, nil
}

func (l *DBLogins) DeleteLoginRecord(ctx context.Context, loginID int) error {
	return fmt.Errorf("not implemented")
}
