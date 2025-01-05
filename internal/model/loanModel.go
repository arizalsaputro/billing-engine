package model

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strings"
	"time"
)

type (
	LoansModel interface {
		CreateLoan(ctx context.Context, loan *Loans, schedules []*PaymentSchedule) (int64, error)
		GetLoanByID(ctx context.Context, loanID int64) (*Loans, error)
	}

	Loans struct {
		LoanId             int64           `db:"loan_id"`
		PrincipalAmount    decimal.Decimal `db:"principal_amount"`
		InterestRate       decimal.Decimal `db:"interest_rate"`
		TermWeeks          int64           `db:"term_weeks"`
		WeeklyPayment      decimal.Decimal `db:"weekly_payment"`
		OutstandingBalance decimal.Decimal `db:"outstanding_balance"`
		Delinquent         bool            `db:"delinquent"`
		LateFeePercentage  decimal.Decimal `db:"late_fee_percentage"`
		GracePeriodDays    int64           `db:"grace_period_days"`
		CreatedAt          time.Time       `db:"created_at"`
		UpdatedAt          time.Time       `db:"updated_at"`
	}

	PaymentSchedule struct {
		ScheduleID      int64           `db:"schedule_id"`       // Primary Key
		LoanID          int64           `db:"loan_id"`           // Foreign Key to Loans table
		WeekNumber      int64           `db:"week_number"`       // Must be > 0
		DueAmount       decimal.Decimal `db:"due_amount"`        // NUMERIC(15, 2), Must be > 0
		DueDate         time.Time       `db:"due_date"`          // Due date of the payment
		Paid            bool            `db:"paid"`              // Default false
		PaymentDate     sql.NullTime    `db:"payment_date"`      // Nullable payment date
		LateFeeApplied  bool            `db:"late_fee_applied"`  // Default false
		LateFeeAmount   decimal.Decimal `db:"late_fee_amount"`   // NUMERIC(15, 2), Must be >= 0
		GracePeriodDays int64           `db:"grace_period_days"` // Default 3, Must be >= 0
		GracePeriodEnd  sql.NullTime    `db:"grace_period_end"`  // Generated always as (due_date + grace_period_days)
		CreatedAt       time.Time       `db:"created_at"`        // Timestamp with default CURRENT_TIMESTAMP
		UpdatedAt       time.Time       `db:"updated_at"`        // Timestamp with default CURRENT_TIMESTAMP
	}

	loansModel struct {
		conn sqlx.SqlConn
	}
)

func NewLoansModel(conn sqlx.SqlConn) LoansModel {
	return &loansModel{conn: conn}
}

func (m *loansModel) CreateLoan(ctx context.Context, loan *Loans, schedules []*PaymentSchedule) (int64, error) {
	var newLoanID int64
	err := m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// Step 1: Insert Loan
		insertLoanQuery := `
			INSERT INTO loan_schema.loans (principal_amount, interest_rate, term_weeks, weekly_payment, outstanding_balance, late_fee_percentage, grace_period_days)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING loan_id
		`

		err := session.QueryRowCtx(ctx, &newLoanID,
			insertLoanQuery,
			loan.PrincipalAmount.StringFixed(2),
			loan.InterestRate.StringFixed(2),
			loan.TermWeeks,
			loan.WeeklyPayment.StringFixed(2),
			loan.OutstandingBalance.StringFixed(2),
			loan.LateFeePercentage.StringFixed(2),
			loan.GracePeriodDays,
		)
		if err != nil {
			return err
		}

		// Step 2: Insert Payment Schedules
		// batch insert
		insertQuery := `INSERT INTO loan_schema.paymentschedules (loan_id, week_number, due_amount, due_date, grace_period_days) VALUES `
		var values []string
		var args []interface{}
		for i, schedule := range schedules {
			// Create a placeholder group (e.g., ($1, $2, $3, $4, $5))
			values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
			// Append the values to args
			args = append(args, newLoanID, schedule.WeekNumber, schedule.DueAmount, schedule.DueDate, schedule.GracePeriodDays)
		}
		finalQuery := insertQuery + strings.Join(values, ", ")
		_, err = session.ExecCtx(ctx, finalQuery, args...)
		if err != nil {
			return nil
		}

		return nil
	})

	return newLoanID, err
}

func (m *loansModel) GetLoanByID(ctx context.Context, loanID int64) (*Loans, error) {
	var loan Loans
	query := `
		SELECT loan_id, principal_amount, interest_rate, term_weeks, weekly_payment,
		       outstanding_balance, delinquent, late_fee_percentage, grace_period_days,
		       created_at, updated_at
		FROM loan_schema.loans
		WHERE loan_id = $1
	`
	err := m.conn.QueryRowCtx(ctx, &loan, query, loanID)
	if err != nil {
		return nil, err
	}

	return &loan, nil
}
