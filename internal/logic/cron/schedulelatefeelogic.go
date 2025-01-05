package cron

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScheduleLateFeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScheduleLateFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleLateFeeLogic {
	return &ScheduleLateFeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleLateFeeLogic) ScheduleLateFee(req *types.CronLateFeeReq) (resp *types.CronLateFeeResp, err error) {
	// todo: add your logic here and delete this line

	return
}
