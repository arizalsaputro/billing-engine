package mock

import (
	"context"
	"github.com/arizalsaputro/billing-engine/internal/model"
	"github.com/stretchr/testify/mock"
)

type LoanModelMock struct {
	mock.Mock
}

func (m *LoanModelMock) GetLoanByID(ctx context.Context, loanID int64) (*model.Loans, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) GetPaymentByPaymentID(ctx context.Context, paymentID int64) (*model.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) GetRepaymentSchedules(ctx context.Context, loanID int64) ([]model.PaymentSchedule, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) UpsertPaymentWithID(ctx context.Context, payment model.Payment) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) ProcessRepayment(ctx context.Context, paymentID int64) (loanID int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) GetDelinquentLoans(ctx context.Context, limit int) ([]model.SimplePaymentSchedule, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) RecheckLoanDelinquency(ctx context.Context, loanID int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) UpdateLoanDelinquency(ctx context.Context, loanID int64, isDelinquent bool) error {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) GetLateRepaymentSchedules(ctx context.Context, limit int) ([]model.SimplePaymentSchedule, error) {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) ApplyLateFees(ctx context.Context, loanID int64) error {
	//TODO implement me
	panic("implement me")
}

func (m *LoanModelMock) CreateLoan(ctx context.Context, loan *model.Loans, schedules []*model.PaymentSchedule) (int64, error) {
	args := m.Called(ctx, loan, schedules)
	return args.Get(0).(int64), args.Error(1)
}
