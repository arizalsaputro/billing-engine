package loans

import (
	"context"

	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoanDelinquencyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLoanDelinquencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoanDelinquencyLogic {
	return &GetLoanDelinquencyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLoanDelinquencyLogic) GetLoanDelinquency(req *types.GetLoanDelinquencygReq) (resp *types.GetLoanDelinquencygResp, err error) {
	// todo: add your logic here and delete this line

	return
}
