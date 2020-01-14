package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"math"
	"time"
)

// Order 订单
type Order struct {
	gorm.Model
	Uid         uint    `gorm:"column:uid;default:0;comment:'用户id'"`                    // 用户id
	ExchangeUid uint    `gorm:"column:exchange_uid;default:0;comment:'转入用户id'"`         // 转入用户id
	Context     string  `gorm:"type:text;column:context;comment:'文本'"`                  // 文本
	CurrencyId  uint    `gorm:"column:currency_id;default:0;comment:'币种id'"`            // 币种id
	ExchangeId  uint    `gorm:"column:exchange_id;default:0;comment:'兑换id'"`            // 兑换币种id
	Balance     float64 `gorm:"column:balance;default:0;comment:'余额';"`                 // 余额
	Ratio       float64 `gorm:"column:ratio;default:0;comment:'比例'"`                    // 比率
	Form        string  `gorm:"column:form;comment:'来源'"`                               // 来源
	Status      int8    `gorm:"size(3);column:status;default:0;comment:'状态(0未成功,1成功)'"` // 订单状态
	Type        int8    `gorm:"size(3);column:type;default:0;comment:'订单类型(0兑换本地的币)'"`  // 订单类型
	OrderUuid   string  `gorm:"size(200);column:order_uuid;comment:'订单uuid'"`           //订单号
	Symbol      string  `gorm:"size:255;column:symbol;comment:'代币代号';"`                 // 代币代号
}

var (
	// OrderStatusNot 未完成
	OrderStatusNot int8 = 0
	// OrderStatusOk 完成
	OrderStatusOk int8 = 1
	// OrderTypeChange 兑换类型 兑换
	OrderTypeChange int8 = 0
	// OrderTypeShare 分红
	OrderTypeShare int8 = 1
	// OrderTypeTransfer 转账
	OrderTypeTransfer int8 = 2
	// OrderTypePayment 支付
	OrderTypePayment int8 = 3
	//  OrderFormUsdd 算力
	OrderFormUsdd string = "usdd"
	// OrderFormTransfer 个人转账
	OrderFormTransfer string = "transfer"
	// OrderFormPayment 支付
	OrderFormPayment string = "payment"
)

// GetOrderTable 订单地址
func GetOrderTable() string {
	return viper.GetString("mysql.prefix") + "order"
}

// NewOrder 初始化
func NewOrder() *Order {
	return &Order{}
}

// CreateOrder 创建订单
func (r *Order) CreateOrder(o *gorm.DB) error {
	r.OrderUuid = fmt.Sprintf("%s", uuid.NewV4())
	return o.Create(r).Error
}

// IsOrderUuid 判断订单是否存在
func (r *Order) IsOrderUuid(o *gorm.DB) error {
	return o.Table(GetOrderTable()).
		Where("order_uuid = ?", r.OrderUuid).
		Find(r).
		Error
}

// RemoveOrder 删除账单
func (r *Order) RemoveOrder(o *gorm.DB) error {
	timer := time.Now()
	if err := o.Table(GetOrderTable()).Where("id = ? ", r.ID).Update(map[string]interface{}{
		"updated_at": timer,
		"deleted_at": timer,
	}).Error; err != nil {
		return err
	}
	return nil
}

// RemoveOrder 删除账单
func (r *Order) RemoveOrderUuid(o *gorm.DB) error {
	return o.Where("order_uuid = ?", r.OrderUuid).
		Delete(r).
		Error
}

// GetAllTransOrder GetAllTransOrder
func (r *Order) GetAllTransOrder(o *gorm.DB, page, pageSize int, endTime, startTime int, key string) (resp.UserTransOrderResp, error) {
	data := resp.UserTransOrderResp{}
	if endTime == 0 {
		endTime = 10000000000000
	}

	count_sql := fmt.Sprintf("select count(*) as num from %s o left join %s u on o.uid = u.id  where o.type = %d and UNIX_TIMESTAMP(o.updated_at) >= %d and UNIX_TIMESTAMP(o.updated_at) <= %d ", GetOrderTable(), GetUserTable(), OrderTypeChange, startTime, endTime)
	sql := fmt.Sprintf("SELECT o.id,o.uid,o.currency_id,o.exchange_id,o.updated_at,TRUNCATE(o.balance,6) as value,o.status,o.ratio,u.name,u.email FROM %s o LEFT JOIN %s u on u.id = o.uid  where o.type = %d and UNIX_TIMESTAMP(o.updated_at) >= %d and UNIX_TIMESTAMP(o.updated_at) <= %d  ", GetOrderTable(), GetUserTable(), OrderTypeChange, startTime, endTime)
	if key != "" {
		count_sql += " and u.name like  '%" + key + "%' "
		sql += " and u.name like  '%" + key + "%' "
	}
	sql = sql + fmt.Sprintf("order by o.id desc limit %d offset %d", pageSize, (page-1)*pageSize)
	rows, err := o.Raw(sql).Rows()
	defer rows.Close()
	var (
		item  resp.UserTransInfoOrderResp
		timer time.Time
	)
	data.Items = make([]resp.UserTransInfoOrderResp, 0)

	if err == nil {
		currencty := make(map[int]string, 0)
		list, _ := NewCurrency().GetAll(o)
		// Fixme: 优化
		for _, val := range list {
			currencty[int(val.CurrencyId)] = val.Symbol
		}

		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				item.UpdatedAt = tools.TimerConvert(timer, item.UpdatedAt)
				if _, ok := currencty[item.CurrencyId]; ok == true {
					item.CurrencySymbol = currencty[item.CurrencyId]
				}
				if _, ok := currencty[item.ExchangeId]; ok == true {
					item.ExchangeSymbol = currencty[item.ExchangeId]
				}
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

// UpdateStatusOk 更新订单状态
func (r *Order) UpdateStatusOk(o *gorm.DB) error {
	timer := time.Now()
	if err := o.Table(GetOrderTable()).Where("id = ? ", r.ID).Update(map[string]interface{}{
		"updated_at":   timer,
		"balance":      r.Balance,
		"context":      r.Context,
		"exchange_uid": r.ExchangeUid,
		"status":       OrderStatusOk,
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetAccountTransferList 获取转账列表
func (r *Order) GetAccountTransferList(o *gorm.DB, page, pageSize int, endTime, startTime int, key string) (resp.AccountTransferListResp, error) {
	data := resp.AccountTransferListResp{}
	if endTime == 0 {
		endTime = 10000000000000
	}

	count_sql := fmt.Sprintf("select count(*) as num from %s o left join %s u on o.uid = u.id LEFT JOIN %s us on us.id = o.exchange_uid where o.type = %d and UNIX_TIMESTAMP(o.updated_at) >= %d and UNIX_TIMESTAMP(o.updated_at) <= %d ", GetOrderTable(), GetUserTable(), GetUserTable(), OrderTypeTransfer, startTime, endTime)
	sql := fmt.Sprintf("SELECT o.order_uuid as order_id,u.uid as uid,u.name as user_name,u.email as user_email,us.uid as adverse_id,us.name as adverse_name,us.email as adverse_email,TRUNCATE(o.balance,6) as balance,o.status,o.updated_at,o.symbol FROM %s o LEFT JOIN %s u on u.id = o.uid LEFT JOIN %s us on us.id = o.exchange_uid  where o.type = %d and UNIX_TIMESTAMP(o.updated_at) >= %d and UNIX_TIMESTAMP(o.updated_at) <= %d  ", GetOrderTable(), GetUserTable(), GetUserTable(), OrderTypeTransfer, startTime, endTime)
	if key != "" {
		count_sql += " and ((u.name like  '%" + key + "%') or (us.name like  '%" + key + "%') ) "
		sql += " and ((u.name like  '%" + key + "%' ) or (us.name like  '%" + key + "%' ) ) "
	}
	sql = sql + fmt.Sprintf("order by o.id desc limit %d offset %d", pageSize, (page-1)*pageSize)
	rows, err := o.Raw(sql).Rows()
	defer rows.Close()
	var (
		item  resp.AccountTransferInfoResp
		timer time.Time
	)
	data.Items = make([]resp.AccountTransferInfoResp, 0)

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
