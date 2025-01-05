package consumer

import (
	"context"
	"net/http"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeLateFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConsumeLateFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeLateFeeLogic {
	return &ConsumeLateFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConsumeLateFeeLogic) ConsumeLateFee(req *types.ConsumeLateFeeReq) (resp *types.ConsumeLateFeeResp, err error) {
	err = l.svcCtx.LoanModel.ApplyLateFees(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}
	resp = &types.ConsumeLateFeeResp{}
	resp.Msg = "success"
	resp.Code = http.StatusOK

	return resp, nil
}
