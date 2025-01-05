package consumer

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeCheckDelinquencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConsumeCheckDelinquencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeCheckDelinquencyLogic {
	return &ConsumeCheckDelinquencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConsumeCheckDelinquencyLogic) ConsumeCheckDelinquency(req *types.ConsumeCheckDelinquencyReq) (resp *types.ConsumeCheckDelinquencyResp, err error) {
	count, err := l.svcCtx.LoanModel.RecheckLoanDelinquency(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.LoanModel.UpdateLoanDelinquency(l.ctx, req.LoanID, count >= 2)
	if err != nil {
		return nil, err
	}

	resp = &types.ConsumeCheckDelinquencyResp{
		Base: types.Base{
			Code: 200,
			Msg:  "success",
		},
	}

	return resp, nil
}
