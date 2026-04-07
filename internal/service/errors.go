package service

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUnauthorized = errors.New("unauthorized")
var ErrEmptyFields = errors.New("empty fields")
var ErrNoRows = errors.New("no rows affected")
