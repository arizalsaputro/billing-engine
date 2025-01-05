package consumer

import (
	"context"
	"net/http"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeDelinquencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConsumeDelinquencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeDelinquencyLogic {
	return &ConsumeDelinquencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConsumeDelinquencyLogic) ConsumeDelinquency(req *types.ConsumeDelinquencyReq) (resp *types.ConsumeDelinquencyResp, err error) {
	err = l.svcCtx.LoanModel.UpdateLoanDelinquency(l.ctx, req.LoanID, true)
	if err != nil {
		return nil, err
	}
	resp = &types.ConsumeDelinquencyResp{}

	resp.Msg = "success"
	resp.Code = http.StatusOK

	return resp, nil
}
