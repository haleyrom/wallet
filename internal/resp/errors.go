package resp

import (
	"net/http"
	"sync"
)

var (
	// okDataPool http响应成功数据池化
	okDataPool *sync.Pool

	// statusCodeMsgs 错误代码消息
	statusCodeMsgs map[StatusCode]string

	// emptyData 空数据
	emptyData = &struct{}{}
)

const (
	// CodeUnknow 未知
	CodeUnknow StatusCode = -1

	// CodeOk 请求响应
	CodeOk StatusCode = http.StatusOK

	// CodeAuth 暂无权限
	CodeAuth StatusCode = http.StatusUnauthorized

	// CodeInternalServerError 内部服务出错
	CodeInternalServerError StatusCode = http.StatusInternalServerError

	// CodeIllegalParam 参数不合法
	CodeIllegalParam StatusCode = 101100

	// CodeNoToken 请求参数必需要有token
	CodeNoToken StatusCode = 101101

	// CodeIllegalToken token不合法
	CodeIllegalToken StatusCode = 101102

	// CodeNotTeam 团队不存在
	CodeNotTeam StatusCode = 101103

	// CodeExistTeam 团队存在
	CodeExistTeam StatusCode = 101104

	// CodeNotUser 用户不存在
	CodeNotUser StatusCode = 101105

	// CodeExistUser 用户存在
	CodeExistUser StatusCode = 101106

	// CodeNotProject 项目不存在
	CodeNotProject StatusCode = 101107

	// CodeExistProject 项目存在
	CodeExistProject StatusCode = 101108

	// CodeExtractCurrency 暂不支持提取
	CodeExtractCurrency StatusCode = 101109

	// CodeLessMoney 余额不足
	CodeLessMoney StatusCode = 101110

	// CodeEmptyPayPassword 支付密码为空
	CodeEmptyPayPassword StatusCode = 101111

	// CodeErrorPayPassword 支付密码错误
	CodeErrorPayPassword StatusCode = 101112

	// CodeNotData 暂无数据
	CodeNotData StatusCode = 101113

	// CodeNotAccount 账本不存在
	CodeNotAccount StatusCode = 101114

	// CodeNotAddr 地址不存在
	CodeNotAddr StatusCode = 101115

	// CodeNotCoin 代币不可用
	CodeNotCoin StatusCode = 101116

	// CodeNotChain 链不存在
	CodeNotChain StatusCode = 101117

	// CodeNotCurrency  货币不存在
	CodeNotCurrency StatusCode = 101118

	// CodeDepositErr 充值不成功
	CodeDepositErr StatusCode = 101119

	//CodeCodeError  邮箱验证码错误
	CodeCodeError StatusCode = 101120

	// CodeAlreadyAudio 已审核
	CodeAlreadyAudio StatusCode = 101121

	// CodeCustomerNotAudio  客服未审核
	CodeCustomerNotAudio StatusCode = 101122

	// CodeSetPayPassword 设置支付密码失败，已被初始化
	CodeSetPayPassword StatusCode = 101123

	// CodeErrSign 签名cutout
	CodeErrSign StatusCode = 101124

	// CodeExistQuote 货币汇率存在
	CodeExistQuote StatusCode = 101125

	// CodeMinWithdrawal 最小提笔数量
	CodeMinWithdrawal StatusCode = 101126

	// CodeNotCompanyAddress 公司地址不存在
	CodeNotCompanyAddress StatusCode = 101127

	// CodeIllegalAddr 地址不合法
	CodeIllegalAddr StatusCode = 101128

	// CodeFinancialNotAudio  客服未审核
	CodeFinancialNotAudio StatusCode = 101129

	// CodeNotDepositDetail 充值记录不存在
	CodeNotDepositDetail StatusCode = 101130

	// CodeIllegalPassword 密码不合法
	CodeIllegalPassword StatusCode = 101131

	// CodeWithdrawalNotCurrency 不支持该币提现
	CodeWithdrawalNotCurrency StatusCode = 101132

	// CodeNotOrderId 订单号不存在
	CodeNotOrderId StatusCode = 101133

	// CodeOrderStatusOK 订单号已支付
	CodeOrderStatusOK StatusCode = 101135
)

// StatusCode 状态码
type StatusCode int

// Error 实现error接口
func (c StatusCode) Error() string {
	if msg, ok := statusCodeMsgs[c]; ok {
		return msg
	}
	return statusCodeMsgs[CodeUnknow]
}

func init() {
	okDataPool = &sync.Pool{
		New: func() interface{} {
			return &ResData{
				Code: int(CodeOk),
				Msg:  "ok",
			}
		},
	}

	statusCodeMsgs = map[StatusCode]string{
		CodeUnknow:                "unknow status",
		CodeOk:                    "请求成功",
		CodeInternalServerError:   "服务繁忙,请稍后！",
		CodeIllegalParam:          "参数错误",
		CodeNoToken:               "请求参数必需要有token",
		CodeIllegalToken:          "token不合法",
		CodeAuth:                  "暂无权限",
		CodeNotTeam:               "该团队不存在/已解散",
		CodeExistTeam:             "已加入该团队,您可以直接进入",
		CodeNotUser:               "用户不存在",
		CodeExistUser:             "用户存在",
		CodeNotProject:            "该项目不存在/已解散",
		CodeExistProject:          "已加入该项目,您可以直接进入",
		CodeExtractCurrency:       "暂不支持提取该币种",
		CodeLessMoney:             "余额不足",
		CodeEmptyPayPassword:      "支付密码为空",
		CodeErrorPayPassword:      "支付密码错误",
		CodeNotData:               "暂无数据",
		CodeNotAccount:            "该钱包不存在/已冻结",
		CodeNotAddr:               "该地址不存在/不可用",
		CodeNotCoin:               "该代币不存在/不可用",
		CodeNotChain:              "该链不存在/不可用",
		CodeNotCurrency:           "货币不存在/不可用",
		CodeDepositErr:            "充值不成功",
		CodeCodeError:             "邮箱验证码错误",
		CodeAlreadyAudio:          "已审核，请勿重复操作",
		CodeCustomerNotAudio:      "客服未审核，请耐心等待",
		CodeSetPayPassword:        "支付密码已设置，请重置",
		CodeErrSign:               "签名错误",
		CodeExistQuote:            "存在兑换汇率",
		CodeMinWithdrawal:         "提取数量不小于最低提币数量",
		CodeNotCompanyAddress:     "公司地址不存在",
		CodeIllegalAddr:           "地址不合法/不可用",
		CodeFinancialNotAudio:     "财务未审核",
		CodeNotDepositDetail:      "充值记录不存在",
		CodeIllegalPassword:       "密码不合法",
		CodeWithdrawalNotCurrency: "暂不支持该币提现",
		CodeNotOrderId:            "订单号不存在",
		CodeOrderStatusOK:         "订单号已支付",
	}
}

// ResData http响应数据封包
type ResData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// EmptyData 空数据
func EmptyData() *struct{} {
	return emptyData
}

// Ok 返回成功响应数据封包
func Ok(data interface{}) *ResData {
	resData := okDataPool.Get().(*ResData)
	resData.Data = data
	return resData
}

// RecycleOk 回收Ok响应数据封包
func RecycleOk(data *ResData) {
	data.Data = nil // Notice:一定要赋nil避免泄内存
	okDataPool.Put(data)
}

// ErrCodeMsg 返回错误消息响应数据
func ErrCodeMsg(code StatusCode, msg ...string) *ResData {
	var errMsg string
	var exist bool
	if len(msg) > 0 {
		errMsg = msg[0]
	} else if errMsg, exist = statusCodeMsgs[code]; !exist {
		errMsg = statusCodeMsgs[CodeInternalServerError]
	}

	return &ResData{
		Code: int(code),
		Msg:  errMsg,
		Data: emptyData,
	}
}
