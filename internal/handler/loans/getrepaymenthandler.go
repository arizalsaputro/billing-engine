package loans

import (
	"net/http"

	"github.com/arizalsaputro/billing-engine/internal/logic/loans"
	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRepaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRepaymentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loans.NewGetRepaymentLogic(r.Context(), svcCtx)
		resp, err := l.GetRepayment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
