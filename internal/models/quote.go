package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// Quote 兑换
type Quote struct {
	gorm.Model
	Code          string  `gorm:"index;size:250;column:code;unique;comment:'标示';"`           // 标示
	BaseCurrency  string  `gorm:"index;size:250;column:base_currency;comment:'基础货币/大写字母';"`  // 基础货币
	QuoteCurrency string  `gorm:"index;size:250;column:quote_currency;comment:'报价货币/大写字母';"` // 基础货币
	Price         float64 `gorm:"column:price;default:0;comment:'金额';"`                      // 金额
}

// Table 表
func GetQuoteTable() string {
	return viper.GetString("mysql.prefix") + "quote"
}

// NewQuote NewQuote
func NewQuote() *Quote {
	return &Quote{}
}

// CreateQuote 创建汇率
func (q *Quote) CreateQuote(o *gorm.DB) error {
	return o.Create(q).Error
}

// IsExistQuote 判断汇率是否存在
func (q *Quote) IsExistQuote(o *gorm.DB) error {
	return o.Table(GetQuoteTable()).Where("code = ?", q.Code).Find(q).Error
}

// GetOrderIdByInfo 根据id获取信息
func (q *Quote) GetOrderIdByInfo(o *gorm.DB) error {
	return o.Table(GetQuoteTable()).Where("id = ?", q.ID).First(q).Error
}

// GetQuoteCurrencyByList 根据QuoteCurrency 获取列表
func (q *Quote) GetQuoteCurrencyByList(o *gorm.DB) ([]resp.CurrencyQuoteInfoResp, error) {
	rows, err := o.Model(q).
		Where("quote_currency = ?", q.QuoteCurrency).
		Select("id,code,base_currency,quote_currency,price,updated_at").Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.CurrencyQuoteInfoResp
			timer time.Time
		)
		data := make([]resp.CurrencyQuoteInfoResp, 0)

		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				timer, _ = time.Parse("2006-01-02T15:04:05+08:00", item.UpdatedAt)
				item.UpdatedAt = timer.Format("2006-01-02 15:04:05")
				data = append(data, item)
			}
		}
		return data, err
	}
	return nil, err
}

// UpdateQuotePrice 更新金额
func (q *Quote) UpdateQuotePrice(o *gorm.DB) error {
	if err := o.Table(GetQuoteTable()).Where("id = ? ", q.ID).Update(map[string]interface{}{
		"price":      q.Price,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetAllPageList 获取分页列表
func (q *Quote) GetAllPageList(o *gorm.DB, page, pageSize int, start_time, end_timer int, keyword string) (resp.CurrencyQuoteListResp, error) {
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

	sql = fmt.Sprintf("%s ORDER BY id desc LIMIT %d,%d", sql, (page-1)*pageSize, pageSize)

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
