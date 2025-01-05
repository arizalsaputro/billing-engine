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
	loans, err := l.svcCtx.LoanModel.GetDelinquentLoans(l.ctx, req.QueryLimit)
	if err != nil {
		return nil, err
	}

	// TODO:
	// kalau dari system yang dibuat sih harusnya ini publish ke kafka ya, tapi biar simple
	// buat simulate coba hit v1/billing/delinquency

	resp = &types.CronDelinquencyResp{
		Data: make([]*types.DataDelinquency, 0),
	}

	for _, loan := range loans {
		resp.Data = append(resp.Data, &types.DataDelinquency{
			LoanID: loan.LoanId,
		})
	}

	return resp, nil
}
