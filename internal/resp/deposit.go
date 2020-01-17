package resp

import "github.com/shopspring/decimal"

// ReadDepositAddrListResp 充值列表
type ReadDepositAddrListResp struct {
	DepositAddrId uint   `json:"deposit_addr_id"` // 充值地址
	BlockChainId  string `json:"block_chain_id"`  // 区块链名称
	ChainCode     string `json:"chain_code"`      // 链代号
	ChainName     string `json:"chain_name"`      // 链名称
	Address       string `json:"address"`         // 地址
	UpdatedAt     string `json:"updated_at"`      // 更新时间
	Type          string `json:"type"`            // 标识 coin,token
	Status        int8   `json:"status"`          // 状态：0开启;1:停用;2:删除
}

// ReadDepositAddrResp 充值地址详情
type ReadDepositAddrResp struct {
	DepositAddrId uint   `json:"deposit_addr_id"` // 充值地址
	BlockChainId  uint   `json:"block_chain_id"`  // 区块链id
	Address       string `json:"address"`         // 地址
}

// ReadDepositDetailInfoResp 读取充值明细
type ReadDepositDetailInfoResp struct {
	Address    string          `json:"address"`     // 地址
	Value      decimal.Decimal `json:"value"`       // 金额
	Symbol     string          `json:"symbol"`      // 代币代号
	Type       string          `json:"type"`        // 标识 coin,token
	Status     int8            `json:"status"`      // 状态 0确认中,1已确定
	UpdatedAt  string          `json:"updated_at"`  // 更新时间
	BlockCount int             `json:"block_count"` // 确认次数
}

// ReadDepositDetailResp 提现详情明细
type ReadDepositDetailResp struct {
	Items []ReadDepositDetailInfoResp `json:"items"` // 帐号详情
	Page  BasePageResp                `json:"page"`  // 分页
}

//后台查看全部用户明细
type ReadAllDepositDetailInfoResp struct {
	OrderId         int             `json:"order_id"`         // 订单id
	Name            string          `json:"name"`             //用户名
	Uid             string          `json:"uid"`              //用户id
	Value           decimal.Decimal `json:"value"`            // 金额
	Status          int8            `json:"status"`           // 状态 0确认中,1已确定
	UpdatedAt       string          `json:"updated_at"`       // 更新时间
	TransactionHash string          `json:"transaction_hash"` // hash
	Symbol          string          `json:"symbol"`           // 代币代号
	Type            string          `json:"type"`             // 标识 coin,token
	Source          int8            `json:"source"`           // 来源 0：充值1：后台充值
}

//后台查看全部用户，分页
type ReadAllDepositDetailResp struct {
	Items []ReadAllDepositDetailInfoResp `json:"items"` // 数据
	Page  BasePageResp                   `json:"page"`  // 分页
}
