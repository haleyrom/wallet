package params

// ReadCurrencyListParam 币种列表参数
type ReadCurrencyListParam BaseBindParam

// ReadCurrencyTransferListParam 币种兑换列表参数
type ReadCurrencyTransferListParam BaseBindParam

// UpdateCurrencyParam 更新币种参数
type UpdateCurrencyParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id" binding:"required"`
	Symbol     string     `json:"symbol" form:"symbol" binding:"required"`
	Name       string     `json:"name"  form:"name" binding:"required"`
	Decimals   int        `json:"decimals" form:"decimals" binding:"required"`
}

// AddCurrencyParam 创建币种参数
type AddCurrencyParam struct {
	Base     *BaseParam `json:"claims" form:"claims"`
	Symbol   string     `json:"symbol" form:"symbol" binding:"required"`
	Name     string     `json:"name"  form:"name" binding:"required"`
	Decimals int        `json:"decimals" form:"decimals" binding:"required"`
}

// UpdateCurrencyStatus 更新币种状态参数
type UpdateCurrencyStatusParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id" binding:"required"`
	Status     int8       `json:"status" form:"status"`
}

// RmCurrencyParam 删除币种参数
type RmCurrencyParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id" binding:"required"`
}

// CurrencyQuoteParam 货币兑换参数解析
type CurrencyQuoteParam struct {
	Base          *BaseParam `json:"claims" form:"claims"`
	QuoteCurrency string     `json:"quote_currency" form:"quote_currency" binding:"required"`
}
