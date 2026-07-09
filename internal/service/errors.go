package service

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrEmptyFields = errors.New("empty fields")
var ErrNoRowsAffected = errors.New("no rows affected")
var ErrUserExists = errors.New("user already exists")
var ErrInvalidInput = errors.New("invalid input")
var ErrNotFound = errors.New("not found")
