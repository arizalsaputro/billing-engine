package consumer

import (
	"net/http"

	"github.com/arizalsaputro/billing-engine/internal/logic/consumer"
	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ConsumeRepaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConsumeRepaymentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := consumer.NewConsumeRepaymentLogic(r.Context(), svcCtx)
		resp, err := l.ConsumeRepayment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
