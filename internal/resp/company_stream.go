package resp

import "github.com/shopspring/decimal"

// CompanyStreamInfoResp 公司流水详情信息
type CompanyStreamInfoResp struct {
	Id          int             `json:"id"`           // 明细id
	Address     string          `json:"address"`      // 地址
	Uid         int             `json:"uid"`          // 用户id
	Name        string          `json:"name"`         // 用户帐号
	Email       string          `json:"email"`        // 邮件
	Income      decimal.Decimal `json:"income"`       // 入账
	Spend       decimal.Decimal `json:"spend"`        // 支出
	Balance     decimal.Decimal `json:"balance"`      // 现余额
	LastBalance decimal.Decimal `json:"last_balance"` // 之前余额
	Symbol      string          `json:"symbol"`       // 币种
	OrderId     string          `json:"order_id"`     // 订单号
	UpdatedAt   string          `json:"updated_at"`   // 时间
}

// CompanyStreamListResp 公司流水列表
type CompanyStreamListResp struct {
	Items []CompanyStreamInfoResp `json:"items"` // 帐号详情
	Page  BasePageResp            `json:"page"`  // 分页
}
