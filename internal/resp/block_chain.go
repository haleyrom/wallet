package resp

// ReadBlockChainListResp 读取链列表结果
type ReadBlockChainListResp struct {
	Id        int    `json:"id"`         // 链id
	ChainCode string `json:"chain_code"` // 链标识
	Name      string `json:"name"`       // 链名称
	Type      string `json:"type"`       // 链类型
}

// ReadOrderSymbolByChainResp 读取链列表结果
type ReadOrderSymbolByChainResp struct {
	CoinId    int    `json:"coin_id"`    // 代币id
	Id        int    `json:"id"`         // 链id
	ChainCode string `json:"chain_code"` // 链标识
	Name      string `json:"name"`       // 链名称
	Type      string `json:"type"`       // 链类型
}
