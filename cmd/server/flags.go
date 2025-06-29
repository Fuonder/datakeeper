package server

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrNotFullIP   = errors.New("given ip address and port incorrect")
	ErrInvalidIP   = errors.New("incorrect ip address")
	ErrInvalidPort = errors.New("incorrect port number")
)

type netAddress struct {
	ipaddr string
	port   int
}

func (n *netAddress) String() string {
	return fmt.Sprintf("%s:%d", n.ipaddr, n.port)
}
func (n *netAddress) Set(value string) error {
	value = strings.TrimPrefix(value, "http://")
	values := strings.Split(value, ":")
	if len(values) != 2 {
		return fmt.Errorf("%w: \"%s\"", ErrNotFullIP, value)
	}
	n.ipaddr = values[0]
	if n.ipaddr == "" {
		return fmt.Errorf("%w: \"%s\"", ErrInvalidIP, values[0])
	}
	var err error
	n.port, err = strconv.Atoi(values[1])
	if err != nil {
		return fmt.Errorf("%w: \"%s\"", ErrInvalidPort, values[1])
	}
	return nil
}

type CliOptions struct {
	APIAddr     netAddress
	LogLevel    string
	DatabaseDSN string
}

func (f *CliOptions) String() string {
	return fmt.Sprintf("APIAddress: %s, "+
		"DatabaseDSN: %s, "+
		"LogLevel: %s",
		f.APIAddr.String(),
		f.DatabaseDSN,
		f.LogLevel,
	)
}

var (
	Flags = CliOptions{
		APIAddr: netAddress{
			ipaddr: "localhost",
			port:   8080,
		},
		LogLevel:    "debug",
		DatabaseDSN: "postgres://datakeeper:12345678@localhost:5432/datakeeper?sslmode=disable",
	}
)

func parseFlags() error {
	flag.StringVar(&Flags.LogLevel, "l", "debug", "log level (debug, info, warn, error, fatal, panic)")
	flag.StringVar(&Flags.DatabaseDSN, "d", "postgres://datakeeper:12345678@localhost:5432/datakeeper?sslmode=disable", "Database DSN")
	flag.Var(&Flags.APIAddr, "a", "ip and port of server in format <ip>:<port>")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		err := Flags.APIAddr.Set(envRunAddr)
		if err != nil {
			return err
		}
	}
	if envDatabaseDSN := os.Getenv("DATABASE_URI"); envDatabaseDSN != "" {
		Flags.DatabaseDSN = envDatabaseDSN
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		Flags.LogLevel = envLogLevel
	}

	return nil
}
