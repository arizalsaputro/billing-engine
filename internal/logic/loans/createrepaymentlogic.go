package loans

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRepaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRepaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRepaymentLogic {
	return &CreateRepaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRepaymentLogic) CreateRepayment(req *types.CreateRepaymentReq) (resp *types.CreateRepaymentResp, err error) {
	// todo: add your logic here and delete this line

	return
}
