package resp

import "github.com/shopspring/decimal"

//后台查看全部用户明细
type UserTransInfoOrderResp struct {
	Id             int             `json:"id"`              // 订单id
	Name           string          `json:"name"`            // 用户名
	Uid            string          `json:"uid"`             // 用户id
	Email          string          `json:"email"`           // 邮箱
	Value          decimal.Decimal `json:"value"`           // 目标金额
	Status         int8            `json:"status"`          // 状态 0确认中,1已确定
	UpdatedAt      string          `json:"updated_at"`      // 更新时间
	CurrencyId     int             `json:"currency_id"`     // 源币种id
	ExchangeId     int             `json:"exchange_id"`     // 目标币种id
	CurrencySymbol string          `json:"currency_symbol"` // 源币种标示
	ExchangeSymbol string          `json:"exchange_symbol"` // 目标币种标示
	Ratio          string          `json:"ratio"`           // 比率

}

//后台查看全部用户，分页
type UserTransOrderResp struct {
	Items []UserTransInfoOrderResp `json:"items"` // 数据
	Page  BasePageResp             `json:"page"`  // 分页
}

// AccountTransferInfoResp 钱包转账信息
type AccountTransferInfoResp struct {
	OrderId      string          `json:"order_id"`      // 订单id
	Uid          string          `json:"uid"`           // 用户id
	UserName     string          `json:"user_name"`     // 用户名
	UserEmail    string          `json:"user_email"`    // 用户邮箱
	AdverseId    string          `json:"adverse_id"`    // 对方id
	AdverseName  string          `json:"adverse_name"`  // 对方名
	AdverseEmail string          `json:"adverse_email"` // 对方邮箱
	Status       int             `json:"status"`        // 状态 0未完成 1完成
	Balance      decimal.Decimal `json:"balance"`       // 金额
	Symbol       string          `json:"symbol"`        // 币种
	UpdatedAt    string          `json:"updated_at"`    // 更新时间
}

// AccountTransferListResp 后台查看转正列表，分页
type AccountTransferListResp struct {
	Items []AccountTransferInfoResp `json:"items"` // 数据
	Page  BasePageResp              `json:"page"`  // 分页
}
