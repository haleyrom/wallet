package params

// ReadListBlockChain 读取链列表
type ReadListBlockChainParam BaseBindParam

// ReadListAccountBlockChainParam 读取账本链列表
type ReadListSymbolBlockChainParam struct {
	Base   *BaseParam `json:"claims" form:"claims"`
	Symbol string     `json:"symbol" form:"symbol"`
}

// UpdateBlockChainParam 更新参数
type UpdateBlockChainParam struct {
	Base         *BaseParam `json:"claims" form:"claims"`
	BlockChainId uint       `json:"chain_id" form:"chain_id"  binding:"required"`
	ChainCode    string     `json:"chain_code" form:"chain_code"  binding:"required"`
	Name         string     `json:"name" form:"name"  binding:"required"`
}

// AddBlockChainParam 更新参数
type AddBlockChainParam struct {
	Base      *BaseParam `json:"claims" form:"claims"`
	ChainCode string     `json:"chain_code" form:"chain_code"  binding:"required"`
	Name      string     `json:"name" form:"name"  binding:"required"`
}

// RemoveBlockChainParam 删除参数
type RemoveBlockChainParam struct {
	Base         *BaseParam `json:"claims" form:"claims"`
	BlockChainId uint       `json:"chain_id" form:"chain_id" binding:"required"`
}
