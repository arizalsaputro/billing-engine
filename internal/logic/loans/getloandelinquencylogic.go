package loans

import (
	"context"
	"github.com/arizalsaputro/billing-engine/internal/model"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanDelinquencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoanDelinquencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanDelinquencyLogic {
	return &GetLoanDelinquencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanDelinquencyLogic) GetLoanDelinquency(req *types.GetLoanDelinquencygReq) (resp *types.GetLoanDelinquencygResp, err error) {
	if req.LoanID <= 0 {
		return nil, model.ErrNotFound
	}

	loan, err := l.svcCtx.LoanModel.GetLoanByID(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	resp = &types.GetLoanDelinquencygResp{
		Base: types.Base{
			Code: 200,
			Msg:  "ok",
		},
		LoanID:       loan.LoanId,
		IsDelinquent: loan.Delinquent,
	}
	return
}
