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

// DepositDetail  充值记录
type DepositDetail struct {
	gorm.Model
	Uid             uint    `gorm:"column:uid;comment:'用户id'"`                                   // 用户id
	CoinId          uint    `gorm:"column:coin_id;default:0comment:'币种id'"`                      // 币种id
	CurrencyId      uint    `gorm:"column:currency_id;default:0comment:'货币id'"`                  //
	Address         string  `gorm:"size:255;column:address;index;comment:'地址'"`                  // 地址
	Value           float64 `gorm:"column:value;default:0;comment:'充值金额'"`                       // 充值金额
	BlockNumber     int     `gorm:"column:block_number;default:0;comment:'充值区块高度'"`              // 充值区块高度
	BlockCount      int     `gorm:"column:block_count;default:0;comment:'区块确认数'"`                // 当区块确认数达到最小确认时，入账。
	TransactionHash string  `gorm:"size:255;column:transaction_hash;index;comment:'事务hash'"`     // 事务hash
	Symbol          string  `gorm:"size:255;column:symbol;comment:'代币代号';"`                      // 代币代号
	Type            string  `gorm:"size:255;column:type;comment:'币类型';"`                         // 标识 coin,token
	Status          int8    `gorm:"size:3;column:status;default:0;comment:'入账状态'"`               // 入账状态：0未入账,1已入账,2审核不通过。
	ContractAddress string  `gorm:"size:200;column:contract_address;comment:'合约地址'"`             // 合约地址
	Source          int8    `gorm:"size:3;column:source;default:0;comment:'来源(0充值，1手充)'"`        // 来源(0充值1手充)
	FinancialStatus int8    `gorm:"column:financial_status;index;default:0;comment:'财务状态'"`      // 财务状态:0 待审核1：通过2：不通过
	FinancialId     uint    `gorm:"column:financial_id;index;default:0;comment:'财务id'"`          // 财务id
	FinancialRemark string  `gorm:"column:financial_remark;comment:'财务备注'"`                      // 财务备注
	Deleted         int8    `gorm:"column:deleted;default:0;comment:'删除状态（0不删除1删除）'"`            // 删除状态（0不删除1删除）
	Md5Keys         string  `gorm:"column:md5key;default:0;unique_index;comment:'md5key'"`       //
	AddressSource   int8    `gorm:"size:3;column:address_source;default:0;commit:'来源0未知1本站2外站'"` // 来源0未知1本站2外站
}

const (
	// DepositStatusNotBooked 未入账
	DepositStatusNotBooked int8 = 0
	// DepositStatusNOtBooked  已入账
	DepositStatusBooked int8 = 1
	// DepositStatusNotDeleted 不删除
	DepositStatusNotDeleted int8 = 0
	// DepositStatusYesDeleted 删除
	DepositStatusYesDeleted int8 = 1
	// DepositSourceRecharge 充值
	DepositSourceRecharge int8 = 0
	// DepositSourceAdmin 后台充值
	DepositSourceAdmin int8 = 1
)

// Table 表
func GetDepositDetailTable() string {
	return viper.GetString("mysql.prefix") + "deposit_detail"
}

// NewDepositDetail 初始化
func NewDepositDetail() *DepositDetail {
	return &DepositDetail{}
}

// CreateDepositDetail 创建充值记录
func (d *DepositDetail) CreateDepositDetail(o *gorm.DB) error {
	return o.Create(d).Error
}

// IsDepositDetail 判断
func (d *DepositDetail) IsTransactionHash(o *gorm.DB) error {
	return o.Table(GetDepositDetailTable()).Where("transaction_hash = ? and deleted = ?", d.TransactionHash, DepositStatusNotDeleted).Order("status desc").Find(d).Error
}

// IsKey 判断
func (d *DepositDetail) IsKey(o *gorm.DB) error {
	return o.Table(GetDepositDetailTable()).Where("md5key = ? and deleted = ?", d.Md5Keys, DepositStatusNotDeleted).First(d).Error
}

// IsDepositBooked 判断是否充值
func (d *DepositDetail) IsDepositBooked(o *gorm.DB) bool {
	var num int
	_ = o.Table(GetDepositDetailTable()).Where("transaction_hash = ? and status = ? and deleted = ?", d.TransactionHash, DepositStatusBooked, DepositStatusNotDeleted).Count(&num)
	if num == 0 {
		return false
	}
	return true
}

// UpdateBlockCount 更新确认数
func (d *DepositDetail) UpdateBlockCount(o *gorm.DB) error {
	if err := o.Table(GetDepositDetailTable()).Where("transaction_hash = ? ", d.TransactionHash).Update(map[string]interface{}{
		"status":      d.Status,
		"updated_at":  time.Now(),
		"block_count": d.BlockCount,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetPageList 获取分页列表
func (d *DepositDetail) GetPageList(o *gorm.DB, page, pageSize int) (resp.ReadDepositDetailResp, error) {
	data := resp.ReadDepositDetailResp{}
	rows, err := o.Raw(fmt.Sprintf("SELECT TRUNCATE(address,value,6) as value,symbol,status,type,updated_at,block_count FROM %s  where uid = ? and deleted = ? and source = ? ORDER BY id desc LIMIT ?,?", GetDepositDetailTable()), d.Uid, DepositStatusNotDeleted, (page-1)*pageSize, pageSize, DepositSourceRecharge).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.ReadDepositDetailInfoResp
			timer time.Time
		)

		data.Items = make([]resp.ReadDepositDetailInfoResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data.Items = append(data.Items, item)
			}
		}

		o.Table(GetDepositDetailTable()).Where("uid = ? and deleted = ? and source = ?", d.Uid, DepositStatusNotDeleted, DepositSourceRecharge).Count(&data.Page.Count)
		data.Page.PageSize = len(data.Items)
		data.Page.CurrentPage = page
		data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	}
	return data, err
}

// GetAllPageList 获取提现分页也表
func (d *DepositDetail) GetAllPageList(o *gorm.DB, page, pageSize int, endTime, startTime int, key string) (resp.ReadAllDepositDetailResp, error) {
	data := resp.ReadAllDepositDetailResp{}
	if endTime == 0 {
		endTime = 10000000000000
	}
	count_sql := fmt.Sprintf("select count(*) as num from %s a left join %s b on a.uid = b.id where UNIX_TIMESTAMP(a.updated_at) >= %d and UNIX_TIMESTAMP(a.updated_at) <= %d and a.deleted = %d", GetDepositDetailTable(), GetUserTable(), startTime, endTime, DepositStatusNotDeleted)
	sql := fmt.Sprintf("select a.id as order_id,a.uid,b.name,a.symbol,a.type,TRUNCATE(a.value,6) as value,a.transaction_hash,a.status,a.updated_at,a.source from %s a left join %s b on a.uid = b.id where UNIX_TIMESTAMP(a.updated_at) >= %d and UNIX_TIMESTAMP(a.updated_at) <= %d  and a.deleted = %d ", GetDepositDetailTable(), GetUserTable(), startTime, endTime, DepositStatusNotDeleted)
	if key != "" {
		count_sql += "and  ((b.name like '%" + key + "%') or (b.email like '%" + key + "%') or (b.uid like '%" + key + "%') ) "
		sql += "and  ((b.name like '%" + key + "%') or (b.email like '%" + key + "%') or (b.uid like '%" + key + "%') ) "
	}
	sql = sql + fmt.Sprintf("order by a.id desc limit %d , %d", (page-1)*pageSize, pageSize)
	rows, err := o.Raw(sql).Rows()
	defer rows.Close()
	var (
		item  resp.ReadAllDepositDetailInfoResp
		timer time.Time
	)
	data.Items = make([]resp.ReadAllDepositDetailInfoResp, 0)
	if err == nil {
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data.Items = append(data.Items, item)
			}
		}
		_ = o.Raw(count_sql).Row().Scan(&data.Page.Count)
		data.Page.PageSize = len(data.Items)
		data.Page.CurrentPage = page
		data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	}
	return data, err
}

// GetAllRechargePageList 获取提现分页也表
func (d *DepositDetail) GetAllRechargePageList(o *gorm.DB, page, pageSize int, endTime, startTime int, key string) (resp.ReadAllDepositDetailResp, error) {
	data := resp.ReadAllDepositDetailResp{}
	if endTime == 0 {
		endTime = 10000000000000
	}
	count_sql := fmt.Sprintf("select count(*) as num from %s a left join %s b on a.uid = b.id where UNIX_TIMESTAMP(a.updated_at) >= %d and UNIX_TIMESTAMP(a.updated_at) <= %d and a.deleted = %d and a.source = %d and a.uid > 0 ", GetDepositDetailTable(), GetUserTable(), startTime, endTime, DepositStatusNotDeleted, DepositSourceAdmin)
	sql := fmt.Sprintf("select a.id as order_id,a.uid,b.name,a.symbol,a.type,TRUNCATE(a.value,6) as value,a.transaction_hash,a.status,a.updated_at,b.name,a.source from %s a left join %s b on a.uid = b.id where UNIX_TIMESTAMP(a.updated_at) >= %d and UNIX_TIMESTAMP(a.updated_at) <= %d  and a.deleted = %d  and a.source = %d and a.uid > 0 ", GetDepositDetailTable(), GetUserTable(), startTime, endTime, DepositStatusNotDeleted, DepositSourceAdmin)
	if key != "" {
		count_sql += "and ((b.name = '%" + key + "%') or (b.email = '%" + key + "%') or (b.uid = '%" + key + "%'))"
		sql += "and ((b.name = '%" + key + "%') or (b.email = '%" + key + "%') or (b.uid = '%" + key + "%'))"
	}
	sql = sql + fmt.Sprintf("order by a.id desc limit %d , %d", (page-1)*pageSize, pageSize)
	rows, err := o.Raw(sql).Rows()
	defer rows.Close()
	var (
		item  resp.ReadAllDepositDetailInfoResp
		timer time.Time
	)
	data.Items = make([]resp.ReadAllDepositDetailInfoResp, 0)
	if err == nil {
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data.Items = append(data.Items, item)
			}
		}
		_ = o.Raw(count_sql).Row().Scan(&data.Page.Count)
		data.Page.PageSize = len(data.Items)
		data.Page.CurrentPage = page
		data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	}
	return data, err
}

// GetSourceAdminOrderIdByInfo 后台充值根据id获取info
func (d *DepositDetail) GetSourceAdminOrderIdByInfo(o *gorm.DB) error {
	return o.Table(GetDepositDetailTable()).Where("id = ? and deleted = ? and source = ?", d.ID, DepositStatusNotDeleted, DepositSourceAdmin).First(d).Error
}

// UpdateDeleted 更新删除状态
func (d *DepositDetail) UpdateDeleted(o *gorm.DB) error {
	timer := time.Now()
	if err := o.Table(GetDepositDetailTable()).Where("id = ? ", d.ID).Update(map[string]interface{}{
		"updated_at": timer,
		"deleted_at": timer,
		"deleted":    DepositStatusYesDeleted,
	}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateFinancial 更新财务状态
func (d *DepositDetail) UpdateFinancial(o *gorm.DB) error {
	var err error
	timer := time.Now()
	data := map[string]interface{}{
		"updated_at":       timer,
		"financial_id":     d.FinancialId,
		"financial_status": d.FinancialStatus,
		"status":           d.Status,
	}
	if d.Status == DepositStatusBooked {
		data["status"] = DepositStatusBooked
	}
	if err := o.Table(GetDepositDetailTable()).Where("id = ? ", d.ID).Update(data).Error; err != nil {
		return err
	}
	return err
}
