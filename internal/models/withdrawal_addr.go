package models

import (
	"database/sql"
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

// WithdrawalAddr  提现地址
type WithdrawalAddr struct {
	gorm.Model
	Uid           uint   `gorm:"column:uid;default:0;comment:'用户id';"`
	BlockChainId  uint   `gorm:"column:block_chain_id;default:0;comment:'链id'"`
	Name          string `gorm:"column:name;comment:'名称'"`
	CurrencyId    uint   `gorm:"column:currency_id;default:0;comment:'币种id';"` // 账户id
	Address       string `gorm:"size:255;column:address;comment:'地址'"`
	Status        int8   `gorm:"size:3;column:status;default:0;commit:'状态(0开启,1:停用,2:删除)'"`   // 状态：0开启;1:停用;2:删除
	AddressSource int8   `gorm:"size:3;column:address_source;default:0;commit:'来源0未知1本站2外站'"` // 来源0未知1本站2外站
}

const (
	// WithdrawalAddrUnknown 未知
	WithdrawalAddrUnknown int8 = 0 + iota
	// WithdrawalAddrLocal 本站
	WithdrawalAddrLocal
	// WithdrawalAddrBack 站外
	WithdrawalAddrBack
)

// GetWithdrawalAddrTable 提现地址
func GetWithdrawalAddrTable() string {
	return viper.GetString("mysql.prefix") + "withdrawal_addr"
}

// NewWithdrawalAddr 初始化
func NewWithdrawalAddr() *WithdrawalAddr {
	return &WithdrawalAddr{}
}

// GetAll 获取全部
func (w *WithdrawalAddr) GetAll(o *gorm.DB) ([]resp.WithdrawalAddrResp, error) {
	data := make([]resp.WithdrawalAddrResp, 0)
	var (
		rows *sql.Rows
		err  error
	)
	if w.BlockChainId > 0 || w.CurrencyId > 0 {
		rows, err = o.Raw(fmt.Sprintf("SELECT currency.name as currency_name,withdrawal.name,withdrawal.id as withdrawal_addr_id,withdrawal.address,chain.chain_code, chain.name as chain_name FROM %s withdrawal LEFT JOIN %s chain on chain.id = withdrawal.block_chain_id LEFT JOIN %s currency on currency.id = withdrawal.currency_id WHERE withdrawal.uid = ? AND withdrawal.status < ? AND chain.status < ? AND currency.status < ? AND withdrawal.block_chain_id = ? AND withdrawal.currency_id = ?", GetWithdrawalAddrTable(), GetBlockChain(), GetCurrencyTable()), w.Uid, vStatusRm, vStatusRm, vStatusRm, w.BlockChainId, w.CurrencyId).Rows()
	} else {
		rows, err = o.Raw(fmt.Sprintf("SELECT currency.name as currency_name,withdrawal.name,withdrawal.id as withdrawal_addr_id,withdrawal.address,chain.chain_code, chain.name as chain_name FROM %s withdrawal LEFT JOIN %s chain on chain.id = withdrawal.block_chain_id LEFT JOIN %s currency on currency.id = withdrawal.currency_id WHERE withdrawal.uid = ? AND withdrawal.status < ? AND chain.status < ? AND currency.status < ?", GetWithdrawalAddrTable(), GetBlockChain(), GetCurrencyTable()), w.Uid, vStatusRm, vStatusRm, vStatusRm).Rows()
	}

	defer rows.Close()

	if err == nil {
		var item resp.WithdrawalAddrResp
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				data = append(data, item)
			}
		}
	}
	return data, err
}

// CreateWithdrawalAddr  创建提现地址
func (w *WithdrawalAddr) CreateWithdrawalAddr(o *gorm.DB) error {
	return o.Create(w).Error
}

// UpdateWithdrawalAddr  更新提现地址
func (w *WithdrawalAddr) UpdateWithdrawalAddr(o *gorm.DB) error {
	w.UpdatedAt = time.Now()
	if err := o.Model(w).Where("id = ? and uid = ? and status < ?", w.ID, w.Uid, vStatusRm).Update(w).Error; err != nil {
		return err
	}
	return nil
}

// RmWithdrawalAddr 删除提现地址
func (w *WithdrawalAddr) RmWithdrawalAddr(o *gorm.DB) error {
	if err := o.Model(w).Where("id = ? and uid = ?  and status < ?", w.ID, w.Uid, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
		"status":     vStatusRm,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetInfo 获取信息
func (w *WithdrawalAddr) GetInfo(o *gorm.DB) error {
	return o.Table(GetWithdrawalAddrTable()).Where("id = ?", w.ID).Find(w).Error
}

// GetUsableInfo 获取可用信息
func (w *WithdrawalAddr) GetUsableInfo(o *gorm.DB) error {
	return o.Table(GetWithdrawalAddrTable()).Where("id = ? and status = ?", w.ID, vStatusOk).Find(w).Error
}
