package loans

import (
	"context"
	"github.com/arizalsaputro/billing-engine/internal/model"
	"github.com/shopspring/decimal"
	"net/http"
	"time"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRepaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRepaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRepaymentLogic {
	return &CreateRepaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRepaymentLogic) CreateRepayment(req *types.CreateRepaymentReq) (resp *types.CreateRepaymentResp, err error) {
	paymentAmount := decimal.NewFromFloat(req.PaymentAmount)
	if paymentAmount.IsNegative() || paymentAmount.IsZero() {
		return nil, model.ErrInvalidPaymentAmount
	}

	if req.LoanID <= 0 {
		return nil, model.ErrNotFound
	}

	// Step 1: get loan, validate loan
	loan, err := l.svcCtx.LoanModel.GetLoanByID(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	// Step 1.2: validate loan
	if loan.OutstandingBalance.IntPart() == 0 {
		return nil, model.ErrLoanAlreadyPaid
	}

	if paymentAmount.GreaterThan(loan.OutstandingBalance) {
		return nil, model.ErrPaymentAmountMoreThanOutstanding
	}

	// Step 2: get unpaid repayment schedule
	schedules, err := l.svcCtx.LoanModel.GetRepaymentSchedules(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	// Step 2.1: validate no schedules
	if len(schedules) == 0 {
		return nil, model.ErrNoPaymentDueDate
	}

	var totalMissedPayment decimal.Decimal
	var repaymentAttempts []model.Payment
	for _, schedule := range schedules {
		totalMissedPayment = totalMissedPayment.Add(schedule.DueAmount)
		totalMissedPayment = totalMissedPayment.Add(schedule.LateFeeAmount)
		repaymentAttempts = append(repaymentAttempts, model.Payment{
			LoanID:        schedule.LoanID,
			PaymentAmount: schedule.DueAmount.Add(schedule.LateFeeAmount),
			PaymentDate:   time.Now().UTC(),
			WeekNumber:    schedule.WeekNumber,
			Status:        model.StatusPaymentPending,
		})
	}
	// step 2.2: validate payment amount with missed payment
	if !paymentAmount.Equal(totalMissedPayment) {
		return nil, model.ErrPaymentAmountNotMatchWithUnpaidWeeklyInstallment
	}

	// Step 3: create/insert table in repayment with status pending
	resp = &types.CreateRepaymentResp{
		Base: types.Base{
			Msg:  "your payment is on process",
			Code: http.StatusAccepted,
		},
		PaymentIDs: make([]*types.DataCreateRepayment, 0),
	}

	for _, repaymentAttempt := range repaymentAttempts {
		paymentId, err := l.svcCtx.LoanModel.UpsertPaymentWithID(l.ctx, repaymentAttempt)
		if err != nil {
			return nil, err
		}
		resp.PaymentIDs = append(resp.PaymentIDs, &types.DataCreateRepayment{
			PaymentId: paymentId,
		})
	}

	// Step 4: publish kafka
	// loop of resp.PaymentID the publish to kafka
	// TODO: for simplicity of code this step will be doing by hit API
	// try Hit /v1/billing/consume/pay using payment_id

	return resp, nil
}
