// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.5

package types

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ConsumeCheckDelinquencyReq struct {
	LoanID int `json:"loanId"`
}

type ConsumeCheckDelinquencyResp struct {
	Base
}

type ConsumeDelinquencyReq struct {
	LoanID int `json:"loanId"`
}

type ConsumeDelinquencyResp struct {
	Base
}

type ConsumeLateFeeReq struct {
	ScheduleID int `json:"scheduleId"`
}

type ConsumeLateFeeResp struct {
	Base
}

type ConsumeRepaymentReq struct {
	PaymentID int `json:"paymentId"`
}

type ConsumeRepaymentResp struct {
	Base
}

type CreateLoanReq struct {
	PrincipalAmount int `json:"principalAmount"`
	TermWeeks       int `json:"termWeeks"`
}

type CreateLoanResp struct {
	Base
	LoanID int `json:"loanId"`
}

type CreateRepaymentReq struct {
	LoanID        int `path:"loanId"`
	PaymentAmount int `json:"paymentAmount"`
}

type CreateRepaymentResp struct {
	Base
	PaymentID int `json:"paymentId"`
}

type CronDelinquencyReq struct {
	QueryLimit int `form:"queryLimit"`
}

type CronDelinquencyResp struct {
	Base
	Data []*DataDelinquency `json:"data"`
}

type CronLateFeeReq struct {
	QueryLimit int `form:"queryLimit"`
}

type CronLateFeeResp struct {
	Base
	Data []*DataLoanScheduleLate `json:"data"`
}

type DataDelinquency struct {
	LoanID int `json:"loanId"`
}

type DataLoanScheduleLate struct {
	ScheduleID     int    `json:"scheduleId"`
	DueAmount      int    `json:"dueAmount"`
	GracePeriodEnd string `json:"gracePeriodEnd"`
}

type GetLoanDelinquencygReq struct {
	LoanID int `path:"loanId"`
}

type GetLoanDelinquencygResp struct {
	Base
	LoanID       int  `json:"loanId"`
	IsDelinquent bool `json:"isDelinquent"`
}

type GetLoanOutStandingResp struct {
	Base
	LoanID             int `json:"loanId"`
	OutstandingBalance int `json:"outstandingBalance"`
}

type GetLoanOutstandingReq struct {
	LoanID int `path:"loanId"`
}

type GetRepaymentReq struct {
	PaymentID int `path:"paymentId"`
}

type GetRepaymentResp struct {
	Base
	PaymentID     int    `json:"paymentId"`
	PaymentAmount int    `json:"paymentAmount"`
	PaymentDate   string `json:"paymentDate"`
	Status        string `json:"status"`
}
