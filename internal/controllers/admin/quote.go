package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"strings"
)

// CreateCurrencyQuote 创建货币汇率
// @Tags Account 后台钱包-用户钱包
// @Summary 创建货币汇率接口
// @Description 创建货币汇率
// @Security ApiKeyAuth
// @Produce json
// @Param base_currency formData string true "基础货币symbol"
// @Param quote_currency formData string true "报价货币symbol"
// @Param price formData number true "金额"
// @Success 200
// @Router /admin/currency/quote/create [post]
func CreateCurrencyQuote(c *gin.Context) {
	p := &params.CreateQuoteParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	} else if p.QuoteCurrency == p.BaseCurrency {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	p.BaseCurrency = strings.ToUpper(p.BaseCurrency)
	p.QuoteCurrency = strings.ToUpper(p.QuoteCurrency)

	quote := &models.Quote{
		Code: fmt.Sprintf("%s-%s", p.BaseCurrency, p.QuoteCurrency),
	}

	o := core.Orm.New().Begin()
	if err := quote.IsExistQuote(o); err == nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeExistQuote)
		return
	}

	quote.BaseCurrency, quote.QuoteCurrency = p.BaseCurrency, p.QuoteCurrency
	quote.Price = p.Price
	if err := quote.CreateQuote(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	if err := CreateQuoteHistory(o, quote); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UpdateCurrencyQuote 更新货币汇率
// @Tags Account 后台钱包-用户钱包
// @Summary 更新货币汇率接口
// @Description 更新货币汇率
// @Security ApiKeyAuth
// @Produce json
// @Param id formData int true "汇率id"
// @Param price formData number true "金额"
// @Success 200
// @Router /admin/currency/quote/update [post]
func UpdateCurrencyQuote(c *gin.Context) {
	p := &params.UpdateQuoteParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	quote := models.NewQuote()
	quote.ID = p.Id
	if err := quote.GetOrderIdByInfo(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	quote.Price = p.Price
	if err := quote.UpdateQuotePrice(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	if err := CreateQuoteHistory(o, quote); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// CreateQuoteHistory 创建汇率历史记录
func CreateQuoteHistory(o *gorm.DB, quote *models.Quote) error {
	history := &models.QuoteHistory{
		Code:          quote.Code,
		BaseCurrency:  quote.BaseCurrency,
		QuoteCurrency: quote.QuoteCurrency,
		Price:         quote.Price,
	}
	err := history.CreateQuoteHistory(o)
	return err
}

// ReadQuotePage 获取汇率分页
// @Tags Account 后台钱包-用户钱包
// @Summary 获取汇率分页接口
// @Description 获取汇率分页
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索标示"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.CurrencyQuoteListResp
// @Router /admin/currency/quote/list [get]
func ReadQuotePage(c *gin.Context) {
	p := &params.ReadQuotePageParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	quote := models.NewQuote()
	data, _ := quote.GetAllPageList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// ReadQuoteHistoryPage 获取汇率历史记录分页
// @Tags Account 后台钱包-用户钱包
// @Summary 获取汇率历史记录分页接口
// @Description 获取汇率历史记录分页
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索标示"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.CurrencyQuoteListResp
// @Router /admin/currency/quote_history/list [get]
func ReadQuoteHistoryPage(c *gin.Context) {
	p := &params.ReadQuotePageParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	history := models.NewQuoteHistory()
	data, _ := history.GetAllPageList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}
