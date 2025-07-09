package main

import (
	"fmt"
	"github.com/Fuonder/datakeeper.git/internal/client/cliservice"
	"github.com/Fuonder/datakeeper.git/internal/client/ui"
	"github.com/Fuonder/datakeeper.git/internal/config"
	"github.com/Fuonder/datakeeper.git/internal/logger"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
	"log"
	"time"
)

func main() {
	flags, err := config.ParseClientFlags()
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

func run(flags config.ClientOptions) error {
	srv, err := cliservice.NewService([]byte(flags.AESKey), flags.APIAddr.String(), 5*time.Second)
	if err != nil {
		return err
	}

	p := tea.NewProgram(ui.NewAuthModel(srv))
	_, err = p.Run()
	if err != nil {
		return err
	}
	return nil
}
