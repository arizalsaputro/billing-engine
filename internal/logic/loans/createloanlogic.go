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

func (l *CreateLoanLogic) CreateLoan(req *types.CreateLoanReq) (resp *types.CreateLoanResp, err error) {
	if req.PrincipalAmount <= 0 || req.TermWeeks <= 0 {
		return nil, model.ErrInvalidCreateLoan
	}

	// Step 1: Calculate Loan Details
	interestRate := decimal.NewFromFloat(l.svcCtx.Config.SettingInterestPerAnnum) // e.g., 10%

	// Calculate totalRepayment: PrincipalAmount * (100 + interestRate) / 100
	principalAmount := decimal.NewFromFloat(float64(req.PrincipalAmount))
	hundred := decimal.NewFromInt(100)
	totalRepayment := principalAmount.Mul(hundred.Add(interestRate)).Div(hundred)

	// Calculate weeklyPayment: totalRepayment / TermWeeks
	termWeeks := decimal.NewFromInt(int64(req.TermWeeks))
	weeklyPayment := totalRepayment.Div(termWeeks)

	// Step 2: Prepare Loan Struct
	loan := &model.Loans{
		PrincipalAmount:    principalAmount,
		InterestRate:       interestRate,
		TermWeeks:          int64(req.TermWeeks),
		WeeklyPayment:      weeklyPayment,
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
			DueAmount:       weeklyPayment,
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
