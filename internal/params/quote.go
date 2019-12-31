package params

// CreateQuoteParam 创建汇率参数解析
type CreateQuoteParam struct {
	BaseCurrency  string  `json:"base_currency" form:"base_currency" binding:"required"`   // 基础货币
	QuoteCurrency string  `json:"quote_currency" form:"quote_currency" binding:"required"` // 报价货币
	Price         float64 `json:"price" form:"price" binding:"required"`                   // 金额
}

// UpdateQuoteParam 更新汇率参数解析
type UpdateQuoteParam struct {
	Id    uint    `json:"id" form:"id" binding:"required"`       // 汇率id
	Price float64 `json:"price" form:"price" binding:"required"` // 金额
}

// ReadQuotePageParam 汇率分页
type ReadQuotePageParam struct {
	Base      *BaseParam `json:"claims" form:"claims" `
	StartTime int        `json:"start_time" form:"start_time"`
	EndTime   int        `json:"end_time" form:"end_time"`
	Keyword   string     `json:"keyword" form:"keyword"`
	Page      int        `json:"page" form:"page"  binding:"required"`
	PageSize  int        `json:"pageSize" form:"pageSize" binding:"required"`
}
