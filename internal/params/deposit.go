package params

// DepositParam  充值参数
type DepositParam struct {
	AppId   string `json:"app_id"`
	OrderId string `json:"order_id"`
	Hash    string `json:"hash"`
}

// ReadDepositAddListParam 列表参数
type ReadDepositAddListParam BaseBindParam

// JoinDepositDetailParam 写入充值记录参数
type JoinDepositDetailParam struct {
	Base            *BaseParam `json:"claims" form:"claims"`
	Address         string     `json:"address" form:"address" binding:"required"`
	Value           float64    `json:"value" form:"value" binding:"required"`
	BlockNumber     int        `json:"block_number" form:"block_number" binding:"required"`
	BlockCount      int        `json:"block_count" form:"block_count" binding:"required"`
	TransactionHash string     `json:"transaction_hash" form:"transaction_hash" binding:"required"`
	CoinId          uint       `json:"coin_id" form:"coin_id" binding:"required"`
}

// ReadWithdrawalAddrParam 读取提现地址
type ReadDepositAddrParam struct {
	Base         *BaseParam `json:"claims" form:"claims" `
	BlockChainId uint       `json:"block_chain_id" form:"block_chain_id" binding:"required"`
}

// TopUpDepositParam 充值参数解析
type TopUpDepositParam struct {
	Base            *BaseParam `json:"claims" form:"claims" `
	Address         string     `json:"address" form:"address" binding:"required"`
	Money           string     `json:"money" form:"money" binding:"required"`
	BlockNumber     string     `json:"block_number" form:"block_number" `
	BlockCount      string     `json:"block_count" form:"block_count"`
	TransactionHash string     `json:"transaction_hash" form:"transaction_hash" binding:"required"`
	Symbol          string     `json:"symbol" form:"symbol" binding:"required"`
	Type            string     `json:"type" form:"type" binding:"required"`
	Hash            string     `json:"hash" form:"hash" binding:"required"`
	ContractAddress string     `json:"contract_address" form:"contract_address"`
}

// ReadDepositDetailParam 读取充值明细
type ReadDepositDetailParam struct {
	Base     *BaseParam `json:"claims" form:"claims" `
	Page     int        `json:"page" form:"page"  binding:"required"`
	PageSize int        `json:"pageSize" form:"pageSize" binding:"required"`
}
