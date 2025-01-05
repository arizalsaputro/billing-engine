package model

import (
	"errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	StatusPaymentSuccess = "success"
	StatusPaymentPending = "pending"
)

var (
	ErrNotFound                                         = sqlx.ErrNotFound
	ErrInvalidCreateLoan                                = errors.New("principal amount or term weeks must be greater than zero")
	ErrInvalidPaymentAmount                             = errors.New("payment amount must be greater than zero")
	ErrPaymentAmountMoreThanOutstanding                 = errors.New("payment amount more than outstanding loan")
	ErrPaymentAmountNotMatchWithUnpaidWeeklyInstallment = errors.New("payment amount not match with pas due unpaid")
	ErrLoanAlreadyPaid                                  = errors.New("loan already paid")
	ErrNoPaymentDueDate                                 = errors.New("no payment due date")
	ErrPaymentNotFound                                  = errors.New("payment not found")
)
