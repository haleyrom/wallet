package params

// ReadCoinListParam 读取币列表
type ReadCoinListParam BaseBindParam

// UpdateCoinParam 更新参数
type UpdateCoinParam struct {
	Base              *BaseParam `json:"claims" form:"claims"`
	CurrencyId        uint       `json:"currency_id" form:"currency_id"  binding:"required"`
	Id                uint       `json:"coin_id" form:"coin_id" binding:"required"`                         //  coin id
	Symbol            string     `json:"symbol" form:"symbol" binding:"required"`                           // 币种代号
	Name              string     `json:"name" form:"name" binding:"required"`                               // 币种名称
	BlockChainId      uint       `json:"block_chain_id" form:"block_chain_id" binding:"required"`           // 区块链名称
	Type              string     `json:"type" form:"type" binding:"required"`                               // 标识 coin,token
	ConfirmCount      int        `json:"confirm_count" form:"confirm_count" binding:"required"`             // 充值入帐的区块链确认数
	MinDeposit        float64    `json:"min_deposit" form:"min_deposit" binding:"required"`                 // 最小充值金额，小于该金额不入账
	MinWithdrawal     float64    `json:"min_withdrawal" form:"min_withdrawal" binding:"required"`           // 小于该金额不能提
	WithdrawalFee     float64    `json:"withdrawal_fee" form:"withdrawal_fee" binding:"required"`           // 提现手续费
	WithdrawalFeeType string     `json:"withdrawal_fee_type" form:"withdrawal_fee_type" binding:"required"` // 手续费类型 fixed 按百分百比,percent 固定收取
	ContractAddress   string     `json:"contract_address" form:"contract_address" binding:"required"`       // 合约地址:如该是type=token，这里必须输入
	Abi               string     `json:"abi" form:"abi"`                                                    // 字节数
	WithdrawalStatus  int8       `json:"withdrawal_status" form:"withdrawal_status"`                        // 状态：0开启;1:停用;
	DepositStatus     int8       `json:"deposit_status" form:"deposit_status"`                              // 状态：0开启;1:停用;
	CustomerStatus    int8       `json:"customer_status" form:"customer_status"`                            // 客服状态:0 必须1：不必须
	FinancialStatus   int8       `json:"financial_status" form:"financial_status"`                          // 财务状态:0 必须1：不必须
}

// AddCoinParam 更新参数
type AddCoinParam struct {
	Base              *BaseParam `json:"claims" form:"claims"`
	CurrencyId        uint       `json:"currency_id" form:"currency_id"  binding:"required"`
	Symbol            string     `json:"symbol" form:"symbol" binding:"required"`                           // 币种代号
	Name              string     `json:"name" form:"name" binding:"required"`                               // 币种名称
	BlockChainId      uint       `json:"block_chain_id" form:"block_chain_id" binding:"required"`           // 区块链名称
	Type              string     `json:"type" form:"type" binding:"required"`                               // 标识 coin,token
	ConfirmCount      int        `json:"confirm_count" form:"confirm_count" binding:"required"`             // 充值入帐的区块链确认数
	MinDeposit        float64    `json:"min_deposit" form:"min_deposit" binding:"required"`                 // 最小充值金额，小于该金额不入账
	MinWithdrawal     float64    `json:"min_withdrawal" form:"min_withdrawal" binding:"required"`           // 小于该金额不能提
	WithdrawalFee     float64    `json:"withdrawal_fee" form:"withdrawal_fee" binding:"required"`           // 提现手续费
	WithdrawalFeeType string     `json:"withdrawal_fee_type" form:"withdrawal_fee_type" binding:"required"` // 手续费类型 fixed 按百分百比,percent 固定收取
	ContractAddress   string     `json:"contract_address" form:"contract_address" binding:"required"`       // 合约地址:如该是type=token，这里必须输入
	Abi               string     `json:"abi" form:"abi" binding:"required"`                                 // 字节数
	WithdrawalStatus  int8       `json:"withdrawal_status" form:"withdrawal_status"`                        // 状态：0开启;1:停用;
	DepositStatus     int8       `json:"deposit_status" form:"deposit_status"`                              // 状态：0开启;1:停用;
	CustomerStatus    int8       `json:"customer_status" form:"customer_status"`                            // 客服状态:0 必须1：不必须
	FinancialStatus   int8       `json:"financial_status" form:"financial_status"`                          // 财务状态:0 必须1：不必须
}

// RemoveCoinParam 删除参数
type RemoveCoinParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	CoinId uint       `json:"coin_id" form:"coin_id" binding:"required"`
}

// ReadCoinInfoParam 读取参数
type ReadCoinInfoParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	CoinId uint       `json:"coin_id" form:"coin_id" binding:"required"`
}

// UpdateCoinStatusParam 更新代币状态参数
type UpdateCoinStatusParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	CoinId uint       `json:"coin_id" form:"coin_id" binding:"required"`
	Status int8       `json:"status" form:"status"`
}

// ReadCoinDepositInfoParam  读取代币提现信息参数
type ReadCoinDepositInfoParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	CoinId uint       `json:"coin_id" form:"coin_id" binding:"required"`
}

// ReadListAccountCoinParam  读取代币提现信息参数
type ReadListAccountCoinParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id"  binding:"required"`
}
