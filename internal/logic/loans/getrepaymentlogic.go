package loans

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRepaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRepaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRepaymentLogic {
	return &GetRepaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRepaymentLogic) GetRepayment(req *types.GetRepaymentReq) (resp *types.GetRepaymentResp, err error) {
	payment, err := l.svcCtx.LoanModel.GetPaymentByPaymentID(l.ctx, req.PaymentID)
	if err != nil {
		return nil, err
	}

	resp = &types.GetRepaymentResp{
		Base: types.Base{
			Code: 200,
			Msg:  "ok",
		},
		PaymentID:     payment.PaymentID,
		PaymentAmount: payment.PaymentAmount.IntPart(),
		PaymentDate:   payment.PaymentDate.String(),
		Status:        payment.Status,
	}

	return resp, nil
}
