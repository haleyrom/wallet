package resp

// CurrencyQuoteInfoResp 读取币种汇率列表结果
type CurrencyQuoteInfoResp struct {
	Id            uint    `json:"id"`             // id
	Code          string  `json:"code"`           // 标示
	BaseCurrency  string  `json:"base_currency"`  // 基础货币
	QuoteCurrency string  `json:"quote_currency"` // 基础货币
	Price         float64 `json:"price"`          // 金额
	UpdatedAt     string  `json:"updated_at"`     // 更新时间
}

// CurrencyQuoteListResp 读取币种汇率列表列表
type CurrencyQuoteListResp struct {
	Items []CurrencyQuoteInfoResp `json:"items"` // 帐号详情
	Page  BasePageResp            `json:"page"`  // 分页
}
