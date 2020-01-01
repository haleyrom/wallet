package params

// ReadWithdrawalAddrListParam 读取提现地址列表
type ReadWithdrawalAddrListParam struct {
	Base         *BaseParam `json:"claims" form:"claims"`
	CurrencyId   uint       `json:"currency_id" form:"currency_id"`
	BlockChainId uint       `json:"block_chain_id" form:"block_chain_id"`
}

// CreateWithdrawalAddrParam 创建提现地址
type CreateWithdrawalAddrParam struct {
	Base         *BaseParam `json:"claims" form:"claims" `
	Name         string     `json:"name" form:"name"`
	CurrencyId   uint       `json:"currency_id" form:"currency_id" binding:"required"`
	BlockChainId uint       `json:"block_chain_id" form:"block_chain_id" binding:"required"`
	Address      string     `json:"address" form:"address" binding:"required"`
}

// UpdateWithdrawalAddrParam 更新提现地址
type UpdateWithdrawalAddrParam struct {
	Base             *BaseParam `json:"claims" form:"claims" `
	Name             string     `json:"name" form:"name"`
	WithdrawalAddrId uint       `json:"withdrawal_addr_id" form:"withdrawal_addr_id" binding:"required"`
	Address          string     `json:"address" form:"address" binding:"required"`
}

// RemoveWithdrawalAddrParam  删除提现地址
type RemoveWithdrawalAddrParam struct {
	Base             *BaseParam `json:"claims" form:"claims" `
	WithdrawalAddrId uint       `json:"withdrawal_addr_id" form:"withdrawal_addr_id" binding:"required"`
}

// ReadWithdrawalDetailParam 读取提现明细
type ReadWithdrawalDetailParam struct {
	Base     *BaseParam `json:"claims" form:"claims" `
	Page     int        `json:"page" form:"page"  binding:"required"`
	PageSize int        `json:"pageSize" form:"pageSize" binding:"required"`
}

// WithdrawalCallbackParam 提现回调参数解析
type WithdrawalCallbackParam struct {
	AppId           string `json:"app_id" form:"app_id"`
	OrderId         string `json:"order_id" form:"order_id"`
	TransactionHash string `json:"transaction_hash" form:"transaction_hash"`
	BlockCount      string `json:"block_count" form:"block_count"`
	BlockNumber     string `json:"block_number" form:"block_number"`
	FromAddress     string `json:"from_address" form:"from_address"`
	ToAddress       string `json:"to_address" form:"to_address"`
	Symbol          string `json:"symbol" form:"symbol"`
	ContractAddress string `json:"contract_address" form:"contract_address"`
	Value           string `json:"value" form:"value"`
	Code            string `json:"code" form:"code"`
	Message         string `json:"message" form:"message"`
	Hash            string `json:"hash" form:"hash"`
}
