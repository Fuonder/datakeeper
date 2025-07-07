package dbservices

import (
	"database/sql"
	auth2 "github.com/Fuonder/datakeeper.git/internal/server/auth"
	cards2 "github.com/Fuonder/datakeeper.git/internal/server/objects/cards"
	files2 "github.com/Fuonder/datakeeper.git/internal/server/objects/files"
	logins2 "github.com/Fuonder/datakeeper.git/internal/server/objects/logins"
	text2 "github.com/Fuonder/datakeeper.git/internal/server/objects/text"
	users2 "github.com/Fuonder/datakeeper.git/internal/server/users"
	"sync"
)

type DatabaseServices struct {
	UserSrv  users2.UserService
	AuthSrv  auth2.Service
	CardSrv  cards2.Service
	LoginSrv logins2.Service
	TextSrv  text2.Service
	FileSrv  files2.Service
}

func NewDatabaseServices(secret []byte, db *sql.DB, mu *sync.RWMutex) (*DatabaseServices, error) {
	s := &DatabaseServices{}

	DBUsers, err := users2.NewDBUsers(db, mu)
	if err != nil {
		return s, err
	}

	s.UserSrv = users2.NewUService(DBUsers)

	DBAuth, err := auth2.NewDBAuth(db, mu)
	if err != nil {
		return s, err
	}

	s.AuthSrv = auth2.NewAService(DBUsers, DBAuth, secret)

	DBCards, err := cards2.NewDBCards(db, mu)
	if err != nil {
		return s, err
	}
	s.CardSrv = DBCards

	DBLogins, err := logins2.NewDBLogins(db, mu)
	if err != nil {
		return s, err
	}
	s.LoginSrv = DBLogins

	DBText, err := text2.NewDBText(db, mu)
	if err != nil {
		return s, err
	}
	s.TextSrv = DBText

	DBFile, err := files2.NewDBFiles(db, mu)
	if err != nil {
		return s, err
	}
	s.FileSrv = DBFile

	return s, nil
}
