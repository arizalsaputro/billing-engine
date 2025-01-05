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
	// todo: add your logic here and delete this line

	return
}
