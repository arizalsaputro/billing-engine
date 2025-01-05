package model

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrNotFound          = sqlx.ErrNotFound
	ErrInvalidCreateLoan = errors.New("principal amount or term weeks must be greater than zero")
)
