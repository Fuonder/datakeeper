package dbservices

import (
	"database/sql"
	"github.com/Fuonder/datakeeper.git/internal/auth"
	"github.com/Fuonder/datakeeper.git/internal/objects/cards"
	"github.com/Fuonder/datakeeper.git/internal/objects/files"
	"github.com/Fuonder/datakeeper.git/internal/objects/logins"
	"github.com/Fuonder/datakeeper.git/internal/objects/text"
	"github.com/Fuonder/datakeeper.git/internal/users"
	"sync"
)

type DatabaseServices struct {
	UserSrv  users.UserService
	AuthSrv  auth.Service
	CardSrv  cards.Service
	LoginSrv logins.Service
	TextSrv  text.Service
	FileSrv  files.Service
}

func NewDatabaseServices(secret []byte, db *sql.DB, mu *sync.RWMutex) (*DatabaseServices, error) {
	s := &DatabaseServices{}

	DBUsers, err := users.NewDBUsers(db, mu)
	if err != nil {
		return s, err
	}

	s.UserSrv = users.NewUService(DBUsers)

	DBAuth, err := auth.NewDBAuth(db, mu)
	if err != nil {
		return s, err
	}

	s.AuthSrv = auth.NewAService(DBUsers, DBAuth, secret)

	DBCards, err := cards.NewDBCards(db, mu)
	if err != nil {
		return s, err
	}
	s.CardSrv = DBCards

	DBLogins, err := logins.NewDBLogins(db, mu)
	if err != nil {
		return s, err
	}
	s.LoginSrv = DBLogins

	DBText, err := text.NewDBText(db, mu)
	if err != nil {
		return s, err
	}
	s.TextSrv = DBText

	DBFile, err := files.NewDBFiles(db, mu)
	if err != nil {
		return s, err
	}
	s.FileSrv = DBFile

	return s, nil
}
