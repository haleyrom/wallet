package params

// UserIdParam 用户id参数解析
type UserIdParam BaseBindParam

// UpdatePayPasswordParam 更新支付密码参数
type UpdatePayPasswordParam struct {
	Base     *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Password string     `json:"password" form:"password"  binding:"required"`
}

// SetPayPassWordHandlerParam 设置支付密码
type SetPayPassWordHandlerParam struct {
	Base     *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Password string     `json:"password" form:"password"  binding:"required"`
}

// ReSetPayWordHandlerParam 重置支付密码
type ReSetPayWordHandlerParam struct {
	Base     *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Code     string     `json:"code" form:"code"  binding:"required"`
	Password string     `json:"password" form:"password"  binding:"required"`
	Email    string     `json:"email" form:"email"`
}

// ChargeQrCodeParam 收费二维码生成参数
type ChargeQrCodeParam struct {
	Base   *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Type   int        `json:"type" form:"type" binding:"required"`
	Money  string     `json:"money" form:"money" `
	Symbol string     `json:"symbol" form:"symbol" binding:"required"`
}

// PaymentQrCodeParam 支付二维码生成参数
type PaymentQrCodeParam struct {
	Base   *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Symbol string     `json:"symbol" form:"symbol" binding:"required"`
	Money  string     `json:"money" form:"money" `
	Type   int        `json:"type" form:"type"`
}

// UserPayInfoParam 用户收款信息
type UserPayInfoParam struct {
	Base    *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Code    uint       `json:"code" form:"code" binding:"required"`
	Money   string     `json:"money" form:"money" binding:"gt=0"`
	Symbol  string     `json:"symbol" form:"symbol" binding:"required"`
	From    string     `json:"from" form:"from" binding:"required"`
	OrderId string     `json:"order_id" form:"order_id"`
}

// UserChangeParam 用户收款
type UserChangeParam struct {
	Base        *BaseParam `json:"claims" form:"claims"  binding:"required"`
	Code        uint       `json:"code" form:"code" binding:"required"`
	Money       float64    `json:"money" form:"money" binding:"gt=0"`
	Symbol      string     `json:"symbol" form:"symbol" binding:"required"`
	From        string     `json:"from" form:"from" binding:"required"`
	OrderId     string     `json:"order_id" form:"order_id"`
	PayPassword string     `json:"pay_password" form:"pay_password"`
}
