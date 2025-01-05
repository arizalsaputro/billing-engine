package model

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"testing"
)

func TestGetLoanByID(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewSqlConnFromDB(db)
	loanModel := NewLoansModel(sqlxDB)

	// Define test data
	loanID := int64(1)
	expectedLoan := &Loans{
		LoanId:             loanID,
		PrincipalAmount:    decimal.NewFromInt(100000),
		InterestRate:       decimal.NewFromFloat(5.5),
		TermWeeks:          52,
		OutstandingBalance: decimal.NewFromInt(50000),
		Delinquent:         false,
		LateFeePercentage:  decimal.NewFromFloat(2.5),
		GracePeriodDays:    7,
		//CreatedAt:          sql.NullTime{Valid: true},
		//UpdatedAt:          sql.NullTime{Valid: true},
	}

	// Set up the mock expectation
	mock.ExpectQuery(`SELECT loan_id, principal_amount, interest_rate, term_weeks, outstanding_balance, delinquent, late_fee_percentage, grace_period_days, created_at, updated_at FROM loan_schema\.loans WHERE loan_id = \$1`).
		WithArgs(loanID).
		WillReturnRows(sqlmock.NewRows([]string{
			"loan_id", "principal_amount", "interest_rate", "term_weeks", "outstanding_balance", "delinquent", "late_fee_percentage", "grace_period_days", "created_at", "updated_at",
		}).AddRow(
			expectedLoan.LoanId,
			expectedLoan.PrincipalAmount,
			expectedLoan.InterestRate,
			expectedLoan.TermWeeks,
			expectedLoan.OutstandingBalance,
			expectedLoan.Delinquent,
			expectedLoan.LateFeePercentage,
			expectedLoan.GracePeriodDays,
			expectedLoan.CreatedAt,
			expectedLoan.UpdatedAt,
		))

	// Call the function
	ctx := context.Background()
	actualLoan, err := loanModel.GetLoanByID(ctx, loanID)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, actualLoan)
	assert.Equal(t, expectedLoan, actualLoan)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}

func TestCreateLoan(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewSqlConnFromDB(db)
	loanModel := NewLoansModel(sqlxDB)

	// Define test data
	loan := &Loans{
		PrincipalAmount:    decimal.NewFromInt(100000),
		InterestRate:       decimal.NewFromFloat(5.5),
		TermWeeks:          52,
		OutstandingBalance: decimal.NewFromInt(100000),
		LateFeePercentage:  decimal.NewFromFloat(2.5),
		GracePeriodDays:    7,
	}

	schedules := []*PaymentSchedule{
		{
			WeekNumber: 1,
			DueAmount:  decimal.NewFromInt(2000),
			//DueDate:         sqlmock.AnyArg(),
			GracePeriodDays: 7,
		},
		{
			WeekNumber: 2,
			DueAmount:  decimal.NewFromInt(2000),
			//DueDate:         sqlmock.AnyArg(),
			GracePeriodDays: 7,
		},
	}

	newLoanID := int64(1)

	// Step 1: Mock Insert Loan
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO loan_schema\.loans \(principal_amount, interest_rate, term_weeks, outstanding_balance, late_fee_percentage, grace_period_days\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) RETURNING loan_id`).
		WithArgs(
			loan.PrincipalAmount.StringFixed(2),
			loan.InterestRate.StringFixed(2),
			loan.TermWeeks,
			loan.OutstandingBalance.StringFixed(2),
			loan.LateFeePercentage.StringFixed(2),
			loan.GracePeriodDays,
		).
		WillReturnRows(sqlmock.NewRows([]string{"loan_id"}).AddRow(newLoanID))

	// Step 2: Mock Insert Payment Schedules
	var values []string
	for i := range schedules {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
	}
	mock.ExpectExec(`INSERT INTO loan_schema\.paymentschedules \(loan_id, week_number, due_amount, due_date, grace_period_days\) VALUES `+strings.Join(values, ", ")).
		WithArgs(newLoanID, schedules[0].WeekNumber, schedules[0].DueAmount, schedules[0].DueDate, schedules[0].GracePeriodDays,
			newLoanID, schedules[1].WeekNumber, schedules[1].DueAmount, schedules[1].DueDate, schedules[1].GracePeriodDays).
		WillReturnResult(sqlmock.NewResult(2, 2))

	mock.ExpectCommit()

	// Call the function
	ctx := context.Background()
	actualLoanID, err := loanModel.CreateLoan(ctx, loan, schedules)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, newLoanID, actualLoanID)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet sqlmock expectations: %v", err)
	}
}
