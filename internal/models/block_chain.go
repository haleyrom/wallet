package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

// BlockChain 币种
type BlockChain struct {
	gorm.Model
	ChainCode string `gorm:"size:255;column:chain_code;commit:'链code'"`
	Name      string `gorm:"size:255;column:name;comment:'链名称';"`                       // 币种名称
	Status    int8   `gorm:"size:3;column:status;default:0;commit:'状态(0开启,1:停用,2:删除)'"` // 状态：0开启;1:停用;2:删除

}

// GetBlockChain 表
func GetBlockChain() string {
	return viper.GetString("mysql.prefix") + "block_chain"
}

// NewBlockChain 初始化
func NewBlockChain() *BlockChain {
	return &BlockChain{}
}

// GetIdAll 获取全部id
func (b *BlockChain) GetAll(o *gorm.DB) ([]resp.ReadBlockChainListResp, error) {
	data := make([]resp.ReadBlockChainListResp, 0)
	rows, err := o.Table(GetBlockChain()).Where("status < ?", vStatusRm).Select("id,chain_code,name").Rows()
	defer rows.Close()

	if err == nil {
		var item resp.ReadBlockChainListResp
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				data = append(data, item)
			}
		}
	}
	return data, err
}

// IsExistCurrency 判断是否存在
func (b *BlockChain) IsExistBlockChain(o *gorm.DB) error {
	if err := o.Table(GetBlockChain()).Where("id = ? and status < ?", b.ID, vStatusRm).First(&b).Error; err == nil {
		return nil
	}
	return fmt.Errorf("%s", "block_chain not exist")
}

// UpdateCurrency  更新币种
func (b *BlockChain) UpdateBlockChain(o *gorm.DB) error {
	if err := o.Model(b).Where("id = ?", b.ID).Update(b).Error; err != nil {
		return err
	}
	return nil
}

// CreateBlockChain  创建链
func (b *BlockChain) CreateBlockChain(o *gorm.DB) error {
	return o.Create(b).Error
}

// RmChain 删除链区
func (b *BlockChain) RmChain(o *gorm.DB) error {
	if err := o.Model(b).Where("id = ? and status < ?", b.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
		"status":     vStatusRm,
	}).Error; err != nil {
		return err
	}
	return nil
}
