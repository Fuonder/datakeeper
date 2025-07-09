package main

import (
	"context"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/cipher"
	"github.com/Fuonder/datakeeper.git/internal/config"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/server/connection/postgre"
	"github.com/Fuonder/datakeeper.git/internal/server/dbservices"
	"github.com/Fuonder/datakeeper.git/internal/server/service"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
)

func main() {
	flags, err := config.ParseServerFlags()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize(flags.LogLevel); err != nil {
		panic(fmt.Errorf("method main: %v", err))
	}
	logger.Log.Info("Flags parsed",
		zap.String("flags", flags.String()))

	logger.Log.Info("Starting service")
	if err = run(flags); err != nil {
		logger.Log.Fatal("", zap.Error(err))
	}
}

func run(flags config.ServerOptions) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	DBConn, err := postgre.NewConnection(ctx, flags.DatabaseDSN)
	if err != nil {
		return err
	}

	instance, mu, err := DBConn.GetDBInstance(ctx)
	if err != nil {
		return err
	}

	DBServices, err := dbservices.NewDatabaseServices([]byte(flags.HashKey), instance, mu)
	if err != nil {
		return err
	}

	cipherService, err := cipher.NewAES256Cipher([]byte(flags.AESKey))
	if err != nil {
		return err
	}

	srv, err := service.NewService(flags.APIAddr.String(), DBServices, cipherService)
	if err != nil {
		return err
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		err = srv.Run()
		if err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger.Log.Debug("exit with error", zap.Error(err))
		cancel()
		return err
	}
	return nil
}
