package resp

// ReadCurrencyListResp 读取币种列表结果
type ReadCurrencyListResp struct {
	CurrencyId uint   `json:"currency_id"` // 币种id
	Symbol     string `json:"symbol"`      // 唯一标识
	Name       string `json:"name"`        // 名字
	Decimals   int    `json:"decimals"`    // 小数点
	UpdatedAt  string `json:"updated_at"`  // 更新时间
	Status     int8   `json:"status"`      // 状态：0开启;1:停用;2:删除
	Money      string `json:"money"`       // 可用金额
}

// ReadCurrencyTransferListResp 读取兑换币种列表结果
type ReadCurrencyTransferListResp struct {
	List     []ReadCurrencyListResp `json:"list"`     // 可兑换列表
	Transfer []ReadCurrencyListResp `json:"transfer"` // 兑换目标列表
}
