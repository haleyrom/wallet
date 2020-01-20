package params

// AccountInfoParam 详情参数
type AccountInfoParam BaseBindParam

// AccountTFORInfoParam TFOR钱包信息解析参数
type AccountTFORInfoParam BaseBindParam

// AccountTransferParam 账本转参数
type AccountTransferParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	Symbol string     `json:"symbol" form:"symbol"  binding:"required"`
	Money  float64    `json:"money" form:"money"  binding:"required,gt=0"`
}

// AccountChangeParam 账本转账本参数
type AccountChangeParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id"  binding:"required"`
	Money      string     `json:"money" form:"money"  binding:"required,gt=0"`
	ChangeId   uint       `json:"change_id" form:"change_id" binding:"required"`
	Ratio      string     `json:"ratio" form:"ratio" binding:"required"`
}

// AccountDetailParam 账单明细
type AccountDetailParam struct {
	Base     *BaseParam `json:"claims" form:"claims"`
	Page     int        `json:"page" form:"page"  binding:"required"`
	PageSize int        `json:"pageSize" form:"pageSize" binding:"required"`
}

// AccountShareBonusParam 节点分红
type AccountShareBonusParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	Money  float64    `json:"money" form:"money"  binding:"required,gt=0"`
	Symbol string     `json:"symbol" form:"symbol"  binding:"required"`
}

// AccountWithdrawalParam 钱包提现
type AccountWithdrawalParam struct {
	Base             *BaseParam `json:"claims" form:"claims"`
	Password         string     `json:"password" form:"password" binding:"required"`
	Money            float64    `json:"money" form:"money"  binding:"required,gt=0"`
	CurrencyId       uint       `json:"currency_id" form:"currency_id"  binding:"required"`
	CoinId           uint       `json:"coin_id" form:"coin_id"  binding:"required"`
	WithdrawalAddrId uint       `json:"withdrawal_addr_id" form:"withdrawal_addr_id" binding:"required"`
}

// AccountListParam 钱包列表
type AccountListParam BaseListParam

// AccountUserListParam 用户钱包列表
type AccountUserListParam BaseListParam

// AccountWithdrawalListParam 用户提现钱包列表
type AccountWithdrawalListParam BaseListParam

// AccountOrderListParam  用户订单列表
type AccountOrderListParam BaseListParam

// AccountDepositDetailList 用户充值流水列表
type AccountDepositDetailListParam BaseListParam

// CompanyDepositList 公司充值流水
type CompanyDepositListParam BaseListParam

// CompanyWithdrawalListParam 公司提币流水
type CompanyWithdrawalListParam BaseListParam

// CompanyDepositAddrListParam 公司充值地址
type CompanyDepositAddrListParam BaseListParam

// CompanyWithdrawalAddrListParam 公司提币地址
type CompanyWithdrawalAddrListParam BaseListParam

// AccountWithdrawalDetailParam 用户充值明细
type AccountWithdrawalDetailParam struct {
	Base *BaseParam `json:"claims" form:"claims"`
	Id   uint       `json:"id" form:"id" binding:"required"`
}

// AccountWithdrawalDetailCustomerParam 客服参数解析
type AccountWithdrawalDetailCustomerParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	Id         uint       `json:"id" form:"id" binding:"required"`
	CustomerId uint       `json:"customer_id" form:"customer_id"`
	Status     int8       `json:"status" form:"status"`
}

// AccountWithdrawalDetailFinancialParam 财务参数解析
type AccountWithdrawalDetailFinancialParam struct {
	Base        *BaseParam `json:"claims" form:"claims"`
	Id          uint       `json:"id" form:"id" binding:"required"`
	FinancialId uint       `json:"financial_id" form:"financial_id"`
	Status      int8       `json:"status" form:"status"`
}

// JoinCompanyWithdrawalAddrParam 新增公司提币地址参数解析
type JoinCompanyWithdrawalAddrParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id" binding:"required"`
	Address    string     `json:"address" form:"address" binding:"required"`
}

// JoinCompanyDepositAddrParam 新增公司充值地址参数解析
type JoinCompanyDepositAddrParam struct {
	Base       *BaseParam `json:"claims" form:"claims"`
	CurrencyId uint       `json:"currency_id" form:"currency_id" binding:"required"`
	Address    string     `json:"address" form:"address" binding:"required"`
}

// UpdateCompanyAddrParam 更新公司地址参数解析
type UpdateCompanyAddrParam struct {
	Base    *BaseParam `json:"claims" form:"claims"`
	Id      uint       `json:"id" form:"id" binding:"required"`
	Address string     `json:"address" form:"address" binding:"required"`
}

// UpdateCompanyAddrStatusParam 更新公司地址状态参数解析
type UpdateCompanyAddrStatusParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	Id     uint       `json:"id" form:"id" binding:"required"`
	Status int8       `json:"status" form:"status"`
}

// CreateCompanyAddrParam 创建地址参数解析
type CreateCompanyAddrParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	Symbol string     `json:"symbol" form:"symbol" binding:"required"`
}

// JoinRechargeParam 充值参数解析
type JoinRechargeParam struct {
	Uid    string  `json:"uid" form:"uid" binding:"required"`
	Symbol string  `json:"symbol" form:"symbol" binding:"required"`
	Money  float64 `json:"money" form:"money" binding:"required"`
}

// ReadRechargePageParam 充值列表参数解析
type ReadRechargePageParam BaseListParam

// RemoveRechargeParam 删除充值记录
type RemoveRechargeParam struct {
	Id uint `json:"id" form:"id" binding:"required"`
}

// AudioRechargeParam 审核充值
type AudioRechargeParam struct {
	Id          uint `json:"id" form:"id" binding:"required"`
	FinancialId uint `json:"financial_id" form:"financial_id"`
	Status      int8 `json:"status" form:"status"`
}

// AccountCurrencyDetailParam 账单币种明细
type AccountCurrencyDetailParam struct {
	Base      *BaseParam `json:"claims" form:"claims"`
	AccountId uint       `json:"account_id" form:"account_id" binding:"required"`
	Page      int        `json:"page" form:"page"  binding:"required"`
	PageSize  int        `json:"pageSize" form:"pageSize" binding:"required"`
	Type      string     `json:"type" form:"type" binding:"required"`
}

// AccountPersonTransferParam 个人转账参数
type AccountPersonTransferParam struct {
	Base        *BaseParam `json:"claims" form:"claims"`
	Email       string     `json:"email" form:"email"  binding:"email"`
	Symbol      string     `json:"symbol" form:"symbol"  binding:"required"`
	Money       float64    `json:"money" form:"money"  binding:"required,gt=0"`
	PayPassword string     `json:"pay_password" form:"pay_password" binding:"required"`
}

// AccountTFORListParam TFOR钱包列表解析参数
type AccountTFORListParam struct {
	Base *BaseParam `json:"claims" form:"claims"`
	Uids string     `json:"uids" form:"uids"`
}

// AccountChargeDetailParam 个人收款账单明细
type AccountChargeDetailParam struct {
	Base     *BaseParam `json:"claims" form:"claims"`
	Page     int        `json:"page" form:"page"  binding:"required"`
	PageSize int        `json:"pageSize" form:"pageSize" binding:"required"`
}

// AccountTransferListParam 转账列表解析参数
type AccountTransferListParam BaseListParam
