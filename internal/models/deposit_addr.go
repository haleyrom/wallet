package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// DepositAddr  充值地址
type DepositAddr struct {
	gorm.Model
	Uid          uint   `gorm:"column:uid;comment:'用户id'"`                                 // 用户id
	AccountId    uint   `gorm:"column:account_id;default:0;comment:'账户id';"`               // 账户id
	BlockChainId uint   `gorm:"column:block_chain_id;default:0;comment:'链id'"`             // 币种id
	Address      string `gorm:"size:255;column:address;comment:'充值地址'"`                    // 充值地址
	OrderId      string `gorm:"size:255;column:order_id;comment:'订单号'"`                    // 生成该地址的订单号
	Status       int8   `gorm:"size:3;column:status;default:0;commit:'状态(0开启,1:停用,2:删除)'"` // 状态：0开启;1:停用;2:删除

}

// GetDepositAddrTable 充值地址
func GetDepositAddrTable() string {
	return viper.GetString("mysql.prefix") + "deposit_addr"
}

// NewDepositAddr 初始化
func NewDepositAddr() *DepositAddr {
	return &DepositAddr{}
}

// CreateDepositAddr  创建钱包地址
func (d *DepositAddr) CreateDepositAddr(o *gorm.DB) error {
	return o.Create(d).Error
}

// GetAll 获取全部
func (d *DepositAddr) GetAll(o *gorm.DB) ([]resp.ReadDepositAddrListResp, error) {
	data := make([]resp.ReadDepositAddrListResp, 0)
	rows, err := o.Raw(fmt.Sprintf("SELECT chain.chain_code,chain.name as chain_name,depoist.address,depoist.block_chain_id,depoist.id as deposit_addr_id FROM %s depoist LEFT JOIN %s chain on chain.id = depoist.block_chain_id WHERE depoist.status < ? and chain.status < ?", GetDepositAddrTable(), GetBlockChain()), vStatusRm, vStatusRm).Rows()
	defer rows.Close()

	if err == nil {
		var item resp.ReadDepositAddrListResp
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				data = append(data, item)
			}
		}
	}
	return data, err
}

// IsExistAddress 判断是否存在
func (d *DepositAddr) IsExistAddress(o *gorm.DB) error {
	if err := o.Model(d).Where("address = ? and status  < ?", d.Address, vStatusRm).Select("id").First(&d).Error; err == nil {
		return nil
	}
	return fmt.Errorf("%s", "address not exist")
}

// ReadWithdrawalAddr 读取
func (d *DepositAddr) ReadWithdrawalAddr(o *gorm.DB) (resp.ReadDepositAddrResp, error) {
	var data resp.ReadDepositAddrResp
	row := o.Model(d).Where("block_chain_id = ? and uid = ?", d.BlockChainId, d.Uid).Select("id as deposit_addr_id,block_chain_id,address").Row() // (*sql.Row)
	_ = row.Scan(&data.DepositAddrId, &data.BlockChainId, &data.Address)
	if data.DepositAddrId == 0 {
		return data, fmt.Errorf("%s", "query is not rows")
	}
	return data, nil
}

// GetAddressByInfo 根据地址获取信息
func (d *DepositAddr) GetAddressByInfo(o *gorm.DB) error {
	if err := o.Table(GetDepositAddrTable()).Where("address = ? ", d.Address).First(&d).Error; err == nil && d.Uid > 0 {
		return nil
	}
	return fmt.Errorf("%s", "address not exist")
}
