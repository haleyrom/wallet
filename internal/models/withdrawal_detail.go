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

// WithdrawalDetail  提现记录
type WithdrawalDetail struct {
	gorm.Model
	Uid             uint    `gorm:"column:uid;default:0;comment:'用户id'"`                         // 用户id
	Address         string  `gorm:"size:255;column:address;comment:'地址'"`                        // 地址
	CoinId          uint    `gorm:"column:coin_id;default:0;comment:'代币id'"`                     // 代币id
	CurrencyId      uint    `gorm:"column:currency_id;default:0;comment:'货币id'"`                 // 货币id
	AccountId       uint    `gorm:"column:account_id;default:0;comment:'帐号id'"`                  // 钱包id
	Value           float64 `gorm:"column:value;default:0;comment:'提现金额'"`                       // 提现金额
	Symbol          string  `gorm:"size:255;column:symbol;comment:'代币代号';"`                      // 代币代号
	Type            string  `gorm:"size:255;column:type;comment:'币类型';"`                         // 标识 coin,token
	OrderId         string  `gorm:"size:255;column:order_id;comment:'订单号'"`                      // 生成该地址的订单号
	TransactionHash string  `gorm:"size:255;column:transaction_hash;comment:'事务hash'"`           // 事务
	Status          int8    `gorm:"size:3;column:status;default:0;index;comment:'状态'"`           // 0已提交,1待审核,2审核中,3通过,4不通过,5已完成,6取消,7提现失败
	Poundage        float64 `gorm:"column:poundage;default:0;comment:'手续费'"`                     // 手续费
	CustomerStatus  int8    `gorm:"column:customer_status;index;default:0;comment:'客服状态'"`       // 客服状态:0 待审核1：通过2：不通过
	FinancialStatus int8    `gorm:"column:financial_status;index;default:0;comment:'财务状态'"`      // 财务状态:0 待审核1：通过2：不通过
	BlockCount      int     `gorm:"column:block_count;default:0;comment:'确认数'"`                  // 充值入帐的区块链确认数
	CustomerId      uint    `gorm:"column:customer_id;index;default:0;comment:'客服id'"`           // 客服id
	FinancialId     uint    `gorm:"column:financial_id;index;default:0;comment:'财务id'"`          // 财务id
	Remark          string  `gorm:"size:200;column:remark;comment:'备注'"`                         // 备注
	RefundStatus    int8    `gorm:"size:4;column:refund_status;comment:'退款状态0不可退款1可退款2退款成功'"`    //退款状态（0不可退款1可退款2退款成功）
	CallbackStatus  string  `gorm:"column:callback_status;comment:'回调状态码'"`                      // 回调状态码
	CallbackJson    string  `gorm:"type:text;column:callback_json;comment:'回调json数据'"`           // 回调json数据
	AddressSource   int8    `gorm:"size:3;column:address_source;default:0;commit:'来源0未知1本站2外站'"` // 来源0未知1本站2外站
	Balance         float64 `gorm:"column:balance;default:0;comment:'当前可用余额';"`                  // 当前可用余额
	FromAddress     string  `gorm:"column:from_address;comment:'出金地址';"`                         // 出金地址
	BlockNumber     int     `gorm:"column:block_number;default:0;comment:'区块高度'"`                // 区块高度
}

const (
	// WithdrawalStatusSubmit 已提交
	WithdrawalStatusSubmit int8 = 0 + iota
	// WithdrawalStatusToAudit  待审核
	WithdrawalStatusToAudit
	// WithdrawalStatusInAudit  审核中
	WithdrawalStatusInAudit
	// WithdrawalStatusThrough  通过
	WithdrawalStatusThrough
	// WithdrawalStatusNoThrough  不通过
	WithdrawalStatusNoThrough
	// WithdrawalStatusOk 已完成
	WithdrawalStatusOk
	// WithdrawalStatusCancel 取消
	WithdrawalStatusCancel
)

const (
	// WithdrawalAudioStatusAwait 审核等待
	WithdrawalAudioStatusAwait int8 = 0 + iota
	// WithdrawalAudioStatusOk 通过审核
	WithdrawalAudioStatusOk
	// WithdrawalAudioStatusFailure 不通过
	WithdrawalAudioStatusFailure
)

// GetWithdrawalDetailTable 获取提现地址记录表名
func GetWithdrawalDetailTable() string {
	return viper.GetString("mysql.prefix") + "withdrawal_detail"

}

// NewWithdrawalDetail 初始化
func NewWithdrawalDetail() *WithdrawalDetail {
	return &WithdrawalDetail{}
}

// CreateWithdrawalDetail  创建提现记录
func (w *WithdrawalDetail) CreateWithdrawalDetail(o *gorm.DB) error {
	return o.Create(w).Error
}

// GetPageList 获取分页列表
func (w *WithdrawalDetail) GetPageList(o *gorm.DB, page, pageSize int) (resp.WithdrawalDetailListResp, error) {
	data := resp.WithdrawalDetailListResp{}
	rows, err := o.Raw(fmt.Sprintf("SELECT address,value,symbol,poundage,status,type,updated_at FROM %s  where uid = ? ORDER BY id desc LIMIT ?,?", GetWithdrawalDetailTable()), w.Uid, (page-1)*pageSize, pageSize).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.WithdrawalDetailResp
			timer time.Time
		)

		data.Items = make([]resp.WithdrawalDetailResp, 0)
		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				data.Items = append(data.Items, item)
			}
		}

		o.Table(GetWithdrawalDetailTable()).Where("uid = ?", w.Uid).Count(&data.Page.Count)
		data.Page.PageSize = len(data.Items)
		data.Page.CurrentPage = page
		data.Page.TotalPage = int(math.Ceil(float64(data.Page.Count) / float64(pageSize)))
	}
	return data, err
}

// GetAllPageList 获取全部分页列表
func (w *WithdrawalDetail) GetAllPageList(o *gorm.DB, page, pageSize, start_time, end_timer int, keyword string) (resp.WithdrawalDetailAllListResp, error) {
	data := resp.WithdrawalDetailAllListResp{}
	sql := fmt.Sprintf("select detail.remark,detail.order_id,detail.id,user.id as uid,user.name,user.email,detail.symbol,detail.financial_status,detail.customer_status,detail.value,detail.status,detail.updated_at,detail.address_source,detail.coin_id,detail.currency_id,detail.type,detail.address,detail.from_address,detail.balance,detail.callback_status,detail.callback_json FROM %s detail LEFT JOIN %s user on user.id = detail.uid WHERE detail.id > 0 ", GetWithdrawalDetailTable(), GetUserTable())
	count_sql := fmt.Sprintf("SELECT count(*) as num FROM %s detail LEFT JOIN %s user ON detail.uid = user.id where detail.id > 0 ", GetWithdrawalDetailTable(), GetUserTable())

	if start_time > 0 && end_timer > 0 {
		sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(detail.updated_at) >= %d AND UNIX_TIMESTAMP(detail.updated_at) <= %d ", sql, start_time, end_timer)
		count_sql = fmt.Sprintf("%s AND UNIX_TIMESTAMP(detail.updated_at) >= %d AND UNIX_TIMESTAMP(detail.updated_at) <= %d ", count_sql, start_time, end_timer)
	}

	if len(keyword) > 0 {
		sql = fmt.Sprintf("%s  AND ((user.name like '%s') or (user.email like '%s') or (user.uid like '%s')) ", sql, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		count_sql = fmt.Sprintf("%s  AND ((user.name like '%s') or (user.email like '%s') or (user.uid like '%s')) ", count_sql, "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	sql = fmt.Sprintf("%s ORDER BY detail.id DESC LIMIT %d,%d", sql, (page-1)*pageSize, pageSize)

	rows, err := o.Raw(sql).Rows()
	defer rows.Close()

	if err == nil {
		var (
			item  resp.WithdrawalDetailAdminResp
			timer time.Time
		)

		data.Items = make([]resp.WithdrawalDetailAdminResp, 0)
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

// ReadInfo 读取信息
func (w *WithdrawalDetail) ReadInfo(o *gorm.DB) (resp.AdminWithdrawalDetailResp, error) {
	var data resp.AdminWithdrawalDetailResp
	err := o.Raw(fmt.Sprintf("SELECT id,order_id,symbol,status,customer_status,financial_status,address,value,updated_at FROM %s WHERE id = ?", GetWithdrawalDetailTable()), w.ID).Scan(&data).Error
	return data, err
}

// IsAudioCustomer 判断是否客服审核
func (w *WithdrawalDetail) IsAudioCustomer(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("id = ?", w.ID).
		Find(w).Error
}

// UpdateCustomerStatus 更新客服状态
func (w *WithdrawalDetail) UpdateCustomerStatus(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("id = ? and customer_status = ?", w.ID, WithdrawalAudioStatusAwait).
		Update(map[string]interface{}{
			"updated_at":      time.Now(),
			"customer_status": w.CustomerStatus,
			"customer_id":     w.CustomerId,
		}).Error
}

// UpdateFinancialStatus 更新财务状态
func (w *WithdrawalDetail) UpdateFinancialStatus(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("id = ? and financial_status = ?", w.ID, WithdrawalAudioStatusAwait).
		Update(map[string]interface{}{
			"updated_at":       time.Now(),
			"financial_status": w.FinancialStatus,
			"financial_id":     w.FinancialId,
		}).Error
}

// UpdateRemark 更新备注
func (w *WithdrawalDetail) UpdateRemark(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("id = ? ", w.ID).
		Update(map[string]interface{}{
			"from_address":     w.FromAddress,
			"customer_status":  w.CustomerStatus,
			"financial_status": w.FinancialStatus,
			"status":           w.Status,
			"updated_at":       time.Now(),
			"remark":           w.Remark,
		}).Error
}

// UpdateOrderIdRemark 根据order_id更新备注
func (w *WithdrawalDetail) UpdateOrderIdRemark(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("order_id = ? ", w.OrderId).
		Update(map[string]interface{}{
			"callback_status": w.CallbackStatus,
			"callback_json":   w.CallbackJson,
			"status":          w.Status,
			"updated_at":      time.Now(),
			"remark":          w.Remark,
		}).Error
}

// UpdateStatus 更新状态
func (w *WithdrawalDetail) UpdateStatus(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("id = ? and financial_status = ? and customer_status = ?", w.ID, WithdrawalAudioStatusOk, WithdrawalAudioStatusOk).
		Update(map[string]interface{}{
			"callback_status":  w.CallbackStatus,
			"callback_json":    w.CallbackJson,
			"block_count":      w.BlockCount,
			"transaction_hash": w.TransactionHash,
			"updated_at":       time.Now(),
			"status":           w.Status,
			"remark":           w.Remark,
		}).Error
}

// GetOrderIdBySubmitInfo 根据订单id获取提交信息
func (w *WithdrawalDetail) GetOrderIdBySubmitInfo(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("order_id = ? and financial_status = ? and customer_status = ? and status = ?", w.OrderId, WithdrawalAudioStatusOk, WithdrawalAudioStatusOk, WithdrawalStatusSubmit).Find(w).Error
}

// GetOrderIdByInfo 根据订单id获取信息
func (w *WithdrawalDetail) GetOrderIdByInfo(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("order_id = ? ", w.OrderId).Find(w).Error
}

// WithdrawalStatusCancel 订单取消
func (w *WithdrawalDetail) WithdrawalStatusCancel(o *gorm.DB) error {
	return o.Table(GetWithdrawalDetailTable()).
		Where("id = ?  and status < ?", w.ID, WithdrawalStatusCancel).
		Update(map[string]interface{}{
			"customer_status":  WithdrawalStatusNoThrough,
			"financial_status": WithdrawalStatusNoThrough,
			"updated_at":       time.Now(),
			"status":           WithdrawalStatusCancel,
		}).Error
}
