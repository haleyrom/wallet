package core

import (
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/config"
	"github.com/haleyrom/wallet/pkg/storage"
	"sync"
)

// PaymentChannel 支付channel
type PaymentChannel struct {
	MapChan map[string]chan resp.ChangeInfoResp
	MapTime map[int][]string
}

var (
	// Conf 配置
	Conf config.Configure
	// Orm 数据
	Orm storage.MysqlClient
	// GResp 返回
	GResp *resp.Resp
	// UserInfo 用户信息
	UserInfoPool *sync.Pool

	// DefaultNilString DefString
	DefaultNilString string = ""

	// DefaultNilNum DefaultNilNum
	DefaultNilNum int = 0

	// HttpHeadToken http请求包头的token数据
	HttpHeadToken string = "Authorization"

	// OperateToUp 进账
	OperateToUp string = "+"

	// OperateToOut 入账
	OperateToOut string = "-"

	// PayChan 支付通道
	PayChan PaymentChannel

	// EmptyStruct
	EmptyStruct = struct{}{}
)

// 初始化
func init() {
	GResp = new(resp.Resp)
	// 用户信息磁化
	UserInfoPool = &sync.Pool{
		New: func() interface{} {
			return &params.BaseParam{}
		},
	}

	PayChan = PaymentChannel{
		MapChan: make(map[string]chan resp.ChangeInfoResp, 0),
		MapTime: make(map[int][]string, 0),
	}

}
