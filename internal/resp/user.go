package resp

// UserInfoResp 用户信息解析
type UserInfoResp struct {
	Id        string  `json:"id" mapstructure:"id"`                // id
	Nickname  string  `json:"nickname" mapstructure:"nickname"`    // nickname
	Email     string  `json:"email" mapstructure:"email"`          // email
	Rcode     string  `json:"rcode" mapstructure:"rcode"`          // rcode
	FatherId  string  `json:"father_id" mapstructure:"father_id"`  // father_id
	HeadImage string  `json:"head_image" mapstructure:"headimage"` // head_image
	Fans      float64 `json:"fans" mapstructure:"fans"`            // fans
}

// GetOrderEmailUserInfoResp 获取用户email信息
type GetOrderEmailUserInfoResp struct {
	UserId       string `json:"user_id" mapstructure:"user_id"`             // 用户id
	UserName     string `json:"user_name" mapstructure:"user_name"`         // 用户名
	AccountState int    `json:"account_state" mapstructure:"account_state"` // 账户状态
}

// ChargeQrCodeResp 收费二维码
type ChargeQrCodeResp struct {
	UserName string  `json:"user_name"` // 用户名称
	Email    string  `json:"email"`     // 邮箱
	Qrcode   string  `json:"qrcode"`    // 二维码
	MinMoney float64 `json:"min_money"` // 最小转账金额
}

// UserPayInfoResp 用户支付信息输出
type UserPayInfoResp struct {
	OrderId string `json:"order_id" form:"order_id"` // 订单号
	Symbol  string `json:"symbol" form:"symbol"`     // symbol标示
	Email   string `json:"email" form:"email"`       // 邮件
	Money   string `json:"money" form:"money"`       // money
}
