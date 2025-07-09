package dbservices

import (
	"database/sql"
	auth "github.com/Fuonder/datakeeper.git/internal/server/auth"
	cards "github.com/Fuonder/datakeeper.git/internal/server/objects/cards"
	files "github.com/Fuonder/datakeeper.git/internal/server/objects/files"
	logins "github.com/Fuonder/datakeeper.git/internal/server/objects/logins"
	text "github.com/Fuonder/datakeeper.git/internal/server/objects/text"
	users "github.com/Fuonder/datakeeper.git/internal/server/users"
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

func NewDatabaseServices(secret []byte, db *sql.DB, mu *sync.RWMutex) (IDatabaseService, error) {
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

func (s *DatabaseServices) GetUserService() users.UserService {
	return s.UserSrv
}

func (s *DatabaseServices) GetAuthService() auth.Service {
	return s.AuthSrv
}

func (s *DatabaseServices) GetCardService() cards.Service {
	return s.CardSrv
}

func (s *DatabaseServices) GetLoginService() logins.Service {
	return s.LoginSrv
}

func (s *DatabaseServices) GetTextService() text.Service {
	return s.TextSrv
}

func (s *DatabaseServices) GetFileService() files.Service {
	return s.FileSrv
}
