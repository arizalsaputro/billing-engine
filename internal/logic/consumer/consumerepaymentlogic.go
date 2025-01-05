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
	// Step 1: Process Payment
	loanId, err := l.svcCtx.LoanModel.ProcessRepayment(l.ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}

	// Step 2: do some other business logic: auto check user delinquency
	//TODO:
	// kalau dalah desain, di prosess ini akan publish kafka, demi simplicity kita abaikan dulu ganti ke rest
	// Hit v1//billing/consume/check/delinquency

	resp = &types.ConsumeRepaymentResp{
		Base: types.Base{
			Msg:  "done",
			Code: 200,
		},
		LoanID: loanId,
	}

	return resp, nil
}
