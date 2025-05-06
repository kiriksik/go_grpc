package storage

import "errors"

var (
	ErrUserExists   = errors.New("this user already exists")
	ErrUserNotFound = errors.New("user not exists")
	ErrAppNotFound  = errors.New("app not exists")
)
