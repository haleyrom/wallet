package base

import (
	"fmt"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// AccountInsertDetail 插入钱包明细
func AccountInsertDetail(o *gorm.DB, detail *models.WithdrawalDetail) error {
	detail.Status = models.WithdrawalStatusSubmit
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

	if account.Balance*100 < money*100 || account.BlockedBalance*100 < money*100 || money*100 > (account.Balance-account.BlockedBalance)*100 {
		logrus.Error("money gt account balance or blocked_balance, %f > %f or %f", money, account.Balance, account.BlockedBalance)
		return resp.CodeLessMoney
	}

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

	fmt.Println("++++++++++++++++++++++++")
	// 入账
	if err := account.UpdateWithdrawalBalance(o, money, money, core.OperateToOut, core.OperateToOut); err != nil {
		return err
	}
	fmt.Println("++++++++++++++++++++++++")

	go func() {
		company_stream := &models.CompanyStream{
			Code:           models.CodeWithdrawal,
			Uid:            account_detail.Uid,
			AccountId:      account_detail.AccountId,
			Balance:        account_detail.Balance,
			LastBalance:    account_detail.LastBalance,
			Income:         account_detail.Income,
			Spend:          account_detail.Spend,
			Type:           account_detail.Type,
			Address:        detail.Address,
			OrderId:        detail.OrderId,
			CallbackJson:   detail.CallbackJson,
			CallbackStatus: detail.CallbackStatus,
		}
		_ = company_stream.CreateCompanyStream(core.Orm.New())
	}()
	return nil
}
