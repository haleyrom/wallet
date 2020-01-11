package resp

// ReadCoinListResp 读取币列表结果
type ReadCoinListResp struct {
	CurrencyId        int     `json:"currency_id"`         // 代币id
	CurrencyName      string  `json:"currency_name"`       // 代币名称
	CurrencyDecimals  int     `json:"currency_decimals"`   // 代币小数点
	CoinId            int     `json:"coin_id"`             // 币种id
	Symbol            string  `json:"symbol"`              // 币种代号
	Name              string  `json:"name"`                // 币种名称
	BlockChainId      string  `json:"block_chain_id"`      // 区块链名称
	ChainCode         string  `json:"chain_code"`          // 链代号
	ChainName         string  `json:"chain_name"`          // 链名称
	Type              string  `json:"type"`                // 标识 coin,token
	ConfirmCount      int     `json:"confirm_count"`       // 充值入帐的区块链确认数
	MinDeposit        float64 `json:"min_deposit"`         // 最小充值金额，小于该金额不入账
	MinWithdrawal     float64 `json:"min_withdrawal"`      // 小于该金额不能提
	WithdrawalFee     float64 `json:"withdrawal_fee"`      // 提现手续费
	WithdrawalFeeType string  `json:"withdrawal_fee_type"` // 手续费类型 fixed 按百分百比,percent 固定收取
	ContractAddress   string  `json:"contract_address"`    // 合约地址:如该是type=token，这里必须输入
	UpdatedAt         string  `json:"updated_at"`          // 更新时间
	Status            int     `json:"status"`              // 状态 // 状态：0开启;1:停用;2:删除
	WithdrawalStatus  int8    `json:"withdrawal_status"`   // 充值状态：0开启;1:停用;
	DepositStatus     int8    `json:"deposit_status"`      // 提现状态：0开启;1:停用;
	CustomerStatus    int8    `json:"customer_status"`     // 客服状态:0 必须1：不必须
	FinancialStatus   int8    `json:"financial_status"`    // 财务状态:0 必须1：不必须
}

// ReadCoinInfoResp 代币信息
type ReadCoinInfoResp struct {
	CoinId            int     `json:"coin_id"`             // 代币id
	CurrencyId        int     `json:"currency_id"`         // 货币id
	Symbol            string  `json:"symbol"`              // 币种代号
	Name              string  `json:"name"`                // 币种名称
	BlockChainId      string  `json:"block_chain_id"`      // 区块链名称
	Type              string  `json:"type"`                // 标识 coin,token
	ConfirmCount      int     `json:"confirm_count"`       // 充值入帐的区块链确认数
	MinDeposit        float64 `json:"min_deposit"`         // 最小充值金额，小于该金额不入账
	MinWithdrawal     float64 `json:"min_withdrawal"`      // 小于该金额不能提
	WithdrawalFee     float64 `json:"withdrawal_fee"`      // 提现手续费
	WithdrawalFeeType string  `json:"withdrawal_fee_type"` // 手续费类型 fixed 按百分百比,percent 固定收取
	ContractAddress   string  `json:"contract_address"`    // 合约地址:如该是type=token，这里必须输入
	Abi               string  `json:"abi"`                 //
	Status            int     `json:"status"`              // 状态
	WithdrawalStatus  int8    `json:"withdrawal_status"`   // 充值状态：0开启;1:停用;
	DepositStatus     int8    `json:"deposit_status"`      // 提现状态：0开启;1:停用;
	CustomerStatus    int8    `json:"customer_status"`     // 客服状态:0 必须1：不必须
	FinancialStatus   int8    `json:"financial_status"`    // 财务状态:0 必须1：不必须
}

// ReadCoinDepositInfoResp 代币提现地址
type ReadCoinDepositInfoResp struct {
	Money             float64 `json:"money"`               // 可提现金额
	CurrencyId        uint    `json:"currency_id"`         // 币种id
	CoinId            uint    `json:"coin_id"`             // 代币id
	MinWithdrawal     float64 `json:"min_withdrawal"`      // 小于该金额不能提
	WithdrawalFee     float64 `json:"withdrawal_fee"`      // 提现手续费
	WithdrawalFeeType string  `json:"withdrawal_fee_type"` // 手续费类型 fixed 按百分百比,percent 固定收取
	Symbol            string  `json:"symbol"`              // 币种代号
	Type              string  `json:"type"`                // 标识 coin,token
	Status            int     `json:"status"`              // 状态
	WithdrawalStatus  int8    `json:"withdrawal_status"`   // 充值状态：0开启;1:停用;
	DepositStatus     int8    `json:"deposit_status"`      // 提现状态：0开启;1:停用;
	CustomerStatus    int8    `json:"customer_status"`     // 客服状态:0 必须1：不必须
	FinancialStatus   int8    `json:"financial_status"`    // 财务状态:0 必须1：不必须
}
