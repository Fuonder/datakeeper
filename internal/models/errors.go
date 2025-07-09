package models

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreationFailed = errors.New("user creation failed")
	ErrWrongCredentials   = errors.New("wrong credentials")
	ErrLoginAlreadyExists = errors.New("user with such login already exists")

	ErrNoData = errors.New("no data")
)
