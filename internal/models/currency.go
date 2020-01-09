package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

// Currency 代币
type Currency struct {
	gorm.Model
	Symbol      string  `gorm:"size:255;column:symbol;comment:'币种标识';"`                     // 代币标识
	Name        string  `gorm:"size:255;column:name;comment:'币种名称';"`                       // 代币名称
	Decimals    int     `gorm:"column:decimals;default:0;comment:'币种小数点';"`                 // 代币小数点
	Status      int8    `gorm:"size(3);column:status;default:0;commit:'状态(0开启,1:停用,2:删除)'"` // 状态：0开启;1:停用;2:删除
	MinPayMoney float64 `gorm:"column:min_pay_money;default:0;commit:'最小支付金额'"`             // 最小支付金额
}

const (
	// vStatusOk 状态
	vStatusOk int8 = 0 + iota
	// vStatusStop 停止
	vStatusStop
	// vStatusRm 删除
	vStatusRm
)

// Table 表
func GetCurrencyTable() string {
	return viper.GetString("mysql.prefix") + "currency"
}

// NewAccount 初始化
func NewCurrency() *Currency {
	return &Currency{}
}

// GetIdAll 获取全部id
func (c *Currency) GetIdAll(o *gorm.DB) ([]uint, error) {
	data := make([]uint, 0)
	rows, err := o.Table(GetCurrencyTable()).Select("id").Rows()
	defer rows.Close()

	if err == nil {
		var id uint
		for rows.Next() {
			if err = rows.Scan(&id); err == nil {
				data = append(data, id)
			}
		}
	}
	return data, err
}

// GetIdAll 获取全部id
func (c *Currency) GetAll(o *gorm.DB) ([]resp.ReadCurrencyListResp, error) {
	data := make([]resp.ReadCurrencyListResp, 0)
	rows, err := o.Table(GetCurrencyTable()).Where("status < ?", vStatusRm).Select("id as currency_id,symbol,name,decimals,status,updated_at").Rows()
	defer rows.Close()

	if err == nil {
		var (
			timer time.Time
			item  resp.ReadCurrencyListResp
		)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				timer, _ = time.Parse("2006-01-02T15:04:05+08:00Z", item.UpdatedAt)
				item.UpdatedAt = timer.Format("2006-01-02 15:04:05")
				data = append(data, item)
			}
		}
	}
	return data, err
}

// IsExistCurrency 判断是否存在
func (c *Currency) IsExistCurrency(o *gorm.DB) error {
	if err := o.Table(GetCurrencyTable()).Where("id = ? and status < ?", c.ID, vStatusRm).Select("id").First(&c).Error; err == nil {
		return nil
	}
	return fmt.Errorf("%s", "currency not exist")
}

// UpdateCurrency  更新币种
func (c *Currency) UpdateCurrency(o *gorm.DB) error {
	if err := o.Table(GetCurrencyTable()).Where("id = ?", c.ID).Update(c).Error; err != nil {
		return err
	}
	return nil
}

// CreateCurrency  创建币种
func (c *Currency) CreateCurrency(o *gorm.DB) error {
	return o.Create(c).Error
}

// UpdateCurrencyStatus 停止币种状态
func (c *Currency) UpdateCurrencyStatus(o *gorm.DB) error {
	if err := o.Table(GetCurrencyTable()).Where("id = ? and status  < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"status":     c.Status,
	}).Error; err != nil {
		return err
	}
	return nil
}

// RmCurrency 删除币种
func (c *Currency) RmCurrency(o *gorm.DB) error {
	if err := o.Table(GetCurrencyTable()).Where("id = ? and status < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
		"status":     vStatusRm,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetSymbolById 获取id
func (c *Currency) GetSymbolById(o *gorm.DB) error {
	if err := o.Table(GetCurrencyTable()).Where("symbol = ? and status < ?", c.Symbol, vStatusStop).Select("id,min_pay_money").First(&c).Error; err == nil {
		return nil
	}
	return fmt.Errorf("%s", "currency not exist")
}
