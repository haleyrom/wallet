package models

import (
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// BlockDetail 冻结明细
type BlockDetail struct {
	gorm.Model
	Uid         uint    `gorm:"column:uid;default:0;comment:'用户id'"`           // 用户id
	AccountId   uint    `gorm:"column:account_id;default:0;comment:'账本id';"`   // 账本id
	Balance     float64 `gorm:"column:balance;default:0;comment:'本期余额';"`      // 本期余额
	LastBalance float64 `gorm:"column:last_balance;default:0;comment:'上期余额';"` // 上期余额
	Income      float64 `gorm:"column:income;default:0;comment:'本期余额';"`       // 本期收入
	Spend       float64 `gorm:"column:spend;default:0;comment:'本期余额';"`        // 上期支出
	Type        int8    `gorm:"size(3);column:type;default:0;comment:'明细类型'"`  // 明细类型
}

// GetBlockDetailTable 获取表名
func GetBlockDetailTable() string {
	return viper.GetString("mysql.prefix") + "block_detail"
}

// NewBlockDetail 初始化
func NewBlockDetail() *BlockDetail {
	return &BlockDetail{}
}

// CreateBlockDetail  创建账本明细
func (b *BlockDetail) CreateBlockDetail(o *gorm.DB) error {
	return o.Create(b).Error
}
