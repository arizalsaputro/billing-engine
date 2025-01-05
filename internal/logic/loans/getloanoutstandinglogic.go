package loans

import (
	"context"
	"github.com/arizalsaputro/billing-engine/internal/model"

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
	if req.LoanID <= 0 {
		return nil, model.ErrNotFound
	}

	loan, err := l.svcCtx.LoanModel.GetLoanByID(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	resp = &types.GetLoanOutStandingResp{
		Base: types.Base{
			Code: 200,
			Msg:  "ok",
		},
		LoanID:             loan.LoanId,
		OutstandingBalance: loan.OutstandingBalance.IntPart(),
	}

	return
}
