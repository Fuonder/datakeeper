package main

import (
	"context"
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/cipher"
	"github.com/Fuonder/datakeeper.git/internal/connection/postgre"
	"github.com/Fuonder/datakeeper.git/internal/dbservices"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	"github.com/Fuonder/datakeeper.git/internal/service"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
)

func main() {
	err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Initialize(Flags.LogLevel); err != nil {
		panic(fmt.Errorf("method main: %v", err))
	}
	logger.Log.Info("Flags parsed",
		zap.String("flags", Flags.String()))

	logger.Log.Info("Starting service")
	if err = run(); err != nil {
		logger.Log.Fatal("", zap.Error(err))
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	DBConn, err := postgre.NewConnection(ctx, Flags.DatabaseDSN)
	if err != nil {
		return err
	}

	instance, mu, err := DBConn.GetDBInstance(ctx)
	if err != nil {
		return err
	}

	DBServices, err := dbservices.NewDatabaseServices([]byte(Flags.HashKey), instance, mu)
	if err != nil {
		return err
	}

	cipherService, err := cipher.NewAES256Cipher([]byte(Flags.AESKey))
	if err != nil {
		return err
	}

	srv, err := service.NewService(Flags.APIAddr.String(), DBServices, cipherService)
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
