package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// QuoteHistory 兑换历史记录
type QuoteHistory struct {
	gorm.Model
	Code          string  `gorm:"size:250;column:code;unique;comment:'标示';"`           // 标示
	BaseCurrency  string  `gorm:"size:250;column:base_currency;comment:'基础货币/大写字母';"`  // 基础货币
	QuoteCurrency string  `gorm:"size:250;column:quote_currency;comment:'报价货币/大写字母';"` // 基础货币
	Price         float64 `gorm:"column:price;default:0;comment:'金额';"`                // 金额
}

// GetQuoteHistoryTable 表
func GetQuoteHistoryTable() string {
	return viper.GetString("mysql.prefix") + "quote_history"
}

// NewQuoteHistory NewQuoteHistory
func NewQuoteHistory() *QuoteHistory {
	return &QuoteHistory{}
}

// CreateQuoteHistory 创建历史记录
func (q *QuoteHistory) CreateQuoteHistory(o *gorm.DB) error {
	return o.Create(q).Error
}

// GetAllPageList 获取分页列表
func (q *QuoteHistory) GetAllPageList(o *gorm.DB, page, pageSize int, start_time, end_timer int, keyword string) (resp.CurrencyQuoteListResp, error) {
	data := resp.CurrencyQuoteListResp{}
	sql := fmt.Sprintf("SELECT id,code,base_currency,quote_currency,price,updated_at FROM %s  where id > 0 ", GetQuoteTable())
	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s  where id > 0 ", GetQuoteTable())

	if start_time > 0 && end_timer > 0 {
		sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(updated_at) >= %d AND UNIX_TIMESTAMP(updated_at) <= %d ", sql, start_time, end_timer)
		count_sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(updated_at) >= %d AND UNIX_TIMESTAMP(updated_at) <= %d ", count_sql, start_time, end_timer)
	}

	if len(keyword) > 0 {
		sql = fmt.Sprintf("%s AND code like '%s'", sql, "%"+keyword+"%")
		count_sql = fmt.Sprintf("%s AND code like '%s'", count_sql, "%"+keyword+"%")
	}

	sql = fmt.Sprintf("%s ORDER BY id LIMIT %d,%d", sql, (page-1)*pageSize, pageSize)

	rows, err := o.Raw(sql).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.CurrencyQuoteInfoResp
			timer time.Time
		)

		data.Items = make([]resp.CurrencyQuoteInfoResp, 0)
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
