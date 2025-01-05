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

func CreateLoanHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateLoanReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, &types.Base{
				Code: http.StatusBadRequest,
				Msg:  err.Error(),
			})
			return
		}

		l := loans.NewCreateLoanLogic(r.Context(), svcCtx)
		resp, err := l.CreateLoan(&req)
		if err != nil {
			if errors.Is(err, model.ErrInvalidCreateLoan) {
				httpx.WriteJsonCtx(r.Context(), w, http.StatusBadRequest, &types.Base{
					Code: http.StatusBadRequest,
					Msg:  err.Error(),
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