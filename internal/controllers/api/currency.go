package api

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"strings"
)

// ReadCurrencyList 读取币种列表
// @Tags  currency 币种功能
// @Summary 读取币种列表接口
// @Description 读取币种列表
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} resp.ReadCurrencyListResp
// @Router /currency/list [get]
func ReadCurrencyList(c *gin.Context) {
	p := &params.ReadCurrencyListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	data, err := models.NewCurrency().GetAll(o)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}

	if len(data) > core.DefaultNilNum {
		account := models.NewAccount()
		account.Uid = p.Base.Uid
		for key, val := range data {
			account.CurrencyId = val.CurrencyId
			if money, err := account.GetUserAvailableBalance(o); err == nil {
				data[key].Money = money
			}
		}
	}
	core.GResp.Success(c, data)
	return
}

// ReadCurrencyList 读取兑换币种列表
// @Tags  currency 币种功能
// @Summary 读取兑换币种列表接口
// @Description 读取兑换币种列表
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} resp.ReadCurrencyTransferListResp
// @Router /currency/transfer_list [get]
func ReadCurrencyTransferList(c *gin.Context) {
	p := &params.ReadCurrencyTransferListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	data := resp.ReadCurrencyTransferListResp{
		List:     make([]resp.ReadCurrencyListResp, 0),
		Transfer: make([]resp.ReadCurrencyListResp, 0),
	}
	list, err := models.NewCurrency().GetAll(o)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}

	if len(list) > core.DefaultNilNum {
		account := models.NewAccount()
		account.Uid = p.Base.Uid
		for key, val := range list {
			account.CurrencyId = val.CurrencyId
			if money, err := account.GetUserAvailableBalance(o); err == nil {
				list[key].Money = money
				if val.Name == "USDD" {
					data.Transfer = append(data.Transfer, list[key])
				} else {
					data.List = append(data.List, list[key])
				}
			}
		}
	}
	core.GResp.Success(c, data)
	return
}

// UpdateCurrency 更新币种
// @Tags  currency 币种功能
// @Summary 更新币种接口
// @Description 更新币种
// @Produce json
// @Security ApiKeyAuth
// @Param currency_id query string true "币种id"
// @Param symbol query string true "币种标识"
// @Param name query string true "币种名称"
// @Param decimals query int true "小数点数"
// @Success 200
// @Router /currency/update [post]
func UpdateCurrency(c *gin.Context) {
	p := &params.UpdateCurrencyParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	currency := models.NewCurrency()
	currency.ID = p.CurrencyId
	o := core.Orm.DB.New()
	if err := currency.IsExistCurrency(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotCurrency)
		return
	}

	currency.Symbol = p.Symbol
	currency.Decimals = p.Decimals
	currency.Name = p.Name
	//currency.UpdatedAt = time.Now()
	if err := currency.UpdateCurrency(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AddCurrency 增加币种
// @Tags  currency 币种功能
// @Summary 增加币种接口
// @Description 增加币种
// @Produce json
// @Security ApiKeyAuth
// @Param symbol query string true "币种标识"
// @Param name query string true "币种名称"
// @Param decimals query int true "小数点数"
// @Success 200
// @Router /currency/add [post]
func AddCurrency(c *gin.Context) {
	p := &params.AddCurrencyParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	currency := &models.Currency{
		Symbol:   p.Symbol,
		Name:     p.Name,
		Decimals: p.Decimals,
	}

	o := core.Orm.DB.New()
	if err := currency.CreateCurrency(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	go func() {
		// 创建账本
		user := models.NewUser()
		if ids, err := user.GetAllByUid(o); err == nil {
			account := models.NewAccount()
			account.CurrencyId = currency.ID
			_ = account.CreateAccountOrderUid(o, ids)
		}
	}()

	core.GResp.Success(c, resp.EmptyData())
	return
}

// RemoveCurrency 删除币种
// @Tags  currency 币种功能
// @Summary 删除币种接口
// @Description 删除币种
// @Produce json
// @Security ApiKeyAuth
// @Param currency_id query string true "币种id"
// @Success 200
// @Router /currency/remove [POST]
func RemoveCurrency(c *gin.Context) {
	p := &params.RmCurrencyParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	currency := models.NewCurrency()
	currency.ID = p.CurrencyId
	o := core.Orm.DB.New()
	if err := currency.IsExistCurrency(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotCurrency)
		return
	}

	if err := currency.RmCurrency(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UpdateCurrencyStatus 更新币种状态
// @Tags  currency 币种功能
// @Summary 更新币种状态接口
// @Description 更新币种状态
// @Produce json
// @Security ApiKeyAuth
// @Param currency_id formData string true "币种id"
// @Param status formData int true "状态:0开启;1关闭"
// @Success 200
// @Router /currency/status [POST]
func UpdateCurrencyStatus(c *gin.Context) {
	p := &params.UpdateCurrencyStatusParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	currency := models.NewCurrency()
	currency.ID = p.CurrencyId
	o := core.Orm.DB.New()
	if err := currency.IsExistCurrency(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotCurrency)
		return
	}

	currency.Status = p.Status
	if err := currency.UpdateCurrencyStatus(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// CurrencyQuote 货币报价
// @Tags  currency 币种功能
// @Summary 货币报价接口
// @Description 货币报价
// @Produce json
// @Security ApiKeyAuth
// @Param quote_currency query string true "报价货币"
// @Success 200 {object} resp.CurrencyQuoteInfoResp
// @Router /currency/quote [get]
func CurrencyQuote(c *gin.Context) {
	p := &params.CurrencyQuoteParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	quote := models.NewQuote()
	quote.QuoteCurrency = strings.ToLower(p.QuoteCurrency)
	data, err := quote.GetQuoteCurrencyByList(core.Orm.New())
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}
