package resp

// WithdrawalAddrResp 提现地址列表集
type WithdrawalAddrResp struct {
	CurrencyName     string `json:"currency_name"`      // 货币名称
	WithdrawalAddrId int    `json:"withdrawal_addr_id"` // 地址id
	Name             string `json:"name"`               //  地址名称
	Address          string `json:"address"`            // 地址
	ChainCode        string `json:"chain_code"`         // 链code
	ChainName        string `json:"chain_name"`         // 链name
}

// WithdrawalDetailResp 提现详情信息
type WithdrawalDetailResp struct {
	Address   string `json:"address"`    // 地址
	Value     string `json:"value"`      // 金额
	Symbol    string `json:"symbol"`     // 代币代号
	Type      string `json:"type"`       // 标识 coin,token
	Status    int8   `json:"status"`     // 状态 0已提交,1待审核,2审核中,3通过,4不通过,5已完成,6取消,7提现失败
	Poundage  string `json:"poundage"`   // 手续费
	UpdatedAt string `json:"updated_at"` // 更新时间
}

// WithdrawalDetailListResp 提现详情列表
type WithdrawalDetailListResp struct {
	Items []WithdrawalDetailResp `json:"items"` // 帐号详情
	Page  BasePageResp           `json:"page"`  // 分页
}

// WithdrawalDetailResp 提现详情信息
type WithdrawalDetailAdminResp struct {
	Id              int    `json:"id"`               // 明细id
	Name            string `json:"name"`             // 帐号
	Email           string `json:"email"`            // 邮箱
	Uid             int    `json:"uid"`              // 用户id
	CoinId          int    `json:"coin_id"`          // 代币id
	CurrencyId      int    `json:"currency_id"`      // 货币id
	Type            string `json:"type"`             // 链类型
	Address         string `json:"address"`          // 地址
	Value           string `json:"value"`            // 金额
	Symbol          string `json:"symbol"`           // 代币代号
	Status          int8   `json:"status"`           // 状态 0已提交,1待审核,2审核中,3通过,4不通过,5已完成,6取消
	CustomerStatus  int8   `json:"customer_status"`  // 客服审核状态 0：审核中;1：审核通过;2：审核不通过
	FinancialStatus int8   `json:"financial_status"` // 财务审核状态 0：审核中;1：审核通过;2：审核不通过
	UpdatedAt       string `json:"updated_at"`       // 更新时间
	OrderId         string `json:"order_id"`         // 订单id
	Remark          string `json:"remark"`           // 备注
	AddressSource   int8   `json:"address_source"`   // 来源 0:未知 1:本站 2:外站
	FromAddress     string `json:"from_address"`     // 出金地址
	Balance         string `json:"balance"`          // 此时可用余额
	BlockCount      int    `json:"block_count"`      // 确认数
	CallbackStatus  string `json:"callback_status"`  // 回调状态码
	CallbackJson    string `json:"callback_json"`    // 回调json数据
}

// WithdrawalDetailAllListResp 提现详情列表
type WithdrawalDetailAllListResp struct {
	Items []WithdrawalDetailAdminResp `json:"items"` // 帐号详情
	Page  BasePageResp                `json:"page"`  // 分页
}

// AdminWithdrawalDetailResp  后台提现详情信息
type AdminWithdrawalDetailResp struct {
	Id              uint   `json:"id"`               // 用户的id
	OrderId         string `json:"order_id"`         // 订单id
	Symbol          string `json:"symbol"`           // 代币代号
	Status          int8   `json:"status"`           // 状态 0已提交,1待审核,2审核中,3通过,4不通过,5取消
	CustomerStatus  int8   `json:"customer_status"`  // 客服审核状态 0：审核中;1：审核通过;2：审核不通过
	FinancialStatus int8   `json:"financial_status"` // 财务审核状态 0：审核中;1：审核通过;2：审核不通过
	Address         string `json:"address"`          // 地址
	Value           string `json:"value"`            // 金额
	UpdatedAt       string `json:"updated_at"`       // 更新时间
}

// WithdrawalOrderTypeByAddrResp 根据type获取地址
type WithdrawalOrderTypeByAddrResp struct {
	Address string `json:"address"` // 地址
}
