package models

import (
	"fmt"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// Account 帐号
type Account struct {
	gorm.Model
	Uid            uint    `gorm:"index;column:uid;comment:'用户id';"`                 // 用户id
	CurrencyId     uint    `gorm:"column:currency_id;index;comment:'币种id';"`         // 币种id
	Balance        float64 `gorm:"column:balance;default:0;comment:'余额';"`           // 余额
	BlockedBalance float64 `gorm:"column:blocked_balance;default:0;comment:'冻结余额';"` // 冻结余额
}

// Table 表
func GetAccountTable() string {
	return viper.GetString("mysql.prefix") + "account"
}

// NewAccount 初始化
func NewAccount() *Account {
	return &Account{}
}

// GetUserAvailableBalance 获取用户可用余额
func (a *Account) GetUserAvailableBalance(o *gorm.DB) (float64, error) {
	var money float64
	err := o.Raw(fmt.Sprintf("SELECT TRUNCATE(cur.balance-cur.blocked_balance,6) as balance from %s cur where cur.uid = ? and  currency_id = ?", GetAccountTable()), a.Uid, a.CurrencyId).Row().Scan(&money)
	return money, err
}

// GetUserBalance 获取用户余额
func (a *Account) GetUserBalance(o *gorm.DB, uid uint) ([]resp.AccountInfoResp, error) {
	data := make([]resp.AccountInfoResp, 0)
	rows, err := o.Raw(fmt.Sprintf("SELECT acc.uid,acc.id as account_id,acc.currency_id,TRUNCATE(acc.balance-acc.blocked_balance,6) as balance,TRUNCATE(acc.blocked_balance,6) as blocked_balance,cur.symbol,cur.decimals,cur.name,acc.updated_at from %s cur LEFT JOIN %s acc on acc.currency_id = cur.id where acc.uid = ?", GetCurrencyTable(), GetAccountTable()), uid).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.AccountInfoResp
			timer time.Time
		)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data = append(data, item)
			}
		}
	}
	return data, err
}

// GetUserTFORBalance 获取用户TFOR余额
func (a *Account) GetUserTFORBalance(o *gorm.DB, uid uint) (resp.AccountInfoResp, error) {
	var data resp.AccountInfoResp
	rows := o.Raw(fmt.Sprintf("SELECT acc.uid,acc.id as account_id,acc.currency_id,TRUNCATE(acc.balance-acc.blocked_balance,6) as balance,TRUNCATE(acc.blocked_balance,6) as blocked_balance,cur.symbol,cur.decimals,cur.name,acc.updated_at from %s cur LEFT JOIN %s acc on acc.currency_id = cur.id where acc.uid = ? AND cur.symbol = ?", GetCurrencyTable(), GetAccountTable()), uid, "TFOR").Row()
	err := rows.Scan(&data.Uid, &data.AccountId, &data.CurrencyId, &data.Balance, &data.BlockedBalance, &data.Symbol, &data.Decimals, &data.Name, &data.UpdatedAt)
	timer, _ := time.Parse("2006-01-02T15:04:05:08Z", data.UpdatedAt)
	data.UpdatedAt = timer.Format("2006-01-02 15:04:05")
	return data, err
}

// GetUserTFORBalanceList 获取用户TFOR余额列表
func (a *Account) GetUserTFORBalanceList(o *gorm.DB, uids []string) (resp.AccountTFORListResp, error) {
	rows, err := o.Raw(fmt.Sprintf("SELECT u.uid ,acc.id as account_id,acc.currency_id,TRUNCATE(acc.balance-acc.blocked_balance,6) as balance,TRUNCATE(acc.blocked_balance,6) as blocked_balance,cur.symbol,cur.decimals,cur.name,acc.updated_at from %s cur LEFT JOIN %s acc on acc.currency_id = cur.id LEFT JOIN %s u on acc.uid = u.id where u.uid in(?) AND cur.symbol = ?", GetCurrencyTable(), GetAccountTable(), GetUserTable()), uids, "TFOR").Rows()
	defer rows.Close()

	var (
		data  resp.AccountTFORListResp
		item  resp.AccountTFORListInfoResp
		timer time.Time
	)
	data.Items = make(map[string]resp.AccountTFORListInfoResp, 0)
	for rows.Next() {
		if err = o.ScanRows(rows, &item); err == nil {
			item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
			data.Items[item.Uid] = item
		}
	}

	return data, err
}

// GetUserAccountBalance 获取用户钱包余额
func (a *Account) GetUserAccountBalance(o *gorm.DB) (resp.AccountInfoResp, error) {
	var data resp.AccountInfoResp
	rows := o.Raw(fmt.Sprintf("SELECT acc.uid,acc.id as account_id,acc.currency_id,TRUNCATE(acc.balance-acc.blocked_balance,6) as balance,TRUNCATE(acc.blocked_balance,6) as blocked_balance,cur.symbol,cur.decimals,cur.name,acc.updated_at from %s cur LEFT JOIN %s acc on acc.currency_id = cur.id where acc.uid = ? AND acc.id = ?", GetCurrencyTable(), GetAccountTable()), a.Uid, a.ID).Row()
	err := rows.Scan(&data.Uid, &data.AccountId, &data.CurrencyId, &data.Balance, &data.BlockedBalance, &data.Symbol, &data.Decimals, &data.Name, &data.UpdatedAt)
	timer, _ := time.Parse("2006-01-02T15:04:05:08Z", data.UpdatedAt)
	data.UpdatedAt = timer.Format("2006-01-02 15:04:05")
	return data, err
}

// CreateAccount 创建账本
func (a *Account) CreateAccount(o *gorm.DB, currency_id []uint) error {
	var err error
	for _, val := range currency_id {
		a.CurrencyId = val
		a.ID = uint(core.DefaultNilNum)
		if err = o.Create(a).Error; err != nil {
			return err
		}
	}
	return nil
}

// CreateAccountOrderUid 创建账本
func (a *Account) CreateAccountOrderUid(o *gorm.DB, uid []uint) error {
	var err error
	for _, val := range uid {
		a.Uid = val
		a.ID = uint(core.DefaultNilNum)
		if err = o.Create(a).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetInfo 获取用户账本信息
func (a *Account) GetInfo(o *gorm.DB) error {
	err := o.Model(a).Find(a).Error
	return err
}

// IsExistAccount 判断是否存在
func (a *Account) IsExistAccount(o *gorm.DB) error {
	if err := o.Table(GetAccountTable()).Where("currency_id = ? and uid = ?", a.CurrencyId, a.Uid).Find(a).Error; err == nil {
		return nil
	}
	return fmt.Errorf("%s", "Account not exist")
}

// GetOrderIdsByInfo 根据ids,获取用户账本信息
func (a *Account) GetOrderIdsByInfo(o *gorm.DB, ids []uint) (map[uint]Account, error) {
	data := make(map[uint]Account, 0)
	item := make([]Account, 0)
	err := o.Table(GetAccountTable()).Where("currency_id in(?)", ids).Where("uid = ?", a.Uid).Find(&item).Error
	if err == nil {
		for _, val := range item {
			data[val.CurrencyId] = val
		}
	}
	return data, err
}

// GetOrderUidSByCurrencyInfo 根据uids,获取用户账本信息
func (a *Account) GetOrderUidSByCurrencyInfo(o *gorm.DB, ids []uint) (map[uint]Account, error) {
	data := make(map[uint]Account, 0)
	item := make([]Account, 0)
	err := o.Table(GetAccountTable()).Where("uid in(?)", ids).Where("currency_id = ?", a.CurrencyId).Find(&item).Error
	if err == nil {
		for _, val := range item {
			data[val.Uid] = val
		}
	}
	return data, err
}

// GetOrderUidCurrencyIdByInfo 根据用户id获取信息,获取用户账本信息
func (a *Account) GetOrderUidCurrencyIdByInfo(o *gorm.DB) error {
	err := o.Table(GetAccountTable()).Where("uid = ? and currency_id = ?", a.Uid, a.CurrencyId).Find(a).Error
	return err
}

// UpdateBalance 更新余额
func (a *Account) UpdateBalance(o *gorm.DB, operate string, money float64) error {
	filed := fmt.Sprintf("balance %s ?", operate)

	return o.Model(a).Where("uid = ? and currency_id = ?", a.Uid, a.CurrencyId).Updates(map[string]interface{}{
		"balance":    gorm.Expr(filed, money),
		"updated_at": time.Now(),
	}).Error

}

// UpdateBlockBalance 更新冻结余额
func (a *Account) UpdateBlockBalance(o *gorm.DB, operate string, money float64) error {
	filed := fmt.Sprintf("blocked_balance %s ?", operate)
	return o.Model(a).Where("uid = ? and currency_id = ?", a.Uid, a.CurrencyId).Updates(map[string]interface{}{
		"blocked_balance": gorm.Expr(filed, money),
		"updated_at":      time.Now(),
	}).Error
}

// UpdateWithdrawalBalance 更新提现金额
func (a *Account) UpdateWithdrawalBalance(o *gorm.DB, balance, block_balance float64, balance_operate, block_operate string) error {
	return o.Model(a).Where("uid = ? and currency_id = ?", a.Uid, a.CurrencyId).Updates(map[string]interface{}{
		"balance":         gorm.Expr(fmt.Sprintf("balance %s ?", balance_operate), balance),
		"blocked_balance": gorm.Expr(fmt.Sprintf("blocked_balance %s ?", block_operate), block_balance),
		"updated_at":      time.Now(),
	}).Error
}

// GetAccountUserList 获取用户钱包列表
func (a *Account) GetAccountUserList(o *gorm.DB, page, pageSize, start_time, end_timer int, keyword string) (resp.AccountUserDetailListResp, error) {
	data := resp.AccountUserDetailListResp{}
	sql := fmt.Sprintf("SELECT detail.id as id,user.id as uid,user.name,user.email,TRUNCATE(detail.income,6) as income,TRUNCATE(detail.spend,6) as spend,TRUNCATE(detail.balance,6) as balance, TRUNCATE(detail.last_balance,6) as last_balance, currency.symbol,detail.updated_at FROM %s detail LEFT JOIN %s user ON detail.uid = user.id LEFT JOIN %s account ON detail.account_id = account.id LEFT JOIN %s currency ON currency.id = account.currency_id where detail.id > 0 ", GetAccountDetailTable(), GetUserTable(), GetAccountTable(), GetCurrencyTable())
	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s detail LEFT JOIN %s user ON detail.uid = user.id where detail.id > 0 ", GetAccountDetailTable(), GetUserTable())

	if start_time > 0 && end_timer > 0 {
		sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(detail.updated_at) >= %d AND UNIX_TIMESTAMP(detail.updated_at) <= %d ", sql, start_time, end_timer)
		count_sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(detail.updated_at) >= %d AND UNIX_TIMESTAMP(detail.updated_at) <= %d ", count_sql, start_time, end_timer)
	}

	if len(keyword) > 0 {
		sql = fmt.Sprintf("%s AND user.name like '%s'", sql, "%"+keyword+"%")
		count_sql = fmt.Sprintf("%s AND user.name like '%s'", count_sql, "%"+keyword+"%")
	}

	sql = fmt.Sprintf("%s ORDER BY detail.id desc LIMIT %d,%d", sql, (page-1)*pageSize, pageSize)

	rows, err := o.Raw(sql).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.AccountUserDetailInfoResp
			timer time.Time
		)

		data.Items = make([]resp.AccountUserDetailInfoResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data.Items = append(data.Items, item)
			}
		}

	}
	_ = o.Raw(count_sql).Row().Scan(&data.Page.Count)
	data.Page.PageSize = len(data.Items)
	data.Page.CurrentPage = page
	data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	return data, err
}

// GetAdminAccountList 获取后台用户钱包列表
func (a *Account) GetAdminAccountList(o *gorm.DB, page, pageSize, start_time, end_timer int, keyword string) (resp.AccountListResp, error) {
	data := resp.AccountListResp{}
	sql := fmt.Sprintf("SELECT account.id AS id,user.id AS uid,user.name,user.email, TRUNCATE(account.balance,6) AS balance, TRUNCATE ( account.blocked_balance, 6 ) AS blocked_balance, currency.symbol, account.updated_at FROM %s account LEFT JOIN %s user ON account.uid = user.id LEFT JOIN %s currency ON currency.id = account.currency_id where account.id > 0 ", GetAccountTable(), GetUserTable(), GetCurrencyTable())
	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s account LEFT JOIN %s user ON account.uid = user.id where account.id > 0 ", GetAccountTable(), GetUserTable())

	if start_time > 0 && end_timer > 0 {
		sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(account.updated_at) >= %d AND UNIX_TIMESTAMP(account.updated_at) <= %d ", sql, start_time, end_timer)
		count_sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(account.updated_at) >= %d AND UNIX_TIMESTAMP(account.updated_at) <= %d ", count_sql, start_time, end_timer)
	}

	if len(keyword) > 0 {
		sql = fmt.Sprintf("%s AND user.name like '%s'", sql, "%"+keyword+"%")
		count_sql = fmt.Sprintf("%s AND user.name like '%s'", count_sql, "%"+keyword+"%")
	}

	sql = fmt.Sprintf("%s ORDER BY account.id desc LIMIT %d,%d", sql, (page-1)*pageSize, pageSize)

	rows, err := o.Raw(sql).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.AccountAdminInfoResp
			timer time.Time
		)

		data.Items = make([]resp.AccountAdminInfoResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data.Items = append(data.Items, item)
			}
		}
	}
	_ = o.Raw(count_sql).Row().Scan(&data.Page.Count)
	data.Page.PageSize = len(data.Items)
	data.Page.CurrentPage = page
	data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	return data, err
}
