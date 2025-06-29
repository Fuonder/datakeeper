package storage

type IStorage interface {
}

type IAuth interface {
	Login(username, password string) (token string, err error)
	Register(username, password string) (token string, err error)
}

type IReader interface {
	GetUserObjects()
}
