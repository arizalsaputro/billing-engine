package cron

import (
	"net/http"

	"github.com/arizalsaputro/billing-engine/internal/logic/cron"
	"github.com/arizalsaputro/billing-engine/internal/svc"
	"github.com/arizalsaputro/billing-engine/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ScheduleDelinquencyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CronDelinquencyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := cron.NewScheduleDelinquencyLogic(r.Context(), svcCtx)
		resp, err := l.ScheduleDelinquency(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
