package loans

import (
	"context"
	"github.com/arizalsaputro/billing-engine/internal/config"
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

func TestCalculateInstallmentsAdjustFinal(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		principal      decimal.Decimal
		interestRate   decimal.Decimal
		termWeeks      int64
		expectedTotal  decimal.Decimal
		expectedResult []decimal.Decimal
	}{
		{
			name:          "Exact division without residual",
			principal:     decimal.NewFromInt(10000000),
			interestRate:  decimal.NewFromFloat(0.10), // 10%
			termWeeks:     2,
			expectedTotal: decimal.NewFromInt(11000000), // 10M + 10% interest
			expectedResult: []decimal.Decimal{
				decimal.NewFromInt(5500000), // 11M / 2
				decimal.NewFromInt(5500000),
			},
		},
		{
			name:          "With residual rounding in final installment",
			principal:     decimal.NewFromInt(10000000),
			interestRate:  decimal.NewFromFloat(0.10), // 10%
			termWeeks:     3,
			expectedTotal: decimal.NewFromInt(11000000), // 10M + 10% interest
			expectedResult: []decimal.Decimal{
				decimal.NewFromInt(3666666), // Floor(11M / 3)
				decimal.NewFromInt(3666666),
				decimal.NewFromInt(3666668), // Adjust final installment
			},
		},
		{
			name:          "Single week term",
			principal:     decimal.NewFromInt(10000000),
			interestRate:  decimal.NewFromFloat(0.10), // 10%
			termWeeks:     1,
			expectedTotal: decimal.NewFromInt(11000000), // 10M + 10% interest
			expectedResult: []decimal.Decimal{
				decimal.NewFromInt(11000000), // Single installment matches total
			},
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installments, totalRepayment := calculateInstallmentsAdjustFinal(tt.principal, tt.interestRate, tt.termWeeks)

			// Check total repayment
			assert.Equal(t, tt.expectedTotal, totalRepayment, "Total repayment mismatch")

			// Check installments
			assert.Equal(t, tt.expectedResult, installments, "Installment breakdown mismatch")
		})
	}
}

func TestCreateLoan(t *testing.T) {
	// Mock configuration and dependencies
	mockLoanModel := new(mock2.LoanModelMock)
	mockConfig := config.Config{
		SettingInterestPerAnnum:  0.1, // 10%
		SettingLateFeePercentage: 2,   // 2%
		SettingGracePeriodDay:    7,   // 7 days
	}
	mockSvcCtx := &svc.ServiceContext{
		Config:    mockConfig,
		LoanModel: mockLoanModel,
	}

	// Create logic instance
	createLoanLogic := &CreateLoanLogic{
		svcCtx: mockSvcCtx,
		ctx:    context.Background(),
	}

	// Test case
	t.Run("successful loan creation", func(t *testing.T) {
		// Mock input
		req := &types.CreateLoanReq{
			PrincipalAmount: 10000000, // 10,000,000
			TermWeeks:       3,
		}

		// Expected loan and schedules
		expectedLoan := &model.Loans{
			PrincipalAmount:    decimal.NewFromInt(10000000),
			InterestRate:       decimal.NewFromFloat(0.1),
			TermWeeks:          3,
			OutstandingBalance: decimal.NewFromInt(11000000), // 10M + 10% interest
			Delinquent:         false,
			LateFeePercentage:  decimal.NewFromFloat(2),
			GracePeriodDays:    7,
		}

		//expectedSchedules := []*model.PaymentSchedule{
		//	{
		//		WeekNumber:      1,
		//		DueAmount:       decimal.NewFromInt(3666666), // Floor(11M / 3)
		//		DueDate:         time.Now().AddDate(0, 0, 7*1),
		//		GracePeriodDays: 7,
		//	},
		//	{
		//		WeekNumber:      2,
		//		DueAmount:       decimal.NewFromInt(3666666),
		//		DueDate:         time.Now().AddDate(0, 0, 7*2),
		//		GracePeriodDays: 7,
		//	},
		//	{
		//		WeekNumber:      3,
		//		DueAmount:       decimal.NewFromInt(3666668), // Adjust final installment
		//		DueDate:         time.Now().AddDate(0, 0, 7*3),
		//		GracePeriodDays: 7,
		//	},
		//}

		// Mock the CreateLoan call
		mockLoanModel.On("CreateLoan", mock.Anything, expectedLoan, mock.Anything).Return(int64(1), nil)

		// Call the function
		resp, err := createLoanLogic.CreateLoan(req)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.Base.Code)
		assert.Equal(t, "success", resp.Base.Msg)
		assert.Equal(t, int64(1), resp.LoanID)

		// Verify mock calls
		mockLoanModel.AssertCalled(t, "CreateLoan", mock.Anything, mock.MatchedBy(func(loan *model.Loans) bool {
			return loan.PrincipalAmount.Equal(expectedLoan.PrincipalAmount) &&
				loan.OutstandingBalance.Equal(expectedLoan.OutstandingBalance)
		}), mock.Anything)
	})
}
