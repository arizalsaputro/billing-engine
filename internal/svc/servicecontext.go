package svc

import (
	"github.com/arizalsaputro/billing-engine/internal/config"
	"github.com/arizalsaputro/billing-engine/internal/model"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	LoanModel model.LoansModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// connect db
	conn := sqlx.NewSqlConn("postgres", c.DataSource)

	return &ServiceContext{
		Config:    c,
		LoanModel: model.NewLoansModel(conn),
	}
}
