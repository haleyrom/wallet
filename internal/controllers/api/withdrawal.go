package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/controllers/base"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// ReadWithdrawalAddrList  读取提币地址列表
// @Tags Withdrawal 提现功能
// @Summary 读取提币地址列表接口
// @Description 读取提币地址列表
// @Produce json
// @Security ApiKeyAuth
// @Param block_chain_id query int true "链id"
// @Param currency_id query int true "币种id"
// @Success 200 {object} resp.WithdrawalAddrResp
// @Router /withdrawal/addr/list [get]
func ReadWithdrawalAddrList(c *gin.Context) {
	p := &params.ReadWithdrawalAddrListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	withdrawalAddr := models.NewWithdrawalAddr()
	withdrawalAddr.Uid = p.Base.Uid
	withdrawalAddr.CurrencyId, withdrawalAddr.BlockChainId = p.CurrencyId, p.BlockChainId
	data, err := withdrawalAddr.GetAll(core.Orm.New())
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// CreateWithdrawalAddr  创建提现地址
// @Tags Withdrawal 提现功能
// @Summary 创建提现地址接口
// @Description 创建提现地址
// @Produce json
// @Security ApiKeyAuth
// @Param block_chain_id formData int true "链id"
// @Param currency_id formData int true "币种id"
// @Param address formData string true "地址"
// @Param name formData string true "名称"
// @Success 200
// @Router /withdrawal/addr/add [post]
func CreateWithdrawalAddr(c *gin.Context) {
	p := &params.CreateWithdrawalAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	coin := models.NewCoin()
	coin.BlockChainId, coin.CurrencyId = p.BlockChainId, p.CurrencyId
	if err := coin.GetOrderChainIdByInfo(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotChain)
		return
	}

	currency := models.NewCurrency()
	currency.ID = p.CurrencyId
	if err := currency.IsExistCurrency(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotCurrency)
		return
	}

	// 判断地址是否合法
	if err := consul.IsWalletAddress(p.Address); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalAddr, err)
		return
	}

	withdrawal := &models.WithdrawalAddr{
		Uid:           p.Base.Uid,
		BlockChainId:  p.BlockChainId,
		Address:       strings.ToLower(p.Address),
		CurrencyId:    p.CurrencyId,
		Name:          p.Name,
		AddressSource: models.WithdrawalAddrBack,
		Type:          coin.Type,
	}

	deposit := models.NewDepositAddr()
	deposit.Address = p.Address
	if err := deposit.IsAddress(o); err == nil {
		withdrawal.AddressSource = models.WithdrawalAddrLocal
	}

	if err := withdrawal.CreateWithdrawalAddr(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UpdateWithdrawalAddr 更新提现地址
// @Tags Withdrawal 提现功能
// @Summary 更新提现地址接口
// @Description 更新提现地址
// @Produce json
// @Security ApiKeyAuth
// @Param withdrawal_addr_id formData int true "提现地址id"
// @Param address formData string true "地址"
// @Param name formData string true "名称"
// @Success 200
// @Router /withdrawal/addr/update [post]
func UpdateWithdrawalAddr(c *gin.Context) {
	p := &params.UpdateWithdrawalAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	//chain := models.NewBlockChain()
	//chain.ID = p.BlockChainId
	//if err := chain.IsExistBlockChain(o); err != nil {
	//	core.GResp.Failure(c, resp.CodeNotChain)
	//	return
	//}

	// 判断地址是否合法
	if err := consul.IsWalletAddress(p.Address); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalAddr, err)
		return
	}

	withdrawal := models.NewWithdrawalAddr()
	withdrawal.Name, withdrawal.Address = p.Name, strings.ToLower(p.Address)
	withdrawal.CurrencyId = withdrawal.CurrencyId
	withdrawal.ID, withdrawal.Uid = p.WithdrawalAddrId, p.Base.Uid
	if err := withdrawal.UpdateWithdrawalAddr(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// RemoveWithdrawalAddr 删除提现地址
// @Tags Withdrawal 提现功能
// @Summary 删除提现地址接口
// @Description 删除提现地址
// @Produce json
// @Security ApiKeyAuth
// @Param withdrawal_addr_id query int true "提现地址id"
// @Success 200
// @Router /withdrawal/addr/remove [post]
func RemoveWithdrawalAddr(c *gin.Context) {
	p := &params.RemoveWithdrawalAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	withdrawal := models.NewWithdrawalAddr()
	withdrawal.ID, withdrawal.Uid = p.WithdrawalAddrId, p.Base.Uid
	if err := withdrawal.RmWithdrawalAddr(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// ReadWithdrawalDetail 获取提现明细
// @Tags Withdrawal 提现功能
// @Summary 获取提现明细接口
// @Description 获取提现明细
// @Produce json
// @Security ApiKeyAuth
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Success 200 {object} resp.WithdrawalDetailListResp
// @Router /withdrawal/detail [get]
func ReadWithdrawalDetail(c *gin.Context) {
	p := &params.ReadWithdrawalDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	withdrawal_detail := models.NewWithdrawalDetail()
	withdrawal_detail.Uid = p.Base.Uid
	data, err := withdrawal_detail.GetPageList(core.Orm.New(), p.Page, p.PageSize)
	if err != nil {
		core.GResp.Failure(c, resp.CodeNotData)
		return
	}

	core.GResp.Success(c, data)
	return
}

// WithdrawalCallback 提现回调
// @Tags Withdrawal 提现功能
// @Summary 提现回调接口
// @Description 提现回调
// @Produce json
// @Security ApiKeyAuth
// @Param app_id formData string true "app_id"
// @Param order_id formData string true "order_id"
// @Param transaction_hash formData string true "transaction_hash"
// @Param block_number formData string true "block_number"
// @Param from_address formData string true "from_address"
// @Param to_address formData string true "to_address"
// @Param symbol formData string true "symbol"
// @Param contract_address formData string true "contract_address"
// @Param value formData string true "value"
// @Param code formData string true "code"
// @Param message formData string true "message"
// @Param hash formData string true "hash"
// @Success 200
// @Router /withdrawal/callback [post]
func WithdrawalCallback(c *gin.Context) {
	p := &params.WithdrawalCallbackParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	// 验签
	data := make(map[string]interface{}, 0)
	jsonStr, _ := json.Marshal(p)
	_ = json.Unmarshal(jsonStr, &data)
	if hash := tools.GenerateSign(data, viper.GetString("deposit.Srekey")); hash != p.Hash {
		core.GResp.CustomFailure(c, resp.CodeErrSign)
		return
	}

	o := core.Orm.New().Begin()
	detail := models.NewWithdrawalDetail()
	detail.OrderId = p.OrderId
	// 根据订单ID获取信息
	if err := detail.GetOrderIdByInfo(o); err != nil {
		o.Rollback()
		core.GResp.CustomFailure(c, err)
		return
	}

	postStr, _ := json.Marshal(c.Request.PostForm)
	detail.CallbackStatus, detail.CallbackJson = p.Code, string(postStr)
	logrus.Info(detail.CallbackJson)
	switch p.Code {
	case "105004":
		// 已提交
		detail.TransactionHash = p.TransactionHash
		if detail.Status == models.WithdrawalStatusThrough {
			// 入账
			if err := base.AccountInsertDetail(o, detail); err != nil {
				o.Rollback()
				core.GResp.CustomFailure(c, err)
				return
			}
		}
		o.Commit()
		core.GResp.Success(c, resp.EmptyData())
		return
	case "105005":
		// 已汇出
	case "105006":
		// 异常
		fallthrough
	default:
		// 汇款失败
		detail.Remark = p.Message
		if err := detail.UpdateOrderIdRemark(o); err != nil {
			o.Commit()
			core.GResp.CustomFailure(c, err)
			return
		}
		o.Commit()
		core.GResp.Success(c, resp.EmptyData())
		return
	}

	if detail.FinancialStatus != models.WithdrawalAudioStatusOk || detail.CustomerStatus != models.WithdrawalAudioStatusOk || (detail.Status != models.WithdrawalStatusSubmit && detail.Status != models.WithdrawalStatusThrough) {
		o.Rollback()
		core.GResp.CustomFailure(c, errors.New("order status not submit"))
		return
	} else if detail.Status == models.WithdrawalStatusThrough {
		if err := base.AccountInsertDetail(o, detail); err != nil {
			o.Rollback()
			core.GResp.CustomFailure(c, err)
			return
		}
		o.Commit()
		core.GResp.CustomFailure(c, resp.CodeNotData)
		return
	}

	detail.BlockCount, _ = strconv.Atoi(p.BlockCount)
	detail.Status, detail.TransactionHash = models.WithdrawalStatusOk, p.TransactionHash
	if err := detail.UpdateStatus(o); err != nil {
		o.Rollback()
		core.GResp.CustomFailure(c, err)
		return
	}
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// WithdrawalOrderTypeByAddr 根据type获取充值地址
// @Tags Withdrawal 提现功能
// @Summary 根据type获取充值地址接口
// @Description 根据type获取充值地址
// @Produce json
// @Security ApiKeyAuth
// @Param type query string true "类型"
// @Success 200
// @Router /withdrawal/type [get]
func WithdrawalOrderTypeByAddr(c *gin.Context) {
	p := &params.WithdrawalOrderTypeByAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	withdrawal_addr := models.NewWithdrawalAddr()
	withdrawal_addr.Uid, withdrawal_addr.Type = p.Base.Uid, p.Types
	if err := withdrawal_addr.GetTypeByInfo(core.Orm.New()); err != nil {
		core.GResp.Failure(c, resp.CodeNotData, err)
		return
	}
	core.GResp.Success(c, resp.WithdrawalOrderTypeByAddrResp{
		Address: withdrawal_addr.Address,
	})
	return
}
