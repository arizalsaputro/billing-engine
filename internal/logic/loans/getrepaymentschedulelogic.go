package loans

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRepaymentScheduleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRepaymentScheduleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepaymentScheduleLogic {
	return &GetRepaymentScheduleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRepaymentScheduleLogic) GetRepaymentSchedule(req *types.GetRepaymentScheduleReq) (resp *types.GetRepaymentScheduleResp, err error) {
	schedules, err := l.svcCtx.LoanModel.GetRepaymentSchedules(l.ctx, req.LoanID)
	if err != nil {
		return nil, err
	}

	resp = &types.GetRepaymentScheduleResp{
		Data: make([]*types.RepaymentSchedule, 0),
	}

	for _, schedule := range schedules {
		resp.Data = append(resp.Data, &types.RepaymentSchedule{
			WeekNumber: int(schedule.WeekNumber),
			DueAmount:  schedule.DueAmount.Add(schedule.LateFeeAmount).InexactFloat64(),
			DueDate:    schedule.DueDate.Format("2006-01-02"),
			IsPaid:     schedule.Paid,
		})
	}

	return resp, nil
}
