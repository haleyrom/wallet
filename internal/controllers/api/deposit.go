package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

// CreateDeposit 创建钱包地址
func CreateDeposit(o *gorm.DB, uid uint) error {
	chain := models.NewBlockChain()
	if items, _ := chain.GetAll(core.Orm.New()); len(items) > 0 {
		addr := models.DepositAddr{
			Uid: uid,
		}

		for _, item := range items {
			go func(code string) {
				if data, err := consul.GetWalletAddress(code); err == nil && data != nil && len(data.Data.Address) > 0 {
					addr.BlockChainId = uint(item.Id)
					addr.Address = data.Data.Address
					addr.OrderId = data.Data.OrderId
					_ = addr.CreateDepositAddr(o)
				} else {
					logrus.Errorf("RegisterWalletAddr data :%v,failure :%v", data, err)
				}
			}(item.ChainCode)
		}
	}
	return nil
}

// ReadDepositAddList 读取钱包地址列表
// @Tags  DepositAdd 钱包地址
// @Summary 读取钱包地址列表接口
// @Description 读取钱包地址列表
// @Produce json
// @Success 200 {object} resp.ReadDepositAddrListResp
// @Router /deposit/list [get]
func ReadDepositAddList(c *gin.Context) {
	data, err := models.NewDepositAddr().GetAll(core.Orm.DB.New())
	if err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// JoinDepositDetail 写入充值记录
// @Tags  DepositAdd 钱包地址
// @Summary 写入充值记录接口
// @Description 写入充值记录列表
// @Produce json
// @Security ApiKeyAuth
// @Param address formData string true "地址"
// @Param value formData string true "充值金额"
// @Param block_number formData string true "充值区块高度"
// @Param block_count formData string true "区块确认数"
// @Param transaction_hash formData string true "事务"
// @Success 200
// @Router /deposit/join [post]
func JoinDepositDetail(c *gin.Context) {
	p := &params.JoinDepositDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	o := core.Orm.New().Begin()
	deposit := &models.DepositAddr{
		Address: p.Address,
	}
	if err := deposit.IsExistAddress(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotAddr)
		return
	}

	coin := models.NewCoin()
	coin.ID = p.CoinId
	if err := coin.IsExistCoin(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, resp.CodeNotCoin)
		return
	}

	detail := &models.DepositDetail{
		CoinId:          p.CoinId,
		Address:         p.Address,
		Value:           p.Value,
		BlockNumber:     p.BlockNumber,
		BlockCount:      p.BlockCount,
		TransactionHash: p.TransactionHash,
	}
	if err := detail.CreateDepositDetail(o); err != nil {
		o.Callback()
		core.GResp.Failure(c, err)
		return
	}

	// TODO: 充值计算
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// ReadDepositAddr  读取充值地址
// @Tags  DepositAdd 钱包地址
// @Summary 读取充值地址接口
// @Description 读取提币地址
// @Security ApiKeyAuth
// @Produce json
// @Param block_chain_id formData int true "链id"
// @Success 200 {object} resp.ReadDepositAddrResp
// @Router /deposit/info [get]
func ReadDepositAddr(c *gin.Context) {
	p := &params.ReadDepositAddrParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	o := core.Orm.New()
	deposit := &models.DepositAddr{
		Uid:          p.Base.Uid,
		BlockChainId: p.BlockChainId,
	}
	data, err := deposit.ReadWithdrawalAddr(o)
	if err != nil {
		chain := models.NewBlockChain()
		chain.ID = p.BlockChainId

		if err = chain.IsExistBlockChain(o); err != nil {
			core.GResp.Failure(c, resp.CodeNotChain, err)
			return
		}

		if wallet, err := consul.GetWalletAddress(chain.ChainCode); err == nil {
			deposit.Uid, deposit.BlockChainId = p.Base.Uid, p.BlockChainId
			deposit.Address, deposit.OrderId = wallet.Data.Address, wallet.Data.OrderId
			_ = deposit.CreateDepositAddr(o)

			core.GResp.Success(c, resp.ReadDepositAddrResp{
				DepositAddrId: deposit.ID,
				BlockChainId:  deposit.BlockChainId,
				Address:       wallet.Data.Address,
			})
		} else {
			core.GResp.Failure(c, resp.CodeNotAddr, err)
		}

		return
	}
	core.GResp.Success(c, data)
	return
}

// TopUpDeposit  充值
// @Tags  DepositAdd 钱包地址
// @Summary 充值接口 code=101119 入账成功
// @Description 充值
// @Security ApiKeyAuth
// @Produce json
// @Param address formData string true "地址"
// @Param money formData number true "充值费用"
// @Param block_number formData int true "确认数"
// @Param block_count formData int true "数量"
// @Param transaction_hash formData string true "订单"
// @Param symbol formData string true "币种标识"
// @Param type formData string true "类型"
// @Param hash formData string true "校验"
// @Param contract_address formData string true "合约地址"
// @Success 200
// @Router /deposit/top_up [post]
func TopUpDeposit(c *gin.Context) {
	p := &params.TopUpDepositParam{}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	// 验证签名
	data := make(map[string]interface{}, 0)
	jsonStr, _ := json.Marshal(p)
	_ = json.Unmarshal(jsonStr, &data)
	if hash := tools.GenerateSign(data, viper.GetString("deposit.Srekey")); hash != p.Hash {
		core.GResp.Failure(c, resp.CodeErrSign)
		return
	}

	o := core.Orm.New().Begin()
	coin := &models.Coin{
		Symbol: p.Symbol,
		Type:   p.Type,
	}

	if err := coin.GetOrderSymbolTypeByCoin(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	addr := &models.DepositAddr{
		Address: p.Address,
	}
	if err := addr.GetAddressByInfo(o); err != nil {
		o.Rollback()
		core.GResp.Success(c, resp.CodeNotAddr)
		return
	}

	money, _ := strconv.ParseFloat(p.Money, 64)
	block_number, _ := strconv.Atoi(p.BlockNumber)
	block_count, _ := strconv.Atoi(p.BlockCount)
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%s", p.Address, p.TransactionHash)))

	// 入提现地址
	detail := &models.DepositDetail{
		Uid:             addr.Uid,
		CoinId:          coin.ID,
		CurrencyId:      coin.CurrencyId,
		Address:         p.Address,
		Value:           money,
		BlockNumber:     block_number,
		BlockCount:      block_count,
		TransactionHash: p.TransactionHash,
		Symbol:          p.Symbol,
		Type:            p.Type,
		Status:          models.DepositStatusNotBooked,
		ContractAddress: p.ContractAddress,
		Key:             hex.EncodeToString(h.Sum(nil)),
	}

	var confirm bool
	// 判断是否已经存在订单
	if err := detail.IsKey(o); err != nil {
		//// 是否到确认数量
		//if coin.ConfirmCount <= block_count {
		//	confirm = true
		//	detail.Status = models.DepositStatusBooked
		//}

		// 创建详细 成功获取id，失败说明已存在重新获取
		if err := detail.CreateDepositDetail(core.Orm.New()); err != nil {
			_ = detail.IsKey(o)
		}
	}

	// 确认数量 小于 记录时过滤
	if block_count <= detail.BlockCount {
		o.Commit()
		core.GResp.Success(c, resp.EmptyData())
		return
	} else if coin.ConfirmCount <= block_count && detail.Status == models.DepositStatusNotBooked {
		// 确实数量大于记录且大于代币确实数量且状态未确认
		confirm = true
		detail.Status = models.DepositStatusBooked
	}

	// 更新区块数量
	detail.BlockCount = block_count
	if err := detail.UpdateBlockCount(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}

	if confirm == true {
		// 入账 判断账本
		account := models.NewAccount()
		account.Uid, account.CurrencyId = addr.Uid, coin.CurrencyId
		if err := account.IsExistAccount(o); err != nil {
			o.Rollback()
			core.GResp.Failure(c, resp.CodeNotAccount)
			return
		}

		account_detail := &models.AccountDetail{
			Uid:         addr.Uid,
			AccountId:   account.ID,
			Balance:     account.Balance + money,
			LastBalance: account.Balance,
			Income:      money,
			Type:        resp.AccountDetailUp,
		}
		if err := account_detail.CreateAccountDetail(o); err != nil {
			o.Rollback()
			core.GResp.Failure(c, resp.CodeNotAccount)
			return
		}

		if err := account.UpdateBalance(o, core.OperateToUp, money); err != nil {
			o.Rollback()
			core.GResp.Failure(c, resp.CodeNotAccount)
			return
		}

		go func() {
			company_stream := &models.CompanyStream{
				Code:        models.CodeDeposit,
				Uid:         account_detail.Uid,
				AccountId:   account_detail.ID,
				Balance:     account_detail.Balance,
				LastBalance: account_detail.LastBalance,
				Income:      account_detail.Income,
				Type:        account_detail.Type,
				Address:     detail.Address,
				OrderId:     strconv.Itoa(int(detail.ID)),
			}
			_ = company_stream.CreateCompanyStream(core.Orm.New())
		}()

		o.Commit()
		core.GResp.CustomSuccess(c, resp.CodeDepositOk, resp.EmptyData())
		return
	}

	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}

// ReadDepositDetail 读取充值详情
// @Tags  DepositAdd 钱包地址
// @Summary 读取充值详情接口
// @Description 读取充值详情
// @Produce json
// @Security ApiKeyAuth
// @Param pageSize query int true "长度"
// @Param page query int true "页数"
// @Success 200 {object} resp.ReadDepositDetailResp
// @Router /deposit/detail [get]
func ReadDepositDetail(c *gin.Context) {
	p := &params.ReadDepositDetailParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	deposit_detail := &models.DepositDetail{
		Uid: p.Base.Uid,
	}
	data, err := deposit_detail.GetPageList(core.Orm.New(), p.Page, p.PageSize)
	if err != nil {
		core.GResp.Failure(c, resp.CodeNotData)
		return
	}

	core.GResp.Success(c, data)
	return
}
