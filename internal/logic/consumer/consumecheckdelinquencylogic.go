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
	// todo: add your logic here and delete this line

	return
}
