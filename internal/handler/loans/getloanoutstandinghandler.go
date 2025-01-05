package loans

import (
	"errors"
	"github.com/arizalsaputro/billing-engine/internal/model"
	"net/http"

	"github.com/arizalsaputro/billing-engine/internal/logic/loans"
	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetLoanOutstandingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLoanOutstandingReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loans.NewGetLoanOutstandingLogic(r.Context(), svcCtx)
		resp, err := l.GetLoanOutstanding(&req)
		if err != nil {
			if errors.Is(err, model.ErrNotFound) {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, &types.Base{
					Code: http.StatusBadRequest,
					Msg:  "Loan not found",
				})
			} else {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, &types.Base{
					Code: http.StatusInternalServerError,
					Msg:  err.Error(),
				})
			}
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
