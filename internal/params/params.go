package params

import "github.com/haleyrom/wallet/pkg/jwt"

// BaseParam 通用参数
type BaseParam struct {
	Uid    uint             `json:"uid"`
	Claims jwt.CustomClaims `json:"claims"`
}

// BaseBindParam 通用绑定参数
type BaseBindParam struct {
	Base *BaseParam `json:"claims" form:"claims"`
}

// BaseListParam 用户钱包列表
type BaseListParam struct {
	Base      *BaseParam `json:"claims" form:"claims"`
	StartTime int        `json:"start_time" form:"start_time"`
	EndTime   int        `json:"end_time" form:"end_time"`
	Keyword   string     `json:"keyword" form:"keyword"`
	Page      int        `json:"page" form:"page"  binding:"required"`
	PageSize  int        `json:"pageSize" form:"pageSize" binding:"required"`
}
