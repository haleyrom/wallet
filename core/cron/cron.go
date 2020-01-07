package cron

import (
	"github.com/haleyrom/wallet/pkg/cron"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	// gcCacheOne 控制初始化gc cache
	gcCacheOne = new(sync.Once)

	// m m
	m *gcPayMap

	// Signals 退出的信号
	Signals chan os.Signal
)

// Job's Id
const (
	JobIdGcGcPayMap = "JobIdGcPayMap"
)

func init() {
	// 开启缓存gc
	gcCacheOne.Do(func() {
		cron.StartCron()
	})
	// 退出信号
	//Signals = make(chan os.Signal, 1)
	//signal.Notify(Signals, os.Interrupt, os.Kill)

	//go func() {
	//	CaptureSignal()
	//}()

}

// Init Init
func InitCron() {
	_ = cron.AddSpecJob(m, true)
}

// CaptureSignal CaptureSignal
func CaptureSignal() {
	select {
	case s := <-Signals:
		// TODO: 退出处理
		logrus.Debug(s)
	default:
	}

}
