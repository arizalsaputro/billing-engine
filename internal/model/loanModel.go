package model

import (
	"context"
	"database/sql"
	"errors"
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
		GetPaymentByPaymentID(ctx context.Context, paymentID int64) (*Payment, error)
		GetRepaymentSchedules(ctx context.Context, loanID int64) ([]PaymentSchedule, error)
		UpsertPaymentWithID(ctx context.Context, payment Payment) (int64, error)
		ProcessRepayment(ctx context.Context, paymentID int64) (loanID int64, err error)
		GetDelinquentLoans(ctx context.Context, limit int) ([]Loans, error)
		RecheckLoanDelinquency(ctx context.Context, loanID int64) (int64, error)
		UpdateLoanDelinquency(ctx context.Context, loanID int64, isDelinquent bool) error
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

	Payment struct {
		PaymentID     int64           `db:"payment_id"`
		LoanID        int64           `db:"loan_id"`
		PaymentAmount decimal.Decimal `db:"payment_amount"`
		PaymentDate   time.Time       `db:"payment_date"`
		WeekNumber    int64           `db:"week_number"`
		Status        string          `db:"status"`
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

func (m *loansModel) GetPaymentByPaymentID(ctx context.Context, paymentID int64) (*Payment, error) {
	var payment Payment
	query := `
		SELECT payment_id, loan_id, payment_amount, payment_date, week_number, status
		FROM loan_schema.payments
		WHERE payment_id = $1
	`
	err := m.conn.QueryRowCtx(ctx, &payment, query, paymentID)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (m *loansModel) GetRepaymentSchedules(ctx context.Context, loanID int64) ([]PaymentSchedule, error) {
	var schedules []PaymentSchedule
	query := `
		SELECT schedule_id, loan_id, week_number, due_amount, due_date, paid, payment_date,
		       late_fee_applied, late_fee_amount, grace_period_days, grace_period_end,
		       created_at, updated_at
		FROM loan_schema.paymentschedules
		WHERE loan_id = $1 AND due_date <= CURRENT_DATE AND paid = FALSE
		ORDER BY week_number ASC
	`
	err := m.conn.QueryRowsCtx(ctx, &schedules, query, loanID)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (m *loansModel) UpsertPaymentWithID(ctx context.Context, payment Payment) (int64, error) {
	var paymentID int64
	query := `
		INSERT INTO loan_schema.payments (loan_id, week_number, payment_amount, status)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (loan_id, week_number)
		DO NOTHING
		RETURNING payment_id;
	`

	err := m.conn.QueryRowCtx(ctx, &paymentID,
		query, payment.LoanID, payment.WeekNumber, payment.PaymentAmount.StringFixed(2), payment.PaymentDate, payment.Status)
	if err != nil {
		return 0, err
	}

	// If no row was inserted and no `payment_id` is returned, fetch the existing `payment_id`
	if paymentID == 0 {
		existingQuery := `
			SELECT payment_id
			FROM loan_schema.payments
			WHERE loan_id = $1 AND week_number = $2
		`
		err = m.conn.QueryRowCtx(ctx, &paymentID, existingQuery, payment.LoanID, payment.WeekNumber)
		if err != nil {
			return 0, err
		}
	}

	return paymentID, nil
}

func (m *loansModel) ProcessRepayment(ctx context.Context, paymentID int64) (loanID int64, err error) {
	// Begin Transaction
	// rollback dan commit di wrapping di function ini
	err = m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		// Step 1. Select and Lock the table payment
		var payment Payment
		query := `
			SELECT loan_id, payment_amount, status, week_number
			FROM loan_schema.payments
			WHERE payment_id = $1
			FOR UPDATE
		`
		err := m.conn.QueryRowCtx(ctx, &payment, query, paymentID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrPaymentNotFound
			}
			return fmt.Errorf("failed to fetch payment: %w", err)
		}

		// Check idempotency
		if payment.Status != StatusPaymentPending {
			return nil // Already processed or in a non-pending state, exit early
		}

		loanID = payment.LoanID

		// Step 2: Update Repayment Schedule
		updateScheduleQuery := `
			UPDATE loan_schema.paymentschedules
			SET paid = TRUE, payment_date = CURRENT_DATE
			WHERE loan_id = $1 AND week_number = $2 AND paid = FALSE
		`
		_, err = m.conn.ExecCtx(ctx, updateScheduleQuery, payment.LoanID, payment.WeekNumber)
		if err != nil {
			return fmt.Errorf("failed to update repayment schedule: %w", err)
		}

		// Step 3: Update Loan Outstanding Balance; atomic update
		updateLoanQuery := `
			UPDATE loan_schema.loans
			SET outstanding_balance = outstanding_balance - $1
			WHERE loan_id = $2
		`
		_, err = m.conn.ExecCtx(ctx, updateLoanQuery, payment.PaymentAmount.StringFixed(2), payment.LoanID)
		if err != nil {
			return fmt.Errorf("failed to update loan balance: %w", err)
		}

		// Step 4: Mark Payment as Success
		updatePaymentQuery := `
			UPDATE loan_schema.payments
			SET status = 'success'
			WHERE payment_id = $1
		`
		_, err = m.conn.ExecCtx(ctx, updatePaymentQuery, paymentID)
		if err != nil {
			return fmt.Errorf("failed to update payment status: %w", err)
		}

		return nil
	})

	return loanID, err
}

func (m *loansModel) GetDelinquentLoans(ctx context.Context, limit int) ([]Loans, error) {
	var loans []Loans

	query := `
		SELECT loan_id
		FROM loan_schema.paymentschedules
		WHERE due_date < CURRENT_DATE AND paid = FALSE
		GROUP BY loan_id
		HAVING COUNT(*) >= 2
		LIMIT $1
	`

	err := m.conn.QueryRowsCtx(ctx, &loans, query, limit)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

func (m *loansModel) RecheckLoanDelinquency(ctx context.Context, loanID int64) (int64, error) {
	// Step 1: Check for delinquency condition
	var delinquentCount int64
	query := `
		SELECT COUNT(*)
		FROM loan_schema.paymentschedules
		WHERE loan_id = $1 AND due_date < CURRENT_DATE AND paid = FALSE
	`

	err := m.conn.QueryRowCtx(ctx, &delinquentCount, query, loanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return delinquentCount, nil // No repayment schedules, loan is not delinquent
		}
		return delinquentCount, err
	}

	return delinquentCount, nil
}

func (m *loansModel) UpdateLoanDelinquency(ctx context.Context, loanID int64, isDelinquent bool) error {
	updateQuery := `
		UPDATE loan_schema.loans
		SET delinquent = $1
		WHERE loan_id = $2
	`

	_, err := m.conn.ExecCtx(ctx, updateQuery, isDelinquent, loanID)
	if err != nil {
		return err
	}

	return nil
}
