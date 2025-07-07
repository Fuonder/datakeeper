package files

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Fuonder/datakeeper.git/internal/models"
	"sync"
)

type DBFiles struct {
	db *sql.DB
	mu *sync.RWMutex
}

const (
	CheckFileExistsQuery = `SELECT id FROM files_data WHERE id = $1 AND user_id = $2;`
	InsertNewFileQuery   = `
		INSERT INTO files_data (user_id, path, file_type, last_update, metadata)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	UpdateFileQuery = `
		UPDATE files_data
		SET path = $1, file_type = $2, last_update = $3, metadata = $4
		WHERE id = $5 AND user_id = $6;
	`
	GetFileQuery = `
		SELECT id, user_id, path, file_type, created_at, last_update, metadata
		FROM files_data WHERE id = $1 AND user_id = $2;
	`
	GetFilesQueryByUID = `
		SELECT id, user_id, path, file_type, created_at, last_update, metadata
		FROM files_data WHERE user_id = $1;
	`
)

func NewDBFiles(db *sql.DB, mu *sync.RWMutex) (*DBFiles, error) {
	return &DBFiles{db: db, mu: mu}, nil
}

func (d *DBFiles) AddOrUpdateFileRecord(ctx context.Context, file models.FileData) (models.FileData, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return models.FileData{}, err
	}
	defer tx.Rollback()

	var existingID int
	err = tx.QueryRowContext(ctx, CheckFileExistsQuery, file.ID, file.UserID).Scan(&existingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.FileData{}, err
	}

	if err == nil { // update existing
		file.ID = existingID
		if err := d.UpdateFileRecord(ctx, tx, file); err != nil {
			return models.FileData{}, err
		}
		if err := tx.Commit(); err != nil {
			return models.FileData{}, err
		}

		return file, nil
	}

	// New
	err = tx.QueryRowContext(ctx, InsertNewFileQuery,
		file.UserID, file.Path, file.FileType, file.LastUpdate, file.Metadata,
	).Scan(&file.ID)
	if err != nil {
		return models.FileData{}, err
	}
	if err := tx.Commit(); err != nil {
		return models.FileData{}, err
	}

	return file, nil
}

func (d *DBFiles) UpdateFileRecord(ctx context.Context, tx *sql.Tx, file models.FileData) error {
	res, err := tx.ExecContext(ctx, UpdateFileQuery,
		file.Path, file.FileType, file.LastUpdate, file.Metadata, file.ID, file.UserID,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (d *DBFiles) GetFileRecord(ctx context.Context, fileID int, userID int) (models.FileData, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var file models.FileData
	err := d.db.QueryRowContext(ctx, GetFileQuery, fileID, userID).
		Scan(&file.ID, &file.UserID, &file.Path, &file.FileType, &file.CreatedAt, &file.LastUpdate, &file.Metadata)
	if err != nil {
		return models.FileData{}, err
	}
	return file, nil
}

func (d *DBFiles) GetUserFileRecords(ctx context.Context, userID int) ([]models.FileData, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rows, err := d.db.QueryContext(ctx, GetFilesQueryByUID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []models.FileData
	for rows.Next() {
		var file models.FileData
		if err := rows.Scan(&file.ID, &file.UserID, &file.Path, &file.FileType, &file.CreatedAt, &file.LastUpdate, &file.Metadata); err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return files, nil
}
