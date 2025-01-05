package consumer

import (
	"context"

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
	// todo: add your logic here and delete this line

	return
}
