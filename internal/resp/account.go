package resp

const (
	// AccountDetailUp 明细充值
	AccountDetailUp int8 = iota + 1
	// AccountDetailOut 提币
	AccountDetailOut
	// AccountDetailShare 节点分成
	AccountDetailShare
	// AccountDetailRelease 算力释放
	AccountDetailRelease
	// AccountDetailInto 转入
	AccountDetailInto
	// AccountDetailUpgrade 升级
	AccountDetailUpgrade
	// AccountDetailConvert 兑换
	AccountDetailConvert
	// AccountDetailTransfer 转账
	AccountDetailTransfer
)

// AccountInfoResp 详情返回结果集
type AccountInfoResp struct {
	Uid            int    `json:"uid"`             // 用户uid
	AccountId      string `json:"account_id"`      // 钱包id
	CurrencyId     int    `json:"currency_id"`     // 币种id
	Balance        string `json:"balance"`         // 余额
	Symbol         string `json:"symbol"`          // 币种标识
	Decimals       int    `json:"decimals"`        // 小数点位数
	Name           string `json:"name"`            // 名称
	UpdatedAt      string `json:"updated_at"`      // 时间
	BlockedBalance string `json:"blocked_balance"` // 冻结金额
	BlockChainId   uint   `json:"block_chain_id"`  // 链id
}

// AccountDetailResp 帐号详情
type AccountDetailResp struct {
	CurrencyId int     `json:"currency_id"` // 币种id
	Symbol     string  `json:"symbol"`      // 币种标识
	Decimals   int     `json:"decimals"`    // 小数点位数
	Name       string  `json:"name"`        // 名称
	Income     float64 `json:"income"`      // 本期收入
	Spend      float64 `json:"spend"`       // 上期支出
	Type       int8    `json:"type"`        // 明细类型（账单业务类型 1充值 2提币 3节点分红 4算力释放 5转入 6升级 7兑换 8转账）
	UpdatedAt  string  `json:"updated_at"`  // 时间
}

// AccountDetailListResp 帐号详情列表
type AccountDetailListResp struct {
	Items []AccountDetailResp `json:"items"` // 帐号详情
	Page  BasePageResp        `json:"page"`  // 分页
}

// AccountUserDetailInfoResp 用户帐号详情信息
type AccountUserDetailInfoResp struct {
	Id           int     `json:"id"`             // 明细id
	Uid          int     `json:"uid"`            // 用户id
	Name         string  `json:"name"`           // 用户帐号
	Email        string  `json:"email"`          // 邮件
	Income       float64 `json:"income"`         // 入账
	Spend        float64 `json:"spend"`          // 支出
	Balance      float64 `json:"balance"`        // 现余额
	LastBalance  float64 `json:"last_balance"`   // 之前余额
	Symbol       string  `json:"symbol"`         // 币种
	UpdatedAt    string  `json:"update_at"`      // 时间
	BlockChainId uint    `json:"block_chain_id"` // 链id
}

// AccountUserDetailListResp 用户帐号详情列表
type AccountUserDetailListResp struct {
	Items []AccountUserDetailInfoResp `json:"items"` // 帐号详情
	Page  BasePageResp                `json:"page"`  // 分页
}

// AccountCurrencyDetailListResp 帐号币种详情列表
type AccountCurrencyDetailListResp struct {
	Info  AccountInfoResp     `json:"info"`  // 币种信息
	Items []AccountDetailResp `json:"items"` // 帐号详情
	Page  BasePageResp        `json:"page"`  // 分页
}

// AccountInfoResp 详情返回结果集
type AccountTFORListInfoResp struct {
	Uid            string `json:"uid"`             // 用户uid
	AccountId      string `json:"account_id"`      // 钱包id
	CurrencyId     int    `json:"currency_id"`     // 币种id
	Balance        string `json:"balance"`         // 余额
	Symbol         string `json:"symbol"`          // 币种标识
	Decimals       int    `json:"decimals"`        // 小数点位数
	Name           string `json:"name"`            // 名称
	UpdatedAt      string `json:"updated_at"`      // 时间
	BlockedBalance string `json:"blocked_balance"` // 冻结金额
	BlockChainId   uint   `json:"block_chain_id"`  // 链id
}

// AccountTFORListResp tfor列表
type AccountTFORListResp struct {
	Items map[string]AccountTFORListInfoResp `json:"items"` // 帐号详情
}
