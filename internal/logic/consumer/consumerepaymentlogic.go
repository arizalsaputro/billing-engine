package consumer

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeRepaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConsumeRepaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumeRepaymentLogic {
	return &ConsumeRepaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConsumeRepaymentLogic) ConsumeRepayment(req *types.ConsumeRepaymentReq) (resp *types.ConsumeRepaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
