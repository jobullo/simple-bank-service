package database

import (
	"errors"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrParentNotFound = errors.New("parent not found")
	ErrInvalidType    = errors.New("invalid transaction type")
)
