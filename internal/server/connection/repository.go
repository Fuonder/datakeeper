package connection

import (
	"context"
	"database/sql"
	"sync"
)

type IDBConnection interface {
	ConnectCtx(ctx context.Context) error
	MigrateCtx(ctx context.Context) error
	Close() error
	GetDBInstance() (*sql.DB, *sync.RWMutex, error)
}
