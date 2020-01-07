package cron

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"time"
)

const timer_second = 3

// gcPayMap
type gcPayMap struct{}

func (g *gcPayMap) Id() string {
	return JobIdGcGcPayMap
}

// Spec 每30分钟gc一次
func (g *gcPayMap) Spec() string {
	// Fixme: 优化从配置文件读取值
	return fmt.Sprintf("@every %ds", timer_second)
}

// Job 按时间LRU策略回收支付map
func (g *gcPayMap) Job() {
	logs.Debug("gc pay map")
	// TODO

	timer := time.Now().Unix()
	order := models.NewOrder()
	for i := 0; i < timer_second*2; i++ {
		if _, ok := core.PayChan.MapTime[int(timer)-i]; ok == true {
			for _, v := range core.PayChan.MapTime[int(timer)-i] {
				delete(core.PayChan.MapChan, v)
				order.OrderUuid = v
				_ = order.RemoveOrderUuid(core.Orm.New())
			}
			delete(core.PayChan.MapTime, int(timer)-i)
		}
	}
}
