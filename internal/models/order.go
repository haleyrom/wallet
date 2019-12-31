package models

import (
	"fmt"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"math"
	"time"
)

// Order 订单
type Order struct {
	gorm.Model
	Uid        uint    `gorm:"column:uid;default:0;comment:'用户id'"`                    // 用户id
	Context    string  `gorm:"type(context);column:context;comment:'文本'"`              // 文本
	CurrencyId uint    `gorm:"column:currency_id;default:0;comment:'币种id'"`            // 币种id
	ExchangeId uint    `gorm:"column:exchange_id;default:0;comment:'兑换id'"`            // 兑换币种id
	Balance    float64 `gorm:"column:balance;default:0;comment:'余额';"`                 // 余额
	Ratio      float64 `gorm:"column:ratio;default:0;comment:'比例'"`                    // 比率
	Form       string  `gorm:"column:form;comment:'来源'"`                               // 来源
	Status     int8    `gorm:"size(3);column:status;default:0;comment:'状态(0未成功,1成功)'"` // 订单状态
	Type       int8    `gorm:"size(3);column:type;default:0;comment:'订单类型(0兑换本地的币)'"`  // 订单类型
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
	//  OrderFormUsdd 算力
	OrderFormUsdd string = "usdd"
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
	return o.Model(r).Create(r).Error
}

// GetAllTransOrder GetAllTransOrder
func (r *Order) GetAllTransOrder(o *gorm.DB, page, pageSize int, endTime, startTime int, key string) (resp.RespUserTransOrder, error) {
	data := resp.RespUserTransOrder{}
	if endTime == 0 {
		endTime = 10000000000000
	}

	count_sql := fmt.Sprintf("select count(*) as num from %s o left join %s u on o.uid = u.id  where o.type = %d and UNIX_TIMESTAMP(o.updated_at) >= %d and UNIX_TIMESTAMP(o.updated_at) <= %d  ", GetOrderTable(), GetUserTable(), OrderTypeChange, startTime, endTime)
	sql := fmt.Sprintf("SELECT o.id,o.uid,o.currency_id,o.exchange_id,o.updated_at,o.balance as value,o.status,o.ratio,u.name,u.email FROM %s o LEFT JOIN %s u on u.id = o.uid  where o.type = %d and UNIX_TIMESTAMP(o.updated_at) >= %d and UNIX_TIMESTAMP(o.updated_at) <= %d  ", GetOrderTable(), GetUserTable(), OrderTypeChange, startTime, endTime)
	if key != "" {
		count_sql += " and o.name like  '%" + key + "%' "
		sql += " and o.name like  '%" + key + "%' "
	}
	sql = sql + fmt.Sprintf("order by o.id desc limit %d offset %d", pageSize, (page-1)*pageSize)
	rows, err := o.Raw(sql).Rows()
	defer rows.Close()
	var (
		item  resp.RespUserTransInfoOrder
		timer time.Time
	)
	data.Items = make([]resp.RespUserTransInfoOrder, 0)

	if err == nil {
		currencty := make(map[int]string, 0)
		list, _ := NewCurrency().GetAll(o)
		// Fixme: 优化
		for _, val := range list {
			currencty[int(val.CurrencyId)] = val.Symbol
		}

		for rows.Next() {
			if err = o.ScanRows(rows, &item); err == nil {
				timer, _ = time.Parse("2006-01-02T15:04:05+08:00", item.UpdatedAt)
				item.UpdatedAt = timer.Format("2006-01-02 15:04:05")
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
