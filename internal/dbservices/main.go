package dbservices

import (
	"database/sql"
	"github.com/Fuonder/datakeeper.git/internal/auth"
	"github.com/Fuonder/datakeeper.git/internal/objects/cards"
	"sync"

	//"github.com/Fuonder/datakeeper.git/internal/auth"
	//"github.com/Fuonder/datakeeper.git/internal/models"
	//"github.com/Fuonder/datakeeper.git/internal/orders"
	"github.com/Fuonder/datakeeper.git/internal/users"
)

type DatabaseServices struct {
	UserSrv users.UserService
	AuthSrv auth.Service
	CardSrv cards.Service
	//WalletSrv wallets.WalletService
	//OrderSrv  orders.OrderService
	//AuthSrv   auth.AuthService
}

func NewDatabaseServices(secret []byte, db *sql.DB, mu *sync.RWMutex) (*DatabaseServices, error) {
	// TODO: init user, auth, cards, files, logins, text services
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
	return s, nil
}
