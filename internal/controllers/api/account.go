package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/controllers/base"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/jwt"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"strings"
)

// AccountInfo 钱包详情
// @Tags Account 钱包帐号
// @Summary 钱包详情接口
// @Description 钱包详情
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} resp.AccountInfoResp
// @Router /account/info [get]
func AccountInfo(c *gin.Context) {
	p := &params.AccountInfoParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	data, err := models.NewAccount().GetUserBalance(core.Orm.DB.New(), p.Base.Uid)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// AccountTFORInfo 钱包TFOR详情
// @Tags Account 钱包帐号
// @Summary 钱包TFOR详情接口
// @Description 钱包TFOR详情
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} resp.AccountInfoResp
// @Router /account/tfor/info [get]
func AccountTFORInfo(c *gin.Context) {
	p := &params.AccountTFORInfoParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	data, err := models.NewAccount().GetUserTFORBalance(core.Orm.DB.New(), p.Base.Uid)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// AccountTFORList 钱包TFOR详情列表
// @Tags Account 钱包帐号
// @Summary 钱包TFOR详情列表接口
// @Description 钱包TFOR详情列表
// @Security ApiKeyAuth
// @Produce json
// @Param uids query string true "用户id多个以,隔开"
// @Success 200 {object} resp.AccountTFORListInfoResp
// @Router /account/tfor/list [get]
func AccountTFORList(c *gin.Context) {
	p := &params.AccountTFORListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	ids := strings.Split(p.Uids, ",")
	if len(ids) == core.DefaultNilNum {
		core.GResp.Failure(c, resp.CodeIllegalParam)
		return
	}
	data, err := models.NewAccount().GetUserTFORBalanceList(core.Orm.DB.New(), ids)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// AccountTransfer 转账动态usdd
// @Tags Account 钱包帐号
// @Summary 转账动态usdd接口
// @Description 转账动态usdd
// @Produce json
// @Security ApiKeyAuth
// @param money formData number true "转账金额"
// @param symbol formData string true "币种标识"
// @Success 200
// @Router /account/transfer [post]
func AccountTransfer(c *gin.Context) {
	p := &params.AccountTransferParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	currency := models.NewCurrency()
	currency.Symbol = p.Symbol
	if err := currency.GetSymbolById(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotAccount)
		return
	}

	account := models.NewAccount()
	account.CurrencyId, account.Uid = currency.ID, p.Base.Uid
	// 判断是否存在账本
	if err := account.IsExistAccount(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotCurrency)
		return
	}

	// 判断是否转账
	if (account.Balance*100 - account.BlockedBalance*100 - p.Money*100) < 0 {
		o.Callback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	}

	balance := account.Balance
	if err := account.UpdateBalance(o, core.OperateToOut, p.Money); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	detail := models.NewAccountDetail()
	detail.Uid, detail.Type = p.Base.Uid, resp.AccountDetailInto
	detail.AccountId, detail.LastBalance = account.ID, balance
	detail.Spend, detail.Balance = p.Money, account.Balance-p.Money
	if err := detail.CreateAccountDetail(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AccountChange 兑换代币
// @Tags Account 钱包帐号
// @Summary 兑换代币接口
// @Description 兑换代币
// @Produce json
// @Security ApiKeyAuth
// @Param currency_id formData int true "代币id"
// @Param change_id formData int true "兑换id"
// @param money formData number true "转账金额"
// @param ratio formData number true "转账比例"
// @Success 200
// @Router /account/change [post]
func AccountChange(c *gin.Context) {
	p := &params.AccountChangeParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	if p.ChangeId == p.CurrencyId {
		core.GResp.Failure(c, errors.New("change is not equal to currency"))
		return
	}

	account := models.NewAccount()
	account.Uid = p.Base.Uid
	o := core.Orm.New().Begin()
	ids := []uint{p.CurrencyId, p.ChangeId}
	data, err := account.GetOrderIdsByInfo(o, ids)
	// 判断是否存在账本
	if err != nil || len(data) < len(ids) {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	//  校验金额
	if (data[p.CurrencyId].Balance*100 - data[p.ChangeId].BlockedBalance*100 - p.Money*100) < 0 {
		o.Callback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	}

	ratio, _ := strconv.ParseFloat(p.Ratio, 64)
	jsonStr, _ := json.Marshal(c.Request.PostForm)
	order := &models.Order{
		Uid:         p.Base.Uid,
		Context:     string(jsonStr),
		CurrencyId:  p.CurrencyId,
		ExchangeUid: p.Base.Uid,
		ExchangeId:  p.ChangeId,
		Balance:     p.Money,
		Ratio:       ratio,
		Status:      models.OrderStatusOk,
		Type:        models.OrderTypeChange,
	}
	if err := order.CreateOrder(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, errors.Errorf("join order fail:%s", err))
		return
	}

	details := make([]models.AccountDetail, 0)
	temp := models.AccountDetail{
		Uid:       p.Base.Uid,
		AccountId: account.ID,
		Type:      resp.AccountDetailConvert,
		OrderId:   order.ID,
	}
	for _, val := range data {
		temp.AccountId, temp.LastBalance = val.ID, val.Balance
		if val.CurrencyId == p.CurrencyId {
			temp.Income, temp.Spend, temp.Balance = float64(core.DefaultNilNum), float64(p.Money), float64(val.Balance-p.Money)
		} else {
			money := p.Money * ratio
			temp.Spend, temp.Income, temp.Balance = float64(core.DefaultNilNum), float64(money), float64(val.Balance+money)
		}
		details = append(details, temp)
	}

	if err := models.NewAccountDetail().CreateAccountDetailAll(o, details); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	account.CurrencyId = p.CurrencyId
	if err := account.UpdateBalance(o, core.OperateToOut, p.Money); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	} else {
		account.CurrencyId = p.ChangeId
		if err = account.UpdateBalance(o, core.OperateToUp, p.Money*ratio); err != nil {
			o.Callback()
			core.GResp.Failure(c, err)
			return
		}
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AccountDetail 钱包明细
// @Tags Account 钱包帐号
// @Summary 钱包明接口
// @Description 钱包明
// @Produce json
// @Security ApiKeyAuth
// @Param pageSize query int true "页面条数"
// @Param page query int true "页数"
// @Success 200 {object} resp.AccountDetailListResp
// @Router /account/detail [get]
func AccountDetail(c *gin.Context) {
	p := &params.AccountDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	detail := models.AccountDetail{
		Uid: p.Base.Uid,
	}
	data, err := detail.GetPageList(core.Orm.New(), p.Page, p.PageSize)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// AccountShareBonus 节点分红入账
// @Tags Account 钱包帐号
// @Summary 节点分红入账接口
// @Description 节点分红入账
// @Produce json
// @Security ApiKeyAuth
// @Param money formData int true "金额"
// @Param symbol formData string true "币种标示"
// @Success 200
// @Router /account/share/bonus [post]
func AccountShareBonus(c *gin.Context) {
	p := &params.AccountShareBonusParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	currency := models.NewCurrency()
	currency.Symbol = p.Symbol
	if err := currency.GetSymbolById(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	account := models.NewAccount()
	account.CurrencyId, account.Uid = currency.ID, p.Base.Uid
	// 判断是否存在账本
	if err := account.IsExistAccount(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotAccount)
		return
	}

	jsonStr, _ := json.Marshal(c.Request.PostForm)
	order := models.Order{
		Uid:        p.Base.Uid,
		Context:    string(jsonStr),
		CurrencyId: currency.ID,
		Balance:    p.Money,
		Form:       models.OrderFormUsdd,
		Status:     models.OrderStatusOk,
		Type:       models.OrderTypeShare,
	}
	if err := order.CreateOrder(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, errors.Errorf("join order fail:%s", err))
		return
	}

	details := &models.AccountDetail{
		Uid:         p.Base.Uid,
		AccountId:   account.ID,
		Type:        resp.AccountDetailShare,
		OrderId:     order.ID,
		Income:      p.Money,
		Balance:     account.Balance + p.Money,
		LastBalance: account.Balance,
	}
	if err := details.CreateAccountDetail(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	account.CurrencyId = currency.ID
	if err := account.UpdateBalance(o, core.OperateToUp, p.Money); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AccountDeposit 提现货币
// @Tags Account 钱包帐号
// @Summary 提现货币接口
// @Description 提现货币
// @Produce json
// @Security ApiKeyAuth
// @Param money formData number true "金额"
// @Param coin_id formData int true "代币id"
// @Param currency_id formData int true "币种id"
// @Param withdrawal_addr_id formData int true "提现地址id"
// @Param password formData string true "支付密码"
// @Success 200
// @Router /account/withdrawal [post]
func AccountWithdrawal(c *gin.Context) {
	p := &params.AccountWithdrawalParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()

	user := models.NewUser()
	user.ID = p.Base.Uid
	if err := user.GetInfo(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	} else if user.PayPassword == core.DefaultNilString {
		o.Callback()
		core.GResp.Failure(c, resp.CodeEmptyPayPassword)
		return
	} else if user.PayPassword != tools.Hash256(p.Password, tools.NewPwdSalt(p.Base.Claims.UserID, 1)) {
		o.Callback()
		core.GResp.Failure(c, resp.CodeErrorPayPassword)
		return
	}

	// 判断充值金额
	account := models.NewAccount()
	account.CurrencyId = p.CurrencyId
	account.Uid = p.Base.Uid
	if err := account.GetOrderUidCurrencyIdByInfo(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotAccount)
		return
	}
	// 获取充值信息
	withdrawal_addr := models.NewWithdrawalAddr()
	withdrawal_addr.ID = p.WithdrawalAddrId
	if err := withdrawal_addr.GetInfo(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotAddr)
		return
	}

	coin := models.NewCoin()
	coin.ID = p.CoinId
	coin_info, err := coin.GetDepositInfo(o)
	if err != nil || coin_info.Status != core.DefaultNilNum {
		o.Callback()
		core.GResp.Failure(c, resp.CodeWithdrawalNotCurrency)
		return
	}

	if coin_info.MinWithdrawal*100 > p.Money*100 {
		o.Callback()
		core.GResp.Failure(c, resp.CodeMinWithdrawal)
		return
	}

	var poundage, money float64
	// 手续费
	if coin_info.WithdrawalFeeType != "fixed" {
		poundage = p.Money * coin_info.WithdrawalFee * 0.01
	} else {
		poundage = coin_info.WithdrawalFee
	}

	// 冻结金额
	if (account.Balance*100 - account.BlockedBalance*100 - poundage*100) < p.Money*100 {
		o.Callback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	}

	money = p.Money + poundage
	// 冻结
	block_detail := models.BlockDetail{
		Uid:         p.Base.Uid,
		AccountId:   account.ID,
		Balance:     account.BlockedBalance + money,
		LastBalance: account.BlockedBalance,
		Income:      money,
	}

	if err := block_detail.CreateBlockDetail(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	// 入账
	if err := account.UpdateBlockBalance(o, core.OperateToUp, money); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	// TODO: 对接提现
	withdrawal_detail := &models.WithdrawalDetail{
		Uid:             p.Base.Uid,
		Address:         withdrawal_addr.Address,
		CoinId:          p.CoinId,
		CurrencyId:      p.CurrencyId,
		AccountId:       account.ID,
		Value:           p.Money,
		Symbol:          coin_info.Symbol,
		Type:            coin_info.Type,
		OrderId:         fmt.Sprintf("%s", uuid.NewV4()),
		Status:          models.WithdrawalStatusToAudit,
		Poundage:        poundage,
		FinancialStatus: models.WithdrawalAudioStatusAwait,
		CustomerStatus:  models.WithdrawalAudioStatusAwait,
	}

	// 不需要审核直接提交
	if coin_info.FinancialStatus > int8(core.DefaultNilNum) {
		withdrawal_detail.FinancialStatus = models.WithdrawalAudioStatusOk
	}
	if coin_info.CustomerStatus > int8(core.DefaultNilNum) {
		withdrawal_detail.CustomerStatus = models.WithdrawalAudioStatusOk
	}

	if withdrawal_detail.FinancialStatus == models.WithdrawalAudioStatusOk && coin.CustomerStatus == models.WithdrawalAudioStatusOk {
		withdrawal_detail.CustomerStatus, withdrawal_detail.FinancialStatus = models.WithdrawalAudioStatusOk, models.WithdrawalAudioStatusOk
		withdrawal_detail.Status = models.WithdrawalStatusThrough
		if msg, err := base.WithdrawalAudioOK(o, withdrawal_detail); err != nil {
			withdrawal_detail.Remark = msg
		}
	}

	if err := withdrawal_detail.CreateWithdrawalDetail(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AccountCurrencyDetail 钱包币种详情
// @Tags Account 钱包帐号
// @Summary 钱包币种详情接口
// @Description 钱包币种详情
// @Produce json
// @Security ApiKeyAuth
// @Param account_id query int true "钱包id"
// @Param pageSize query int true "页面条数"
// @Param page query int true "页数"
// @Param type query string true "类型（all全部，income收入，expend支出）"
// @Success 200 {object} resp.AccountCurrencyDetailListResp
// @Router /account/currency/detail [get]
func AccountCurrencyDetail(c *gin.Context) {
	p := &params.AccountCurrencyDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	switch p.Type {
	case models.AccountCurrentClassAll:
		fallthrough
	case models.AccountCurrentClassIn:
		fallthrough
	case models.AccountCurrentClassUp:
		o := core.Orm.New()
		detail := models.AccountDetail{
			Uid:       p.Base.Uid,
			AccountId: p.AccountId,
		}
		data, err := detail.GetCurrencyPageList(o, p.Page, p.PageSize, p.Type)
		if err != nil {
			core.GResp.Failure(c, err)
			return
		}
		account := models.NewAccount()
		account.ID, account.Uid = p.AccountId, p.Base.Uid
		if data.Info, err = account.GetUserAccountBalance(o); err != nil {
			core.GResp.Failure(c, err)
			return
		}
		coin := models.NewCoin()
		coin.Symbol = data.Info.Symbol
		if err := coin.GetOrderSymbolByInfo(o); err != nil {
			core.GResp.Failure(c, err)
			return
		}
		data.Info.BlockChainId = coin.BlockChainId
		core.GResp.Success(c, data)
		return
	default:
		core.GResp.Failure(c, resp.CodeIllegalParam)
		return
	}
}

// AccountPersonTransfer 个人转账
// @Tags Account 钱包帐号
// @Summary 个人转账接口
// @Description 个人转账
// @Produce json
// @Security ApiKeyAuth
// @param money formData number true "转账金额"
// @param symbol formData string true "币种标识"
// @param email formData string true "邮件"
// @param pay_password formData string true "支付密码"
// @Success 200
// @Router /account/person/transfer [post]
func AccountPersonTransfer(c *gin.Context) {
	p := &params.AccountPersonTransferParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()

	// 查询用户创建用户
	result, err := consul.GetOrderEmailByInfo(p.Email, c.Request.Header.Get(core.HttpHeadToken))
	var data resp.GetOrderEmailUserInfoResp

	if err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	} else if err = mapstructure.Decode(result, &data); err != nil || len(data.UserId) == core.DefaultNilNum {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	}

	// 获取用户id
	param := &params.BaseParam{
		Uid: 0,
		Claims: jwt.CustomClaims{
			UserID: data.UserId,
			Name:   data.UserName,
			Email:  p.Email,
		},
	}

	if err = base.CreateUser(param); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotUser)
		return
	}

	user := models.NewUser()
	user.ID = p.Base.Uid
	_ = user.GetInfo(o)

	if user.PayPassword != tools.Hash256(p.PayPassword, tools.NewPwdSalt(p.Base.Claims.UserID, 1)) {
		fmt.Println(user.PayPassword, tools.Hash256(p.PayPassword, tools.NewPwdSalt(p.Base.Claims.UserID, 1)), p.Base.Claims.UserID)
		o.Callback()
		core.GResp.Failure(c, resp.CodeErrorPayPassword)
		return
	}

	// 获取代笔信息
	currency := models.NewCurrency()
	currency.Symbol = p.Symbol
	if err := currency.GetSymbolById(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotCurrency, err)
		return
	}

	// 获取金额
	account := models.NewAccount()
	account.CurrencyId = currency.ID
	list := make(map[uint]models.Account, 0)
	if list, err = account.GetOrderUidSByCurrencyInfo(o, []uint{p.Base.Uid, param.Uid}); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotAccount, err)
		return
	}

	// 余额不足
	if (list[p.Base.Uid].Balance*100 - list[p.Base.Uid].BlockedBalance*100) < p.Money*100 {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	}

	// 创建订单
	jsonStr, _ := json.Marshal(c.Request.PostForm)
	order := &models.Order{
		Uid:         p.Base.Uid,
		Context:     string(jsonStr),
		CurrencyId:  list[p.Base.Uid].CurrencyId,
		ExchangeUid: param.Uid,
		ExchangeId:  list[param.Uid].CurrencyId,
		Balance:     p.Money,
		Ratio:       float64(core.DefaultNilNum),
		Status:      models.OrderStatusOk,
		Type:        models.OrderTypeTransfer,
		Form:        models.OrderFormTransfer,
	}
	if err := order.CreateOrder(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, errors.Errorf("join order fail:%s", err))
		return
	}

	// 扣费
	if err := AccountOperate(o, list[p.Base.Uid], p.Money, core.OperateToOut, resp.AccountDetailTransfer, order.ID); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeLessMoney)
		return
	} else if err = AccountOperate(o, list[param.Uid], p.Money, core.OperateToUp, resp.AccountDetailTransfer, order.ID); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AccountOperate 账本操作
func AccountOperate(o *gorm.DB, account models.Account, money float64, operate string, types int8, order_id uint) error {
	balance := account.Balance
	if err := account.UpdateBalance(o, operate, money); err != nil {
		return err
	}

	detail := models.NewAccountDetail()
	detail.Uid, detail.Type, detail.OrderId = account.Uid, types, order_id
	detail.AccountId, detail.LastBalance = account.ID, balance
	if operate == core.OperateToOut {
		detail.Spend, detail.Balance = money, account.Balance-money
	} else {
		detail.Income, detail.Balance = money, account.Balance+money
	}
	if err := detail.CreateAccountDetail(o); err != nil {
		return err
	}
	return nil
}

// AccountChargeDetail 个人收款明细
// @Tags Account 钱包帐号
// @Summary 个人收款明细接口
// @Description 个人收款明细
// @Produce json
// @Security ApiKeyAuth
// @Param pageSize query int true "页面条数"
// @Param page query int true "页数"
// @Success 200 {object} resp.AccountDetailListResp
// @Router /account/person/charge/detail [get]
func AccountChargeDetail(c *gin.Context) {
	p := &params.AccountChargeDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	detail := models.AccountDetail{
		Uid: p.Base.Uid,
	}
	data, err := detail.GetGatherPageList(core.Orm.New(), p.Page, p.PageSize)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}
