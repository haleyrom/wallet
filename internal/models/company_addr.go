package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// CompanyAddr  公司地址
type CompanyAddr struct {
	gorm.Model
	Code         string `gorm:"column:code;comment:'code'"` // 账户名称
	BlockChainId uint   `gorm:"column:block_chain_id;default:0;comment:'链id'"`
	Symbol       string `gorm:"size:255;column:symbol;comment:'代币代号';"`                    // 代币代号
	Type         string `gorm:"size:255;column:type;comment:'币类型';"`                       // 标识 coin,toke
	Address      string `gorm:"size:255;column:address;comment:'地址'"`                      // 地址
	Status       int8   `gorm:"size:3;column:status;default:0;commit:'状态(0开启,1:停用,2:删除)'"` // 状态：0开启;1:停用;2:删除
}

// GetCompanyAddrTable 公司地址
func GetCompanyAddrTable() string {
	return viper.GetString("mysql.prefix") + "company_addr"
}

// NewCompanyAddr 初始化
func NewCompanyAddr() *CompanyAddr {
	return &CompanyAddr{}
}

// GetList 获取企业充值列表
func (c *CompanyAddr) GetList(o *gorm.DB, page, pageSize, start_time, end_timer int, keyword string) (resp.CompanyAddrListResp, error) {
	data := resp.CompanyAddrListResp{}
	sql := fmt.Sprintf("SELECT id,symbol,type,address,status,updated_at FROM %s where status < %d AND code = '%s' ", GetCompanyAddrTable(), vStatusRm, c.Code)
	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s where status < %d AND code = '%s' ", GetCompanyAddrTable(), vStatusRm, c.Code)

	if start_time > 0 && end_timer > 0 {
		sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(updated_at) >= %d AND UNIX_TIMESTAMP(updated_at) <= %d ", sql, start_time, end_timer)
		count_sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(updated_at) >= %d AND UNIX_TIMESTAMP(updated_at) <= %d ", count_sql, start_time, end_timer)
	}

	if len(keyword) > 0 {
		sql = fmt.Sprintf("%s AND symbol like '%s'", sql, "%"+keyword+"%")
		count_sql = fmt.Sprintf("%s AND symbol like '%s'", count_sql, "%"+keyword+"%")
	}

	sql = fmt.Sprintf("%s ORDER BY id DESC LIMIT %d,%d", sql, (page-1)*pageSize, pageSize)

	rows, err := o.Raw(sql).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.CompanyAddrInfoResp
			timer time.Time
		)

		data.Items = make([]resp.CompanyAddrInfoResp, 0)
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

// CreateCompanyAddr  创建公司地址
func (c *CompanyAddr) CreateCompanyAddr(o *gorm.DB) error {
	c.Status = vStatusOk
	return o.Create(c).Error
}

// UpdateAddrStatus 更新状态
func (c *CompanyAddr) UpdateAddrStatus(o *gorm.DB) error {
	if err := o.Model(c).Where("id = ?  and status < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"status":     c.Status,
	}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAddr 更新地址
func (c *CompanyAddr) UpdateAddr(o *gorm.DB) error {
	if err := o.Model(c).Where("id = ?  and status < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"address":    c.Address,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetOrderSymbolByAddress
func (c *CompanyAddr) GetOrderSymbolByAddress(o *gorm.DB) (string, error) {
	var address string
	err := o.Table(GetCompanyAddrTable()).Where("symbol = ? and status = ? and code = ?", c.Symbol, vStatusOk, c.Code).Select("address").Order("id desc").Limit(1).Row().Scan(&address)
	return address, err
}
