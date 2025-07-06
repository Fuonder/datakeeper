package service

import (
	"github.com/Fuonder/datakeeper.git/internal/cipher"
	"github.com/Fuonder/datakeeper.git/internal/dbservices"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	. "github.com/Fuonder/datakeeper.git/internal/service/handlers"
	. "github.com/Fuonder/datakeeper.git/internal/service/router"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	apiSrv     http.Server
	DBServices *dbservices.DatabaseServices
}

func NewService(APIAddr string, DBServices *dbservices.DatabaseServices, cipherService cipher.Service) (*Service, error) {

	h := NewHandlers(DBServices, cipherService)
	rObj := NewRouterObject(*h)
	router, err := rObj.GetRouter()
	if err != nil {
		return nil, err
	}

	service := &Service{
		apiSrv: http.Server{
			Addr:    APIAddr,
			Handler: router,
		},
		DBServices: DBServices,
	}
	return service, nil
}

func (s *Service) Run() error {
	logger.Log.Info("API Listening at",
		zap.String("Addr", s.apiSrv.Addr))
	return s.apiSrv.ListenAndServe()
}
