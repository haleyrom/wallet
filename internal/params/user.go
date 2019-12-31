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
