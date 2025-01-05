package consumer

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
