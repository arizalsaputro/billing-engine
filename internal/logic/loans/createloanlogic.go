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

type CreateLoanLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLoanLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLoanLogic {
	return &CreateLoanLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func calculateInstallmentsAdjustFinal(principal decimal.Decimal, interestRate decimal.Decimal, termWeeks int64) ([]decimal.Decimal, decimal.Decimal) {
	// Calculate total repayment
	totalRepayment := principal.Mul(decimal.NewFromInt(1).Add(interestRate))

	// Calculate raw weekly installment
	rawInstallment := totalRepayment.Div(decimal.NewFromInt(termWeeks))

	// Round down each installment
	roundedInstallment := rawInstallment.Floor()

	// Calculate the total for all but the last installment
	totalRounded := roundedInstallment.Mul(decimal.NewFromInt(termWeeks - 1))

	// Calculate the final installment
	finalInstallment := totalRepayment.Sub(totalRounded)

	// Generate installment list
	installments := make([]decimal.Decimal, termWeeks)
	for i := int64(0); i < termWeeks-1; i++ {
		installments[i] = roundedInstallment
	}
	installments[termWeeks-1] = finalInstallment

	return installments, totalRepayment
}

func (l *CreateLoanLogic) CreateLoan(req *types.CreateLoanReq) (resp *types.CreateLoanResp, err error) {
	if req.PrincipalAmount <= 0 || req.TermWeeks <= 0 {
		return nil, model.ErrInvalidCreateLoan
	}
	// Step 1: calculate installment
	principalAmount := decimal.NewFromInt(req.PrincipalAmount) // Principal: 10,000,000
	interestRate := decimal.NewFromFloat(l.svcCtx.Config.SettingInterestPerAnnum)
	installments, totalRepayment := calculateInstallmentsAdjustFinal(principalAmount, interestRate, req.TermWeeks)

	// Step 2: Prepare Loan Struct
	loan := &model.Loans{
		PrincipalAmount:    principalAmount,
		InterestRate:       interestRate,
		TermWeeks:          req.TermWeeks,
		OutstandingBalance: totalRepayment,
		Delinquent:         false,
		LateFeePercentage:  decimal.NewFromFloat(l.svcCtx.Config.SettingLateFeePercentage), // Example: 2%
		GracePeriodDays:    int64(l.svcCtx.Config.SettingGracePeriodDay),                   // Example: 7-day grace period
	}

	// Step 3: Prepare Loan Schedule
	var schedules []*model.PaymentSchedule
	startDate := time.Now()
	for i := 1; i <= int(req.TermWeeks); i++ {
		dueDate := startDate.AddDate(0, 0, 7*i)
		schedules = append(schedules, &model.PaymentSchedule{
			WeekNumber:      int64(i),
			DueAmount:       installments[i-1],
			DueDate:         dueDate,
			GracePeriodDays: int64(l.svcCtx.Config.SettingGracePeriodDay),
		})
	}

	newLoanID, err := l.svcCtx.LoanModel.CreateLoan(l.ctx, loan, schedules)
	if err != nil {
		return nil, err
	}

	resp = &types.CreateLoanResp{
		Base: types.Base{
			Code: http.StatusOK,
			Msg:  "success",
		},
		LoanID: newLoanID,
	}

	return resp, nil
}
