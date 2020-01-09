package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// CompanyStream 公司流水
type CompanyStream struct {
	gorm.Model
	Code        string  `gorm:"column:code;comment:'code'"`                    // 账户名称
	Uid         uint    `gorm:"column:uid;default:0;comment:'用户id'"`           // 用户id
	AccountId   uint    `gorm:"column:account_id;default:0;comment:'账本id';"`   // 账本id
	Balance     float64 `gorm:"column:balance;default:0;comment:'本期余额';"`      // 本期余额
	LastBalance float64 `gorm:"column:last_balance;default:0;comment:'上期余额';"` // 上期余额
	Income      float64 `gorm:"column:income;default:0;comment:'本期收入';"`       // 本期收入
	Spend       float64 `gorm:"column:spend;default:0;comment:'本期支出';"`        // 本期支出
	Type        int8    `gorm:"size:3;column:type;default:0;comment:'明细类型'"`   // 明细类型
	OrderId     string  `gorm:"column:order_id;comment:'订单id'"`                // 订单id
	Address     string  `gorm:"column:address;comment:'充值地址'"`                 // 充值地址
}

const (
	// CodeDeposit
	CodeDeposit string = "account_deposit"
	// CodeWithdrawal
	CodeWithdrawal string = "account_withdrawal"
)

// GetCompanyStreamTable 获取表
func GetCompanyStreamTable() string {
	return viper.GetString("mysql.prefix") + "company_stream"
}

// NewCompanyStream 初始化
func NewCompanyStream() *CompanyStream {
	return &CompanyStream{}
}

// CreateCompanyDeposit  创建公司充值流水
func (c *CompanyStream) CreateCompanyStream(o *gorm.DB) error {
	return o.Create(c).Error
}

// GetList 获取企业充值列表
func (c *CompanyStream) GetList(o *gorm.DB, page, pageSize, start_time, end_timer int, keyword string) (resp.CompanyStreamListResp, error) {
	data := resp.CompanyStreamListResp{}
	sql := fmt.Sprintf("SELECT detail.address,detail.id as id,user.id as uid,user.name,user.email,TRUNCATE(detail.income,6) as income,TRUNCATE(detail.spend,6) as spend,TRUNCATE(detail.balance,6) as balance, TRUNCATE(detail.last_balance,6) as last_balance, currency.symbol,detail.updated_at,detail.order_id FROM %s detail LEFT JOIN %s user ON detail.uid = user.id LEFT JOIN %s account ON detail.account_id = account.id LEFT JOIN %s currency ON currency.id = account.currency_id where detail.order_id > 0 AND detail.code = '%s' ", GetCompanyStreamTable(), GetUserTable(), GetAccountTable(), GetCurrencyTable(), c.Code)

	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s detail LEFT JOIN %s user ON detail.uid = user.id where detail.order_id > 0 AND detail.code = '%s' ", GetCompanyStreamTable(), GetUserTable(), c.Code)

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
			item  resp.CompanyStreamInfoResp
			timer time.Time
		)

		data.Items = make([]resp.CompanyStreamInfoResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				timer, _ = time.Parse("2006-01-02T15:04:05+08:00", item.UpdatedAt)
				item.UpdatedAt = timer.Format("2006-01-02 15:04:05")
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
