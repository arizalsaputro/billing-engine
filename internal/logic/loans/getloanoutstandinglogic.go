package loans

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanOutstandingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoanOutstandingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanOutstandingLogic {
	return &GetLoanOutstandingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanOutstandingLogic) GetLoanOutstanding(req *types.GetLoanOutstandingReq) (resp *types.GetLoanOutStandingResp, err error) {
	loan, err := l.svcCtx.LoanModel.GetLoanByID(l.ctx, int64(req.LoanID))
	if err != nil {
		return nil, err
	}
	resp = &types.GetLoanOutStandingResp{
		LoanID: int(loan.LoanId),
		OutstandingBalance: ,
	}

	return
}
