package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

// AccountUserList 用户钱包列表
// @Tags Account 后台钱包-用户钱包
// @Summary 用户钱包列表接口
// @Description 用户钱包列表
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索帐号"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.AccountUserDetailListResp
// @Router /admin/account/user/list [get]
func AccountUserList(c *gin.Context) {
	p := &params.AccountUserListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	account := models.NewAccount()
	data, _ := account.GetAccountUserList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// AccountUserList 用户提币订单列表
// @Tags Account 后台钱包-用户钱包
// @Summary 用户提币订单列表接口
// @Description 用户提币订单列表
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索帐号"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.WithdrawalDetailAllListResp
// @Router /admin/account/withdrawal/list [get]
func AccountWithdrawalList(c *gin.Context) {
	p := &params.AccountWithdrawalListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	detail := models.NewWithdrawalDetail()
	data, _ := detail.GetAllPageList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// AccountUserList 用户充值流水列表
// @Tags Account 后台钱包-用户钱包
// @Summary 获取用户充值流水接口
// @Description 用户充值流水列表
// @Security ApiKeyAuth
// @Produce json
// @Param page_size query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索帐号"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.AccountUserDetailListResp
// @Router /admin/account/user/deposit-list [get]
func DepositDetailList(c *gin.Context) {
	p := &params.AccountDepositDetailListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	detail := models.NewDepositDetail()
	o := core.Orm.New()
	detailResp, e := detail.GetAllPageList(o, p.Page, p.PageSize, p.EndTime, p.StartTime, p.Keyword)
	if e != nil {
		core.GResp.Failure(c, e)
	} else {
		core.GResp.Success(c, detailResp)
	}
	return
}

// AccountUserList 用户代币装换流水
// @Tags Account 后台钱包-用户钱包
// @Summary 获取用户代币装换流水接口
// @Description 用户代币装换
// @Security ApiKeyAuth
// @Produce json
// @Param page_size query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索帐号名"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.AccountUserDetailListResp
// @Router /admin/account/user/order-list [get]
func OrderList(c *gin.Context) {
	p := &params.AccountOrderListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	order := models.NewOrder()
	o := core.Orm.New()
	data, err := order.GetAllTransOrder(o, p.Page, p.PageSize, p.EndTime, p.StartTime, p.Keyword)
	if err != nil {
		core.GResp.Failure(c, err)
	} else {
		core.GResp.Success(c, data)
	}
	return
}

// AccountWithdrawalDetail 提币详情
// @Tags Account 后台钱包-用户钱包
// @Summary 提币详情接口
// @Description 提币详情
// @Security ApiKeyAuth
// @Produce json
// @Param id query int true "id"
// @Success 200 {object} resp.AdminWithdrawalDetailResp
// @Router /admin/account/withdrawal/detail [get]
func AccountWithdrawalDetail(c *gin.Context) {
	p := &params.AccountWithdrawalDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	detail := models.NewWithdrawalDetail()
	detail.ID = p.Id
	if err := detail.ReadInfo(core.Orm.New()); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	timer, _ := time.Parse("2006-01-02 15:04:05 +0800 CST", detail.UpdatedAt.String())

	core.GResp.Success(c, resp.AdminWithdrawalDetailResp{
		Id:              detail.ID,
		OrderId:         detail.OrderId,
		Symbol:          detail.Symbol,
		Status:          detail.Status,
		CustomerStatus:  detail.CustomerStatus,
		FinancialStatus: detail.FinancialStatus,
		Address:         detail.Address,
		Value:           detail.Value,
		UpdatedAt:       timer.Format("2006-01-02 15:04:05"),
	})
	return
}

// WithdrawalDetailCustomer 提现客服审核
// @Tags Account 后台钱包-用户钱包
// @Summary 提现客户审核接口
// @Description 提现客服审核
// @Security ApiKeyAuth
// @Produce json
// @Param id formData int true "明细id"
// @Param customer_id formData int true "客服id"
// @Param status formData int true "状态1：通过；2：不通过"
// @Success 200
// @Router /admin/account/withdrawal/customer [post]
func WithdrawalDetailCustomer(c *gin.Context) {
	p := &params.AccountWithdrawalDetailCustomerParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	detail := models.NewWithdrawalDetail()
	detail.ID = p.Id
	if err := detail.IsAudioCustomer(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotData)
		return
	} else if detail.Status < models.WithdrawalStatusToAudit || detail.Status > models.WithdrawalStatusInAudit {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	// 财务拒绝不处理
	if detail.FinancialStatus == models.WithdrawalAudioStatusFailure {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	detail.CustomerStatus, detail.CustomerId = p.Status, p.CustomerId
	// 审核成功
	if detail.CustomerStatus == models.WithdrawalAudioStatusOk {
		// 同时审核成功进行处理
		if detail.FinancialStatus == models.WithdrawalAudioStatusOk {
			// 调取提现接口
			msg, err := WithdrawalAudioOK(o, detail)
			// 提交成功
			if err != nil {
				detail.Status, detail.Remark = models.WithdrawalStatusCancel, msg
				detail.CustomerStatus = models.WithdrawalAudioStatusAwait
				_ = detail.UpdateRemark(o)
				// 退款
				if err = WithdrawalAudioRefund(o, detail); err != nil {
					o.Rollback()
					core.GResp.CustomFailure(c, err)
					return
				}
				o.Commit()
				core.GResp.CustomFailure(c, err)
				return
			} else {
				detail.Status = models.WithdrawalStatusSubmit
			}
		} else {
			detail.Status = models.WithdrawalStatusInAudit
		}
	} else {
		// 退款
		if err := WithdrawalAudioRefund(o, detail); err != nil {
			o.Rollback()
			core.GResp.CustomFailure(c, err)
			return
		}
		detail.Status, detail.CustomerStatus = models.WithdrawalStatusNoThrough, models.WithdrawalAudioStatusFailure
	}
	_ = detail.UpdateRemark(o)

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// WithdrawalDetailFinancial 提现财务审核
// @Tags Account 后台钱包-用户钱包
// @Summary 提现财务审核接口
// @Description 提现财务审核
// @Security ApiKeyAuth
// @Produce json
// @Param id formData int true "明细id"
// @Param financial_id formData int true "财务id"
// @Param status formData int true "状态1：通过；2：不通过"
// @Success 200
// @Router /admin/account/withdrawal/financial  [post]
func WithdrawalDetailFinancial(c *gin.Context) {
	p := &params.AccountWithdrawalDetailFinancialParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	detail := models.NewWithdrawalDetail()
	detail.ID = p.Id
	if err := detail.IsAudioCustomer(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	} else if detail.Status < models.WithdrawalStatusToAudit || detail.Status > models.WithdrawalStatusInAudit {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	// 客服拒绝不处理
	if detail.CustomerStatus == models.WithdrawalAudioStatusFailure {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	detail.FinancialStatus, detail.FinancialId = p.Status, p.FinancialId
	// 审核成功
	if detail.FinancialStatus == models.WithdrawalAudioStatusOk {
		// 同时审核成功进行处理
		if detail.CustomerStatus == models.WithdrawalAudioStatusOk {
			// 调取提现接口
			msg, err := WithdrawalAudioOK(o, detail)
			// 提交成功
			if err != nil {
				detail.Status, detail.Remark = models.WithdrawalStatusCancel, msg
				detail.FinancialStatus = models.WithdrawalAudioStatusAwait
				_ = detail.UpdateRemark(o)

				if err = WithdrawalAudioRefund(o, detail); err != nil {
					o.Rollback()
					core.GResp.CustomFailure(c, err)
					return
				}

				o.Commit()
				core.GResp.CustomFailure(c, err)
				return
			}
			detail.Status = models.WithdrawalStatusSubmit
		} else {
			detail.Status = models.WithdrawalStatusInAudit
		}
	} else {
		// 退款
		if err := WithdrawalAudioRefund(o, detail); err != nil {
			o.Rollback()
			core.GResp.CustomFailure(c, err)
			return
		}
		detail.Status, detail.FinancialStatus = models.WithdrawalStatusNoThrough, models.WithdrawalAudioStatusFailure
	}
	_ = detail.UpdateRemark(o)

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// WithdrawalAudioRefund 提现退款
func WithdrawalAudioRefund(o *gorm.DB, detail *models.WithdrawalDetail) error {
	money := detail.Value + detail.Poundage

	account := models.NewAccount()
	account.ID, account.Uid, account.CurrencyId = detail.AccountId, detail.Uid, detail.CurrencyId

	_ = account.IsExistAccount(o)

	if account.BlockedBalance*100 < money*100 {
		return resp.CodeLessMoney
	}

	// 冻结
	block_detail := models.BlockDetail{
		Uid:         detail.Uid,
		AccountId:   detail.AccountId,
		Balance:     account.BlockedBalance - money,
		LastBalance: account.BlockedBalance,
		Spend:       money,
	}

	if err := block_detail.CreateBlockDetail(o); err != nil {
		return err
	}

	// 入账
	if err := account.UpdateBlockBalance(o, core.OperateToOut, money); err != nil {
		return err
	}
	return nil
}

// WithdrawalAudioOK 提现处理
func WithdrawalAudioOK(o *gorm.DB, detail *models.WithdrawalDetail) (string, error) {
	// TODO: 等待调试提币接口
	consul_service, err := consul.ConsulGetServer("blockchain-pay.tfor")
	if err != nil {
		return core.DefaultNilString, err
	}

	coin := models.NewCoin()
	coin.Symbol = detail.Symbol
	contract_address, err := coin.GetConTractAddress(o)
	if err != nil {
		return core.DefaultNilString, err
	}

	company_addr := models.NewCompanyAddr()
	company_addr.Symbol, company_addr.Code = detail.Symbol, models.CodeWithdrawal
	address, err := company_addr.GetOrderSymbolByAddress(o)
	if err != nil {
		return fmt.Sprintf("%s", resp.CodeNotCompanyAddress), resp.CodeNotCompanyAddress
	}

	url := fmt.Sprintf("%s%s", consul_service, "/api/v1/blockchain-pay/ethtereum/withdrawal")

	data := map[string]interface{}{
		"app_id":           viper.GetString("appname"),
		"order_id":         detail.OrderId,
		"symbol":           detail.Symbol,
		"contract_address": contract_address,
		"from_address":     address,
		"to_address":       detail.Address,
		"value":            fmt.Sprintf("%.2f", detail.Value),
	}

	result, err := tools.WithdrawalAudio(data, url, viper.GetString("deposit.Srekey"))

	if result == nil || err != nil || result.Code != http.StatusOK {
		return core.DefaultNilString, errors.Errorf("%s", result.Msg)
	}
	fmt.Println(detail, "=========")
	detail.TransactionHash = result.Data.TransactionHash
	_ = AccountInsertDetail(o, detail)
	return result.Msg, nil
}

// AccountInsertDetail 插入钱包明细
func AccountInsertDetail(o *gorm.DB, detail *models.WithdrawalDetail) error {
	detail.Status = models.WithdrawalStatusThrough
	if err := detail.UpdateStatus(o); err != nil {
		return err
	}

	account := models.NewAccount()
	account.ID, account.Uid, account.CurrencyId = detail.AccountId, detail.Uid, detail.CurrencyId
	if err := account.GetOrderUidCurrencyIdByInfo(o); err != nil {
		return err
	}

	// 入账金额
	money := detail.Value + detail.Poundage
	// 冻结支出
	block_detail := &models.BlockDetail{
		Uid:         detail.Uid,
		AccountId:   detail.AccountId,
		Balance:     account.BlockedBalance - money,
		LastBalance: account.BlockedBalance,
		Spend:       money,
	}

	if err := block_detail.CreateBlockDetail(o); err != nil {
		return err
	}

	account_detail := &models.AccountDetail{
		Uid:         detail.Uid,
		AccountId:   detail.AccountId,
		Balance:     account.Balance - money,
		LastBalance: account.Balance,
		Spend:       money,
		Type:        resp.AccountDetailOut,
	}
	if err := account_detail.CreateAccountDetail(o); err != nil {
		return err
	}

	// 入账
	if err := account.UpdateWithdrawalBalance(o, money, money, core.OperateToOut, core.OperateToOut); err != nil {
		return err
	}

	go func() {
		company_stream := &models.CompanyStream{
			Code:        models.CodeWithdrawal,
			Uid:         account_detail.Uid,
			AccountId:   account_detail.ID,
			Balance:     account_detail.Balance,
			LastBalance: account_detail.LastBalance,
			Income:      account_detail.Income,
			Type:        account_detail.Type,
			Address:     detail.Address,
			OrderId:     detail.OrderId,
		}
		_ = company_stream.CreateCompanyStream(core.Orm.New())
	}()
	return nil
}

// CompanyDepositList 公司充值流水
// @Tags Account 后台钱包-用户钱包
// @Summary 公司充值流水接口
// @Description 公司充值流水
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索帐号"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.CompanyStreamListResp
// @Router /admin/account/company/deposit/list [get]
func CompanyDepositList(c *gin.Context) {
	p := &params.CompanyDepositListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	company_stream := models.NewCompanyStream()
	company_stream.Code = models.CodeDeposit
	data, _ := company_stream.GetList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// CompanyWithdrawalList 公司提现流水
// @Tags Account 后台钱包-用户钱包
// @Summary 公司提现流水接口
// @Description 公司提现流水
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索帐号"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.CompanyStreamListResp
// @Router /admin/account/company/withdrawal/list [get]
func CompanyWithdrawalList(c *gin.Context) {
	p := &params.CompanyWithdrawalListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	company_stream := models.NewCompanyStream()
	company_stream.Code = models.CodeWithdrawal
	data, _ := company_stream.GetList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// CompanyDepositAddrList 公司归集地址
// @Tags Account 后台钱包-用户钱包
// @Summary 公司归集地址接口
// @Description 公司归集地址
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索币种"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.CompanyStreamListResp
// @Router /admin/account/company/deposit_addr/list [get]
func CompanyDepositAddrList(c *gin.Context) {
	p := &params.CompanyDepositListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	company_addr := models.NewCompanyAddr()
	company_addr.Code = models.CodeDeposit
	data, _ := company_addr.GetList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// CompanyWithdrawalList 公司出金地址
// @Tags Account 后台钱包-用户钱包
// @Summary 公司出金地址接口
// @Description 公司出金地址
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索币种"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object} resp.CompanyStreamListResp
// @Router /admin/account/company/withdrawal_addr/list [get]
func CompanyWithdrawalAddrList(c *gin.Context) {
	p := &params.CompanyWithdrawalListParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	company_addr := models.NewCompanyAddr()
	company_addr.Code = models.CodeWithdrawal
	data, _ := company_addr.GetList(core.Orm.New(), p.Page, p.PageSize, p.StartTime, p.EndTime, p.Keyword)
	core.GResp.Success(c, data)
	return
}

// JoinCompanyWithdrawalAddr 新增公司充值地址
// @Tags Account 后台钱包-用户钱包
// @Summary 新增公司充值地址接口
// @Description 新增公司提现地址
// @Security ApiKeyAuth
// @Produce json
// @Param currency_id formData int true "币种id"
// @Param address formData string true "地址"
// @Success 200
// @Router /admin/account/company/withdrawal_addr/join [post]
func JoinCompanyWithdrawalAddr(c *gin.Context) {
	p := &params.JoinCompanyWithdrawalAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	coin := models.NewCoin()
	coin.CurrencyId = p.CurrencyId
	if err := coin.GetOrderCurrencyIdByInfo(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	addr := &models.CompanyAddr{
		Code:         models.CodeWithdrawal,
		BlockChainId: coin.BlockChainId,
		Symbol:       coin.Symbol,
		Type:         coin.Type,
		Address:      p.Address,
	}
	if err := addr.CreateCompanyAddr(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// JoinCompanyDepositAddr 新增公司提现地址
// @Tags Account 后台钱包-用户钱包
// @Summary 新增公司提现地址接口
// @Description 新增公司提现地址
// @Security ApiKeyAuth
// @Produce json
// @Param currency_id formData int true "币种id"
// @Param address formData string true "地址"
// @Success 200
// @Router /admin/account/company/deposit_addr/join [post]
func JoinCompanyDepositAddr(c *gin.Context) {
	p := &params.JoinCompanyDepositAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	coin := models.NewCoin()
	coin.CurrencyId = p.CurrencyId
	if err := coin.GetOrderCurrencyIdByInfo(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	// 判断地址是否合法
	if err := consul.IsWalletAddress(p.Address); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeIllegalAddr)
		return
	}

	addr := &models.CompanyAddr{
		Code:         models.CodeDeposit,
		BlockChainId: coin.BlockChainId,
		Symbol:       coin.Symbol,
		Type:         coin.Type,
		Address:      p.Address,
	}
	if err := addr.CreateCompanyAddr(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UpdateCompanyAddr 更新公司出金/归集地址
// @Tags Account 后台钱包-用户钱包
// @Summary 更新公司出金/归集地址接口
// @Description 更新公司出金/归集地址
// @Security ApiKeyAuth
// @Produce json
// @Param id formData int true "id"
// @Param address formData string true "地址"
// @Success 200
// @Router /admin/account/company/addr/update [post]
func UpdateCompanyAddr(c *gin.Context) {
	p := &params.UpdateCompanyAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	// 判断地址是否合法
	if err := consul.IsWalletAddress(p.Address); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalAddr)
		return
	}

	addr := models.NewCompanyAddr()
	addr.ID = p.Id
	addr.Address = p.Address
	if err := addr.UpdateAddr(core.Orm.New()); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// UpdateCompanyAddrStatus 更新公司出金/归集地址状态
// @Tags Account 后台钱包-用户钱包
// @Summary 更新公司出金/归集地址状态接口
// @Description 更新公司出金/归集地址状态
// @Security ApiKeyAuth
// @Produce json
// @Param id formData int true "id"
// @Param status formData string true "状态（0开启1禁止）"
// @Success 200
// @Router /admin/account/company/addr/status [post]
func UpdateCompanyAddrStatus(c *gin.Context) {
	p := &params.UpdateCompanyAddrStatusParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	addr := models.NewCompanyAddr()
	addr.ID = p.Id
	addr.Status = p.Status
	if err := addr.UpdateAddrStatus(core.Orm.New()); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// CreateCompanyAddr 创建公司地址
// @Tags Account 后台钱包-用户钱包
// @Summary 创建公司地址接口
// @Description 创建公司地址
// @Security ApiKeyAuth
// @Produce json
// @Param symbol formData string true "symbol"
// @Success 200
// @Router /admin/account/company/addr/create [post]
func CreateCompanyAddr(c *gin.Context) {
	p := &params.CreateCompanyAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	coin := models.NewCoin()
	coin.Symbol = p.Symbol
	if err := coin.GetOrderSymbolByInfo(o); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	chain := models.NewBlockChain()
	chain.ID = coin.BlockChainId
	if err := chain.IsExistBlockChain(o); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	if data, err := consul.GetWalletAddress(chain.ChainCode); err == nil {
		core.GResp.Success(c, resp.CreateCompanyAddrResp{
			Address: data.Data.Address,
		})
	} else {
		core.GResp.Failure(c, err)
	}
	return
}
