package dbservices

import (
	"github.com/Fuonder/datakeeper.git/internal/server/auth"
	"github.com/Fuonder/datakeeper.git/internal/server/objects/cards"
	"github.com/Fuonder/datakeeper.git/internal/server/objects/files"
	"github.com/Fuonder/datakeeper.git/internal/server/objects/logins"
	"github.com/Fuonder/datakeeper.git/internal/server/objects/text"
	"github.com/Fuonder/datakeeper.git/internal/server/users"
)

type IDatabaseService interface {
	GetUserService() users.UserService
	GetAuthService() auth.Service
	GetCardService() cards.Service
	GetLoginService() logins.Service
	GetTextService() text.Service
	GetFileService() files.Service
}
