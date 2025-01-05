package cron

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ScheduleDelinquencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScheduleDelinquencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScheduleDelinquencyLogic {
	return &ScheduleDelinquencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScheduleDelinquencyLogic) ScheduleDelinquency(req *types.CronDelinquencyReq) (resp *types.CronDelinquencyResp, err error) {
	// todo: add your logic here and delete this line

	return
}
