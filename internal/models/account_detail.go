package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// AccountDetail 账本明细
type AccountDetail struct {
	gorm.Model
	Uid         uint    `gorm:"column:uid;default:0;comment:'用户id'"`           // 用户id
	AccountId   uint    `gorm:"column:account_id;default:0;comment:'账本id';"`   // 账本id
	Balance     float64 `gorm:"column:balance;default:0;comment:'本期余额';"`      // 本期余额
	LastBalance float64 `gorm:"column:last_balance;default:0;comment:'上期余额';"` // 上期余额
	Income      float64 `gorm:"column:income;default:0;comment:'本期收入';"`       // 本期收入
	Spend       float64 `gorm:"column:spend;default:0;comment:'本期支出';"`        // 本期支出
	Type        int8    `gorm:"size:3;column:type;default:0;comment:'明细类型'"`   // 明细类型
	OrderId     uint    `gorm:"column:order_id;default:0;comment:'订单id'"`      // 订单id
}

const (
	// AccountCurrentClassAll 全部
	AccountCurrentClassAll string = "all"
	// AccountCurrentClassIn 入账
	AccountCurrentClassIn string = "income"
	// AccountCurrentClassUp 支出
	AccountCurrentClassUp string = "expend"
)

// GetAccountDetailTable 获取表
func GetAccountDetailTable() string {
	return viper.GetString("mysql.prefix") + "account_detail"
}

// NewAccountDetail 初始化
func NewAccountDetail() *AccountDetail {
	return &AccountDetail{}
}

// CreateAccountDetail  创建账本明细
func (a *AccountDetail) CreateAccountDetail(o *gorm.DB) error {
	return o.Create(a).Error
}

// LastAccountDetail 最后一条数据
func (a *AccountDetail) LastAccountDetail(o *gorm.DB) error {
	return o.Model(a).Where("id = ?", a.ID).Last(&a).Error
}

func (a *AccountDetail) CreateAccountDetailAll(o *gorm.DB, items []AccountDetail) error {
	timer := time.Now().Format("2006-01-02 15:04:05")
	sql := fmt.Sprintf("INSERT INTO `%s` (`uid`,`account_id`,`balance`,`last_balance`,`income`,`spend`,`type`,`order_id`,`created_at`,`updated_at`) VALUES ", GetAccountDetailTable())
	// 循环data数组,组合sql语句
	for key, value := range items {
		if len(items)-1 == key {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("(%d,'%d','%0.2f','%0.2f','%0.2f','%0.2f','%d','%d','%s','%s');", value.Uid, value.AccountId, value.Balance, value.LastBalance, value.Income, value.Spend, value.Type, value.OrderId, timer, timer)
		} else {
			sql += fmt.Sprintf("(%d,'%d','%0.2f','%0.2f','%0.2f','%0.2f','%d','%d','%s','%s'),", value.Uid, value.AccountId, value.Balance, value.LastBalance, value.Income, value.Spend, value.Type, value.OrderId, timer, timer)
		}
	}
	return o.Exec(sql).Error

}

// GetPageList 获取分页列表
func (a *AccountDetail) GetPageList(o *gorm.DB, page, pageSize int) (resp.AccountDetailListResp, error) {
	data := resp.AccountDetailListResp{}
	rows, err := o.Raw(fmt.Sprintf("SELECT account.currency_id,currency.symbol,currency.name,currency.decimals,detail.income,detail.spend,detail.type,detail.updated_at FROM %s AS detail LEFT JOIN %s account ON account.id = detail.account_id LEFT JOIN %s currency ON currency.id = account.currency_id where detail.uid = ? ORDER BY detail.id desc LIMIT ?,?", GetAccountDetailTable(), GetAccountTable(), GetCurrencyTable()), a.Uid, (page-1)*pageSize, pageSize).Rows()
	defer rows.Close()

	if err == nil {
		var (
			timer time.Time
			item  resp.AccountDetailResp
		)
		data.Items = make([]resp.AccountDetailResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				timer, _ = time.Parse("2006-01-02T15:04:05+08:00", item.UpdatedAt)
				item.UpdatedAt = timer.Format("2006-01-02 15:04:05")
				data.Items = append(data.Items, item)
			}
		}

		o.Table(GetAccountDetailTable()).Where("uid = ?", a.Uid).Count(&data.Page.Count)
		data.Page.PageSize = len(data.Items)
		data.Page.CurrentPage = page
		data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	}
	return data, err
}

// GetCurrencyPageList 获取币种分页列表
func (a *AccountDetail) GetCurrencyPageList(o *gorm.DB, page, pageSize int, types string) (resp.AccountCurrencyDetailListResp, error) {
	data := resp.AccountCurrencyDetailListResp{}
	sql := fmt.Sprintf("SELECT account.currency_id,currency.symbol,currency.name,currency.decimals,detail.income,detail.spend,detail.type,detail.updated_at FROM %s AS detail LEFT JOIN %s account ON account.id = detail.account_id LEFT JOIN %s currency ON currency.id = account.currency_id where detail.uid = ? and detail.account_id = ? ", GetAccountDetailTable(), GetAccountTable(), GetCurrencyTable())

	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s AS detail where detail.uid = ? and detail.account_id = ? ", GetAccountDetailTable())

	switch types {
	case AccountCurrentClassAll:
	case AccountCurrentClassIn:
		sql += " and detail.income > 0 "
		count_sql += " and detail.income > 0 "
	case AccountCurrentClassUp:
		sql += " and detail.spend > 0 "
		count_sql += " and detail.spend > 0 "
	}

	sql += " ORDER BY detail.id desc LIMIT ?,? "

	rows, err := o.Raw(sql, a.Uid, a.AccountId, (page-1)*pageSize, pageSize).Rows()
	defer rows.Close()

	if err == nil {
		var (
			timer time.Time
			item  resp.AccountDetailResp
		)
		data.Items = make([]resp.AccountDetailResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				timer, _ = time.Parse("2006-01-02T15:04:05+08:00", item.UpdatedAt)
				item.UpdatedAt = timer.Format("2006-01-02 15:04:05")
				data.Items = append(data.Items, item)
			}
		}

		_ = o.Raw(count_sql, a.Uid, a.AccountId).Row().Scan(&data.Page.Count)
		data.Page.PageSize = len(data.Items)
		data.Page.CurrentPage = page
		data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	}
	return data, err
}
