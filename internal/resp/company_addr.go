package resp

// CompanyStreamInfoResp 公司地址信息
type CompanyAddrInfoResp struct {
	Id        int    `json:"id"`        // 地址id
	Address   string `json:"address"`   // 地址
	Symbol    string `json:"symbol"`    // 币种
	Status    int8   `json:"status"`    // 状态0开启;1:停用;2:删除
	UpdatedAt string `json:"update_at"` // 时间
}

// CompanyStreamListResp 公司地址列表
type CompanyAddrListResp struct {
	Items []CompanyAddrInfoResp `json:"items"` // 帐号详情
	Page  BasePageResp          `json:"page"`  // 分页
}

// CreateCompanyAddrResp 创建公司地址
type CreateCompanyAddrResp struct {
	Address string `json:"address"` // 地址
}
