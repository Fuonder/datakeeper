package router

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/server/service/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouterObject struct {
	h        handlers.Handlers
	chRouter chi.Router
}

func NewRouterObject(h handlers.Handlers) *RouterObject {
	return &RouterObject{h: h, chRouter: chi.NewRouter()}
}

func (r *RouterObject) GetRouter() (chi.Router, error) {
	if r.chRouter == nil {
		return nil, fmt.Errorf("router not initialized")
	}
	logger.Log.Debug("Configuring Router")
	r.chRouter.Use(middleware.Compress(5))

	r.chRouter.Route("/register", func(router chi.Router) {
		router.Post("/", logger.HanlderWithLogger(r.h.RegisterHandlerPost))
	})
	r.chRouter.Route("/login", func(router chi.Router) {
		router.Post("/", logger.HanlderWithLogger(r.h.LoginHandlerPost))
	})

	r.chRouter.Route("/data", func(router chi.Router) {
		router.Use(r.h.AuthMiddleware)
		router.Get("/", logger.HanlderWithLogger(r.h.DataHandlerGet))

		router.Route("/login", func(router chi.Router) {
			router.Post("/", logger.HanlderWithLogger(r.h.SaveLoginHandlerPost))
			router.Get("/{object_id}", logger.HanlderWithLogger(r.h.GetLoginHandlerGet))
		})

		router.Route("/text", func(router chi.Router) {
			router.Post("/", logger.HanlderWithLogger(r.h.SaveTextHandlerPost))
			router.Get("/{object_id}", logger.HanlderWithLogger(r.h.GetTextHandlerGet))
		})

		router.Route("/file", func(router chi.Router) {
			router.Post("/", logger.HanlderWithLogger(r.h.SaveFileHandlerPost))
			router.Get("/{object_id}", logger.HanlderWithLogger(r.h.GetFileHandlerGet))
		})

		router.Route("/card", func(router chi.Router) {
			router.Post("/", logger.HanlderWithLogger(r.h.SaveCardHandlerPost))
			router.Get("/{object_id}", logger.HanlderWithLogger(r.h.GetCardHandlerGet))
		})
	})

	logger.Log.Info("Successfully initialized Router")
	return r.chRouter, nil
}

/*
POST	/register	// reg user -> return token
POST	/login		// login user -> return token
GET		/data		// list all objects -> return json object name, type, id
POST	/data/login	// load new or modify existing login object -> return status code
GET 	/data/login // get login object with given id in url params ->
POST 	/data/text
GET 	/data/text
POST 	/data/file
GET 	/data/file
POST 	/data/card
GET 	/data/card

*/
