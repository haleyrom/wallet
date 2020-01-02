package admin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/controllers/api"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/jwt"
	"github.com/mitchellh/mapstructure"
	"strconv"
	"time"
)

// JoinRecharge 添加充值
// @Tags Account 后台运营-充值
// @Summary 添加充值接口
// @Description 添加充值
// @Security ApiKeyAuth
// @Produce json
// @Param uid formData string true "用户id"
// @Param symbol formData string true "币种标示"
// @Param money formData number false "金额"
// @Success 200
// @Router /admin/operate/recharge/join [post]
func JoinRecharge(c *gin.Context) {
	p := &params.JoinRechargeParam{}

	var (
		o   = core.Orm.New().Begin()
		err error
	)
	// 绑定参数
	if err = c.ShouldBind(p); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	// 获取用户id
	user := &models.User{
		Uid: p.Uid,
	}

	if err = user.IsExistUser(o); err != nil {
		// 查询用户创建用户
		result, err := consul.GetUserInfo(p.Uid, c.Request.Header.Get(core.HttpHeadToken))
		var data resp.UserInfoResp

		if err != nil {
			o.Rollback()
			core.GResp.Failure(c, resp.CodeNotUser)
			return
		} else if err = mapstructure.Decode(result, &data); err != nil {
			core.GResp.Failure(c, resp.CodeNotUser)
			return
		}

		param := &params.BaseParam{
			Uid: 0,
			Claims: jwt.CustomClaims{
				UserID: data.Id,
				Name:   data.Nickname,
				Email:  data.Email,
			},
		}

		if err = api.CreateUser(param); err != nil {
			core.GResp.Failure(c, resp.CodeNotUser)
			return
		}
	}

	// 获取货币规则
	coin := &models.Coin{
		Symbol: p.Symbol,
	}
	if err = coin.GetOrderSymbolByInfo(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotCoin)
		return
	}

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%d%s", p.Uid, time.Now().Unix())))

	// 创建订单
	deposit_detail := &models.DepositDetail{
		Uid:             user.ID,
		CoinId:          coin.ID,
		Value:           p.Money,
		Symbol:          coin.Symbol,
		Type:            coin.Type,
		CurrencyId:      coin.CurrencyId,
		Source:          models.DepositSourceAdmin,
		FinancialStatus: int8(core.DefaultNilNum),
		Deleted:         models.DepositStatusNotDeleted,
		Status:          int8(core.DefaultNilNum),
		Md5Keys:         hex.EncodeToString(h.Sum(nil)),
	}
	if err := deposit_detail.CreateDepositDetail(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// ReadRechargePage 充值列表
// @Tags Account 后台运营-充值
// @Summary 充值列表接口
// @Description 充值列表
// @Security ApiKeyAuth
// @Produce json
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Param keyword query string false "搜索用户"
// @Param start_time query int false "开始时间"
// @Param end_time query int false "结束时间"
// @Success 200 {object}  resp.ReadAllDepositDetailResp
// @Router /admin/operate/recharge/list [get]
func ReadRechargePage(c *gin.Context) {
	p := &params.ReadRechargePageParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	detail := models.NewDepositDetail()
	o := core.Orm.New()
	detailResp, e := detail.GetAllRechargePageList(o, p.Page, p.PageSize, p.EndTime, p.StartTime, p.Keyword)
	if e != nil {
		core.GResp.Failure(c, e)
	} else {
		core.GResp.Success(c, detailResp)
	}
	return
}

// RemoveRecharge 删除充值
// @Tags Account 后台运营-充值
// @Summary 删除充值接口
// @Description 删除充值
// @Security ApiKeyAuth
// @Produce json
// @Param id formData string true "充值id"
// @Success 200
// @Router /admin/operate/recharge/remove [post]
func RemoveRecharge(c *gin.Context) {
	p := &params.RemoveRechargeParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New()
	deposit_detail := models.NewDepositDetail()
	deposit_detail.ID = p.Id
	if err := deposit_detail.GetSourceAdminOrderIdByInfo(o); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	} else if deposit_detail.FinancialStatus == models.WithdrawalAudioStatusOk {
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	if err := deposit_detail.UpdateDeleted(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// AudioRecharge 审核充值
// @Tags Account 后台运营-充值
// @Summary 审核充值接口
// @Description 审核充值
// @Security ApiKeyAuth
// @Produce json
// @Param id formData string true "充值id"
// @Param financial_id formData int true "财务id"
// @Param status formData int true "状态1：通过；2：不通过"
// @Success 200
// @Router /admin/operate/recharge/audio [post]
func AudioRecharge(c *gin.Context) {
	p := &params.AudioRechargeParam{}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	deposit_detail := models.NewDepositDetail()
	deposit_detail.ID = p.Id
	if err := deposit_detail.GetSourceAdminOrderIdByInfo(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeNotDepositDetail, err)
		return
	} else if deposit_detail.FinancialStatus > models.WithdrawalAudioStatusAwait {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	deposit_detail.Status = p.Status
	if err := deposit_detail.UpdateFinancial(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeAlreadyAudio)
		return
	}

	if p.Status == models.DepositStatusBooked {
		account := models.NewAccount()
		account.Uid, account.CurrencyId = deposit_detail.Uid, deposit_detail.CurrencyId
		if err := account.GetOrderUidCurrencyIdByInfo(o); err != nil {
			o.Callback()
			core.GResp.Failure(c, err)
			return
		}

		// 入账金额
		money := deposit_detail.Value
		account_detail := &models.AccountDetail{
			Uid:         deposit_detail.Uid,
			AccountId:   account.ID,
			Balance:     account.Balance + money,
			LastBalance: account.Balance,
			Income:      money,
			Type:        resp.AccountDetailUp,
		}
		if err := account_detail.CreateAccountDetail(o); err != nil {
			o.Callback()
			core.GResp.Failure(c, err)
			return
		}

		// 入账
		if err := account.UpdateBalance(o, core.OperateToUp, money); err != nil {
			o.Callback()
			core.GResp.Failure(c, err)
			return
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
				Address:     core.DefaultNilString,
				OrderId:     strconv.Itoa(int(deposit_detail.ID)),
			}
			_ = company_stream.CreateCompanyStream(core.Orm.New())
		}()
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return

}
