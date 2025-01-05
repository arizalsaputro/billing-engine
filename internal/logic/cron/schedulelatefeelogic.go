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
	schedules, err := l.svcCtx.LoanModel.GetLateRepaymentSchedules(l.ctx, req.QueryLimit)
	if err != nil {
		return nil, err
	}

	resp = &types.CronLateFeeResp{
		Data: make([]*types.DataLoanScheduleLate, 0),
	}

	for _, schedule := range schedules {
		resp.Data = append(resp.Data, &types.DataLoanScheduleLate{
			LoanID: schedule.LoanID,
		})

		// TODO: kalau dari design harusnya ini publish kafka
		// demi simplicity sampai sini aja
		// untuk simulate coba /v1/billing/consume/late
	}

	return resp, nil
}
