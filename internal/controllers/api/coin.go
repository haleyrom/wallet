package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/pkg/errors"
)

// ReadCoinList 读取代币列表
// @Tags  Coin 代币功能
// @Summary 读取代币列表接口
// @Description 读取代币列表
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} resp.ReadCoinListResp
// @Router /coin/list [get]
func ReadCoinList(c *gin.Context) {
	data, err := models.NewCoin().GetAll(core.Orm.New())
	if err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	core.GResp.Success(c, map[string]interface{}{
		"items":       data,
		"page":        1,
		"totalPage":   1,
		"currentPage": 1,
		"count":       len(data),
		"pageSize":    1,
	})
	return
}

// ReadCoinInfo 读取代币
// @Tags  Coin 代币功能
// @Summary 读取代币接口
// @Description 读取代币
// @Produce json
// @Security ApiKeyAuth
// @Param coin_id query int true "代币id"
// @Success 200 {object} resp.ReadCoinInfoResp
// @Router /coin/info [get]
func ReadCoinInfo(c *gin.Context) {
	p := &params.ReadCoinInfoParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	coin := models.NewCoin()
	coin.ID = p.CoinId
	data, err := coin.GetInfo(core.Orm.New())
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// UpdateCoin 更新代币
// @Tags  Coin 代币功能
// @Summary 更新代币接口
// @Description 更新代币
// @Produce json
// @Security ApiKeyAuth
// @param currency_id formData int true "币种id"
// @Param coin_id formData int true "代币id"
// @Param symbol formData string true "代币代号"
// @Param name formData string true "代币名称"
// @Param block_chain_id formData int true "链id"
// @Param type formData  string true "币类型（coin,token）"
// @Param confirm_count formData int true "确认数"
// @Param min_deposit formData number true "最小充值金额"
// @Param min_withdrawal formData number true "最小提现金额"
// @Param withdrawal_fee formData number true "提现手续费"
// @Param withdrawal_fee_type formData int true "手续费类型 1：百分比;2：固定类型"
// @Param contract_address formData string true "合约地址 如该是type=token，这里必须输入"
// @Param withdrawal_status formData number true "充值状态0：开启；1:关闭"
// @Param deposit_status formData number  true "提笔状态0：开启；1:关闭"
// @Param customer_status formData number true "客服状态:0 必须1：不必须"
// @Param financial_status formData number true "财务状态:0 必须1：不必须"
// @Success 200
// @Router /coin/update [POST]
func UpdateCoin(c *gin.Context) {
	p := &params.UpdateCoinParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.DB.New()
	chain := models.NewBlockChain()
	chain.ID = p.BlockChainId
	if err := chain.IsExistBlockChain(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotChain)
		return
	}

	coin := models.NewCoin()
	coin.ID = p.Id
	if err := coin.IsExistCoin(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotCoin)
		return
	}

	coin.Symbol, coin.Name = p.Symbol, p.Name
	coin.CurrencyId, coin.BlockChainId = p.CurrencyId, p.BlockChainId
	coin.Type, coin.ConfirmCount = p.Type, p.ConfirmCount
	coin.MinDeposit, coin.MinWithdrawal = p.MinDeposit, p.MinWithdrawal
	coin.WithdrawalFee, coin.WithdrawalFeeType = p.WithdrawalFee, p.WithdrawalFeeType
	coin.ContractAddress, coin.Abi = p.ContractAddress, p.Abi
	coin.WithdrawalStatus, coin.DepositStatus = p.WithdrawalStatus, p.DepositStatus
	coin.CustomerStatus, coin.FinancialStatus = p.CustomerStatus, p.FinancialStatus
	if err := coin.UpdateCoin(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AddCoin 增加代币
// @Tags  Coin 代币功能
// @Summary 增加代币接口
// @Description 增加代币
// @Produce json
// @Security ApiKeyAuth
// @param currency_id formData int true "币种id"
// @Param symbol formData string true "代币代号"
// @Param name formData string true "代币名称"
// @Param block_chain_id query int true "链id"
// @Param type formData string true "币类型（coin,token）"
// @Param confirm_count formData int true "确认数"
// @Param min_deposit formData number true "最小充值金额"
// @Param min_withdrawal formData number true "最小提现金额"
// @Param withdrawal_fee formData number true "提现手续费"
// @Param withdrawal_fee_type formData int true "手续费类型 1：百分比;2：固定类型"
// @Param contract_address formData string true "合约地址 如该是type=token，这里必须输入"
// @Success 200
// @Router /coin/add [post]
func AddCoin(c *gin.Context) {
	p := &params.AddCoinParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.DB.New()
	chain := models.NewBlockChain()
	chain.ID = p.BlockChainId
	if err := chain.IsExistBlockChain(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	coin := models.Coin{
		CurrencyId:        p.CurrencyId,
		Symbol:            p.Symbol,
		Name:              p.Name,
		BlockChainId:      p.BlockChainId,
		Type:              p.Type,
		ConfirmCount:      p.ConfirmCount,
		MinDeposit:        p.MinDeposit,
		MinWithdrawal:     p.MinWithdrawal,
		WithdrawalFee:     p.WithdrawalFee,
		WithdrawalFeeType: p.WithdrawalFeeType,
		ContractAddress:   p.ContractAddress,
		Abi:               p.Abi,
		WithdrawalStatus:  p.WithdrawalStatus,
		DepositStatus:     p.DepositStatus,
		CustomerStatus:    p.CustomerStatus,
		FinancialStatus:   p.FinancialStatus,
	}
	if err := coin.CreateCoin(core.Orm.DB.New()); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	core.GResp.Success(c, resp.EmptyData())
	return
}

// RemoveCoin 删除代币
// @Tags  Coin 代币功能
// @Summary 删除代币接口
// @Description 删除代币
// @Produce json
// @Security ApiKeyAuth
// @Param coin_id formData int true "代币id"
// @Success 200
// @Router /coin/remove [POST]
func RemoveCoin(c *gin.Context) {
	p := &params.RemoveCoinParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	coin := models.NewCoin()
	coin.ID = p.CoinId
	o := core.Orm.DB.New()
	if err := coin.IsExistCoin(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	if err := coin.RmChain(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UpdateCoinStatus 更新代币状态
// @Tags  Coin 代币功能
// @Summary 更新代币状态接口
// @Description 更新代币状态
// @Produce json
// @Security ApiKeyAuth
// @Param coin_id formData int true "代币id"
// @Param status formData int true "状态:0开启;1关闭"
// @Success 200
// @Router /coin/status [POST]
func UpdateCoinStatus(c *gin.Context) {
	p := &params.UpdateCoinStatusParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	coin := models.NewCoin()
	coin.ID = p.CoinId
	if err := coin.IsExistCoin(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotCoin)
		return
	}

	coin.Status = p.Status
	if err := coin.UpdateCoinStatus(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	core.GResp.Success(c, resp.EmptyData())
	return
}

// ReadCoinDepositInfo 读取代币提现信息
// @Tags  Coin 代币功能
// @Summary 读取代币提现信息接口
// @Description 读取代币提现信息
// @Produce json
// @Security ApiKeyAuth
// @Param coin_id query int true "代币id"
// @Success 200 {object} resp.ReadCoinDepositInfoResp
// @Router /coin/deposit [get]
func ReadCoinDepositInfo(c *gin.Context) {
	p := &params.ReadCoinDepositInfoParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, fmt.Errorf("%d", resp.CodeIllegalParam))
		return
	}

	o := core.Orm.New()
	coin := models.NewCoin()
	coin.ID = p.CoinId
	data, err := coin.GetDepositInfo(o)
	if err != nil {
		core.GResp.Failure(c, errors.Errorf("%d", resp.CodeExtractCurrency))
		return
	}

	account := models.NewAccount()
	account.Uid, account.CurrencyId = p.Base.Uid, data.CurrencyId
	data.Money, _ = account.GetUserAvailableBalance(o)
	core.GResp.Success(c, data)
	return

}
