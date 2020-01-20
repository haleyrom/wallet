package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

// Coin 代币
type Coin struct {
	gorm.Model
	CurrencyId        uint    `gorm:"column:currency_id;comment:'币种id'"`                               // 币种id
	Symbol            string  `gorm:"size:255;column:symbol;comment:'代币代号';"`                          // 代币代号
	Name              string  `gorm:"size:255;column:name;comment:'代币名称';"`                            // 代币名称
	BlockChainId      uint    `gorm:"column:block_chain_id;default:0;comment:'区块链id';"`                // 区块链名称
	Type              string  `gorm:"size:255;column:type;comment:'币类型';"`                             // 标识 coin,token
	ConfirmCount      int     `gorm:"column:confirm_count;default:0;comment:'确认数'"`                    // 充值入帐的区块链确认数
	MinDeposit        float64 `gorm:"column:min_deposit;default:0;comment:'最小充值金额'"`                   // 最小充值金额，小于该金额不入账
	MinWithdrawal     float64 `gorm:"column:min_withdrawal;default:0;comment:'最小提现金额'"`                // 小于该金额不能提
	WithdrawalFee     float64 `gorm:"column:withdrawal_fee;default:0;comment:'提现手续费'"`                 // 提现手续费
	WithdrawalFeeType string  `gorm:"size:255;column:withdrawal_fee_type;comment:'手续费类型'"`             // 手续费类型 fixed 按百分百比,percent 固定收取
	ContractAddress   string  `gorm:"size:255;column:contract_address;comment:'合约地址'"`                 // 合约地址:如该是type=token，这里必须输入
	Abi               string  `gorm:"size:255;column:abi;comment:'字节数'"`                               // 字节长度
	Status            int8    `gorm:"size(3);column:status;default:0;commit:'状态(0开启,1:停用,2:删除)'"`      // 状态：0开启;1:停用;2:删除
	WithdrawalStatus  int8    `gorm:"size(3);column:withdrawal_status;default:0;commit:'充值状态-0开启1停用'"` // 状态：0开启;1:停用;
	DepositStatus     int8    `gorm:"size(3);column:deposit_status;default:0;commit:'提笔状态-0开启1停用'"`    // 状态：0开启;1:停用;
	CustomerStatus    int8    `gorm:"size(3);column:customer_status;index;default:0;comment:'客服状态'"`   // 客服状态:0 必须1：不必须
	FinancialStatus   int8    `gorm:"size(3);column:financial_status;index;default:0;comment:'财务状态'"`  // 财务状态:0 必须1：不必须
	Precision         int     `gorm:"column:precision;default:6;comment:'小数点精度6位'"`                    // 小数点精度默认6位
}

// GetCoinTable 表
func GetCoinTable() string {
	return viper.GetString("mysql.prefix") + "coin"
}

// NewCoin 初始化
func NewCoin() *Coin {
	return &Coin{}
}

// GetAll 获取全部
func (c *Coin) GetAll(o *gorm.DB) ([]resp.ReadCoinListResp, error) {
	data := make([]resp.ReadCoinListResp, 0)
	rows, err := o.Raw(fmt.Sprintf("SELECT coin.currency_id,currency.symbol as currency_symbol,currency.name as currency_name,currency.decimals as currency_decimals,coin.id as coin_id,coin.symbol,coin.type,coin.status,coin.name,block_chain_id,chain.chain_code,chain.name as chain_name,coin.type,confirm_count,min_deposit,min_withdrawal,withdrawal_fee,withdrawal_fee_type,contract_address,coin.updated_at,coin.withdrawal_status,coin.deposit_status,coin.customer_status,coin.financial_status FROM %s coin left JOIN %s chain on chain.id = coin.block_chain_id LEFT JOIN %s currency ON currency.id = coin.currency_id WHERE coin.status < ?", GetCoinTable(), GetBlockChain(), GetCurrencyTable()), vStatusRm).Rows()
	defer rows.Close()

	if err == nil {
		var (
			timer time.Time
			item  resp.ReadCoinListResp
		)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data = append(data, item)
			}
		}
	}
	return data, err
}

// IsExistCoin 判断是否存在
func (c *Coin) IsExistCoin(o *gorm.DB) error {
	if err := o.Model(c).Where("id = ? and status < ?", c.ID, vStatusRm).Select("id").First(&c).Error; err == nil {
		return nil
	}
	return fmt.Errorf("%s", "Coin not exist")
}

// UpdateCoin  更新代币
func (c *Coin) UpdateCoin(o *gorm.DB) error {
	c.UpdatedAt = time.Now()
	fmt.Println(c)
	if err := o.Table(GetCoinTable()).Where("id = ? and status < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"symbol":              c.Symbol,
		"name":                c.Name,
		"currency_id":         c.CurrencyId,
		"block_chain_id":      c.BlockChainId,
		"type":                c.Type,
		"confirm_count":       c.ConfirmCount,
		"min_deposit":         c.MinDeposit,
		"min_withdrawal":      c.MinWithdrawal,
		"withdrawal_fee":      c.WithdrawalFee,
		"withdrawal_fee_type": c.WithdrawalFeeType,
		"contract_address":    c.ContractAddress,
		"abi":                 c.Abi,
		"withdrawal_status":   c.WithdrawalStatus,
		"deposit_status":      c.DepositStatus,
		"customer_status":     c.CustomerStatus,
		"financial_status":    c.FinancialStatus,
		"updated_at":          time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// CreateCoin  创建链
func (c *Coin) CreateCoin(o *gorm.DB) error {
	return o.Create(c).Error
}

// RmChain 删除链区
func (c *Coin) RmChain(o *gorm.DB) error {
	if err := o.Model(c).Where("id = ? and status < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
		"status":     vStatusRm,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetInfo  获取消息
func (c *Coin) GetInfo(o *gorm.DB) (resp.ReadCoinInfoResp, error) {
	var data resp.ReadCoinInfoResp
	row := o.Table(GetCoinTable()).Where("id = ? and status < ?", c.ID, vStatusRm).Select("id as coin_id,currency_id,symbol,name,block_chain_id,type,confirm_count,min_deposit,min_withdrawal,withdrawal_fee,withdrawal_fee_type,contract_address,abi,status,withdrawal_status,deposit_status,customer_status,financial_status").Row()
	_ = row.Scan(&data.CoinId, &data.CurrencyId, &data.Symbol, &data.Name, &data.BlockChainId, &data.Type, &data.ConfirmCount, &data.MinDeposit, &data.MinWithdrawal, &data.WithdrawalFee, &data.WithdrawalFeeType, &data.ContractAddress, &data.Abi, &data.Status, &data.WithdrawalStatus, &data.DepositStatus, &data.CustomerStatus, &data.FinancialStatus)
	if data.CoinId == 0 {
		return data, fmt.Errorf("%s", "Coin not exist")
	}
	return data, nil
}

// UpdateCoinStatus  更新代币状态
func (c *Coin) UpdateCoinStatus(o *gorm.DB) error {
	if err := o.Table(GetCoinTable()).Where("id = ? and status < ?", c.ID, vStatusRm).Update(map[string]interface{}{
		"updated_at": time.Now(),
		"status":     c.Status,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetDepositInfo 获取提现信息
func (c *Coin) GetDepositInfo(o *gorm.DB) (resp.ReadCoinDepositInfoResp, error) {
	var data resp.ReadCoinDepositInfoResp
	row := o.Table(GetCoinTable()).Where("id = ? and status < ? and deposit_status = ?", c.ID, vStatusStop, vStatusOk).Select("id as coin_id,currency_id,min_withdrawal,withdrawal_fee,withdrawal_fee_type,symbol,type,status,withdrawal_status,deposit_status,customer_status,financial_status").Row()

	_ = row.Scan(&data.CoinId, &data.CurrencyId, &data.MinWithdrawal, &data.WithdrawalFee, &data.WithdrawalFeeType, &data.Symbol, &data.Type, &data.Status, &data.WithdrawalStatus, &data.DepositStatus, &data.CustomerStatus, &data.FinancialStatus)

	if data.CoinId == 0 {
		return data, fmt.Errorf("%s", "Coin not exist")
	}
	return data, nil
}

// GetOrderSymbolByChain  根据symbol获取链
func (c *Coin) GetOrderSymbolByChain(o *gorm.DB) ([]resp.ReadOrderSymbolByChainResp, error) {
	data := make([]resp.ReadOrderSymbolByChainResp, 0)
	rows, err := o.Raw(fmt.Sprintf("SELECT chain.id,coin.id as coin_id,chain.chain_code,chain.name,coin.type FROM %s coin LEFT JOIN %s chain on chain.id = coin.block_chain_id WHERE symbol = ?", GetCoinTable(), GetBlockChain()), c.Symbol).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item resp.ReadOrderSymbolByChainResp
		)
		for rows.Next() {
			_ = o.ScanRows(rows, &item)
			data = append(data, item)
		}
	}
	return data, nil
}

// GetOrderSymbolByCoin  根据symbol获取coin
func (c *Coin) GetOrderSymbolByCoin(o *gorm.DB, uid uint) ([]resp.ReadOrderSymbolByCoinResp, error) {
	data := make([]resp.ReadOrderSymbolByCoinResp, 0)
	rows, err := o.Raw(fmt.Sprintf("SELECT chain.chain_code,coin.symbol,coin.type,addr.address,coin.min_deposit FROM %s coin LEFT JOIN %s addr on addr.block_chain_id = coin.block_chain_id LEFT JOIN %s chain ON chain.id = coin.block_chain_id WHERE coin.symbol = ? and addr.uid = ? and coin.deposit_status = ? ", GetCoinTable(), GetDepositAddrTable(), GetBlockChain()), c.Symbol, uid, 0).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item resp.ReadOrderSymbolByCoinResp
		)
		for rows.Next() {
			_ = o.ScanRows(rows, &item)
			data = append(data, item)
		}
	}
	return data, nil
}

// GetOrderSymbolTypeByCoin 获取symbol/type获取信息
func (c *Coin) GetOrderSymbolTypeByCoin(o *gorm.DB) error {
	return o.Table(GetCoinTable()).Where("symbol = ? and type = ?", c.Symbol, c.Type).Find(c).Error
}

// GetConTractAddress 获取合约地址
func (c *Coin) GetConTractAddress(o *gorm.DB) (string, error) {
	var data string
	err := o.Table(GetCoinTable()).Where("symbol= ?", c.Symbol).Select("contract_address").Row().Scan(&data)
	return data, err
}

// GetOrderSymbolByInfo 根据symbol获取信息
func (c *Coin) GetOrderSymbolByInfo(o *gorm.DB) error {
	return o.Table(GetCoinTable()).Where("symbol = ?", c.Symbol).Find(c).Error
}

// GetOrderSymbolByInfo 根据symbol获取信息
func (c *Coin) GetOrderCurrencyIdByInfo(o *gorm.DB) error {
	return o.Table(GetCoinTable()).Where("currency_id = ?", c.CurrencyId).Find(c).Error
}

// GetOrderChainIdByInfo 根据coin获取信息
func (c *Coin) GetOrderChainIdByInfo(o *gorm.DB) error {
	return o.Table(GetCoinTable()).Where("block_chain_id = ? and currency_id = ?", c.BlockChainId, c.CurrencyId).Find(c).Error
}
