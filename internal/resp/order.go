package resp

//后台查看全部用户明细
type RespUserTransInfoOrder struct {
	Id             int     `json:"id"`              // 订单id
	Name           string  `json:"name"`            // 用户名
	Uid            string  `json:"uid"`             // 用户id
	Email          string  `json:"email"`           // 邮箱
	Value          float64 `json:"value"`           // 目标金额
	Status         int8    `json:"status"`          // 状态 0确认中,1已确定
	UpdatedAt      string  `json:"updated_at"`      // 更新时间
	CurrencyId     int     `json:"currency_id"`     // 源币种id
	ExchangeId     int     `json:"exchange_id"`     // 目标币种id
	CurrencySymbol string  `json:"currency_symbol"` // 源币种标示
	ExchangeSymbol string  `json:"exchange_symbol"` // 目标币种标示
	Ratio          float64 `json:"ratio"`           // 比率

}

//后台查看全部用户，分页
type RespUserTransOrder struct {
	Items []RespUserTransInfoOrder `json:"items"` // 数据
	Page  BasePageResp             `json:"page"`  // 分页
}
