package loans

import (
	"context"
	mock2 "github.com/arizalsaputro/billing-engine/internal/logic/mock"
	"github.com/arizalsaputro/billing-engine/internal/model"
	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestCreateRepayment(t *testing.T) {
	// Mock dependencies
	mockLoanModel := new(mock2.LoanModelMock)
	mockSvcCtx := &svc.ServiceContext{
		LoanModel: mockLoanModel,
	}

	// Create logic instance
	createRepaymentLogic := &CreateRepaymentLogic{
		svcCtx: mockSvcCtx,
		ctx:    context.Background(),
	}

	// Test case: Successful repayment creation
	t.Run("successful repayment creation", func(t *testing.T) {
		// Mock input
		req := &types.CreateRepaymentReq{
			LoanID:        1,
			PaymentAmount: 200000.00,
		}

		// Mock loan data
		loan := &model.Loans{
			LoanId:             1,
			OutstandingBalance: decimal.NewFromInt(200000),
		}
		mockLoanModel.On("GetLoanByID", mock.Anything, req.LoanID).Return(loan, nil)

		// Mock schedules
		schedules := []model.PaymentSchedule{
			{
				LoanID:        1,
				DueAmount:     decimal.NewFromInt(150000),
				LateFeeAmount: decimal.NewFromInt(50000),
				WeekNumber:    1,
			},
		}
		mockLoanModel.On("GetRepaymentSchedules", mock.Anything, req.LoanID).Return(schedules, nil)

		// Mock payment upsert
		mockLoanModel.On("UpsertPaymentWithID", mock.Anything, mock.Anything).Return(int64(1001), nil)

		// Call the function
		resp, err := createRepaymentLogic.CreateRepayment(req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusAccepted, resp.Base.Code)
		assert.Equal(t, "your payment is on process", resp.Base.Msg)
		assert.Len(t, resp.PaymentIDs, 1)
		assert.Equal(t, int64(1001), resp.PaymentIDs[0].PaymentId)

		// Verify mock calls
		mockLoanModel.AssertCalled(t, "GetLoanByID", mock.Anything, req.LoanID)
		mockLoanModel.AssertCalled(t, "GetRepaymentSchedules", mock.Anything, req.LoanID)
		mockLoanModel.AssertCalled(t, "UpsertPaymentWithID", mock.Anything, mock.Anything)
	})

	// Test case: Invalid payment amount
	t.Run("invalid payment amount", func(t *testing.T) {
		req := &types.CreateRepaymentReq{
			LoanID:        1,
			PaymentAmount: -100.00, // Invalid negative amount
		}

		resp, err := createRepaymentLogic.CreateRepayment(req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, model.ErrInvalidPaymentAmount, err)
	})

	// Test case: Loan already paid
	t.Run("loan already paid", func(t *testing.T) {
		req := &types.CreateRepaymentReq{
			LoanID:        1,
			PaymentAmount: 200000.00,
		}

		loan := &model.Loans{
			LoanId:             1,
			OutstandingBalance: decimal.Zero, // No outstanding balance
		}
		mockLoanModel.On("GetLoanByID", mock.Anything, req.LoanID).Return(loan, nil)

		resp, err := createRepaymentLogic.CreateRepayment(req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, model.ErrLoanAlreadyPaid, err)
	})

	// Test case: Payment amount does not match missed payment
	t.Run("payment amount mismatch", func(t *testing.T) {
		req := &types.CreateRepaymentReq{
			LoanID:        1,
			PaymentAmount: 300000.00, // More than missed payment
		}

		loan := &model.Loans{
			LoanId:             1,
			OutstandingBalance: decimal.NewFromInt(200000),
		}
		mockLoanModel.On("GetLoanByID", mock.Anything, req.LoanID).Return(loan, nil)

		schedules := []model.PaymentSchedule{
			{
				LoanID:        1,
				DueAmount:     decimal.NewFromInt(150000),
				LateFeeAmount: decimal.NewFromInt(50000),
				WeekNumber:    1,
			},
		}
		mockLoanModel.On("GetRepaymentSchedules", mock.Anything, req.LoanID).Return(schedules, nil)

		resp, err := createRepaymentLogic.CreateRepayment(req)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, model.ErrPaymentAmountNotMatchWithUnpaidWeeklyInstallment, err)
	})
}
