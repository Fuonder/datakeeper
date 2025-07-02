package dbservices

import (
	"database/sql"
	"fmt"
	"sync"

	//"github.com/Fuonder/datakeeper.git/internal/auth"
	//"github.com/Fuonder/datakeeper.git/internal/models"
	//"github.com/Fuonder/datakeeper.git/internal/orders"
	"github.com/Fuonder/datakeeper.git/internal/users"
)

type DatabaseServices struct {
	UserSrv users.UserService
	//WalletSrv wallets.WalletService
	//OrderSrv  orders.OrderService
	//AuthSrv   auth.AuthService
}

func NewDatabaseServices(secret []byte, db *sql.DB, mu *sync.RWMutex) (*DatabaseServices, error) {
	// init user, auth, cards, files, logins, text services

	return nil, fmt.Errorf("not implemented")
}
