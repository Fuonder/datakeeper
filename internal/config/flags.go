package config

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

type NetAddress struct {
	IPAddr string
	Port   int
}

func (n *NetAddress) String() string {
	return fmt.Sprintf("%s:%d", n.IPAddr, n.Port)
}

func (n *NetAddress) Set(value string) error {
	value = strings.TrimPrefix(value, "http://")
	values := strings.Split(value, ":")
	if len(values) != 2 {
		return fmt.Errorf("%w: \"%s\"", ErrNotFullIP, value)
	}
	n.IPAddr = values[0]
	if n.IPAddr == "" {
		return fmt.Errorf("%w: \"%s\"", ErrInvalidIP, values[0])
	}
	var err error
	n.Port, err = strconv.Atoi(values[1])
	if err != nil {
		return fmt.Errorf("%w: \"%s\"", ErrInvalidPort, values[1])
	}
	return nil
}

type CommonOptions struct {
	APIAddr  NetAddress
	LogLevel string
	HashKey  string
	AESKey   string
}

type ServerOptions struct {
	CommonOptions
	DatabaseDSN string
}

type ClientOptions = CommonOptions

func (f CommonOptions) String() string {
	return fmt.Sprintf(
		"API address: %s, Log level: %s, Hash key: %s, AES key: %s",
		f.APIAddr.String(),
		f.LogLevel,
		f.HashKey,
		f.AESKey,
	)
}

func (s ServerOptions) String() string {
	return fmt.Sprintf(
		"%s, DB DSN: %s",
		s.CommonOptions.String(),
		s.DatabaseDSN,
	)
}

func (f *CommonOptions) bindFlags() {
	flag.StringVar(&f.LogLevel, "l", "debug", "log level (debug, info, warn, error, fatal, panic)")
	flag.Var(&f.APIAddr, "a", "ip and port of server in format <ip>:<port>")
	flag.StringVar(&f.HashKey, "h", "TEST123", "hash key")
	flag.StringVar(&f.AESKey, "c", "32-byte-long-encryption-key-1234", "aes256 32 byte long key")
}

func (f *CommonOptions) parseEnv() error {
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		if err := f.APIAddr.Set(envRunAddr); err != nil {
			return err
		}
	}
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		f.LogLevel = envLogLevel
	}
	if envHashSecret := os.Getenv("HASH_SECRET"); envHashSecret != "" {
		f.HashKey = envHashSecret
	}
	if envAESSecret := os.Getenv("AES_SECRET"); envAESSecret != "" {
		f.AESKey = envAESSecret
	}
	return nil
}

func ParseClientFlags() (ClientOptions, error) {
	flags := ClientOptions{
		APIAddr: NetAddress{
			IPAddr: "localhost",
			Port:   8080,
		},
		LogLevel: "debug",
		HashKey:  "TEST123",
		AESKey:   "32-byte-long-encryption-key-1234",
	}
	flags.bindFlags()
	flag.Parse()
	if err := flags.parseEnv(); err != nil {
		return flags, err
	}
	return flags, nil
}

func ParseServerFlags() (ServerOptions, error) {
	flags := ServerOptions{
		CommonOptions: CommonOptions{
			APIAddr: NetAddress{
				IPAddr: "localhost",
				Port:   8080,
			},
			LogLevel: "debug",
			HashKey:  "TEST123",
			AESKey:   "32-byte-long-encryption-key-1234",
		},
		DatabaseDSN: "postgres://datakeeper:12345678@localhost:5432/datakeeper?sslmode=disable",
	}
	flags.CommonOptions.bindFlags()
	flag.StringVar(&flags.DatabaseDSN, "d", flags.DatabaseDSN, "Database DSN")
	flag.Parse()

	if err := flags.CommonOptions.parseEnv(); err != nil {
		return flags, err
	}
	if envDSN := os.Getenv("DATABASE_URI"); envDSN != "" {
		flags.DatabaseDSN = envDSN
	}

	return flags, nil
}
