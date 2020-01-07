package core

import (
	"fmt"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/core/cron"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/pkg/jwt"
	"github.com/haleyrom/wallet/pkg/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// cfg  配置
	configFilePath = pflag.StringP("config", "c", "assets/config/conf.yaml", "apiserver config file path.")
)

// Init init
func Init() {
	InitConf()
	InitJwt()
	InitStorage()
	cron.InitCron()
}

// InitConf  初始化配置
func InitConf() {
	pflag.Parse()
	// 获取配置
	if err := core.Conf.Init(*configFilePath); err != nil {
		logrus.Error("err parsing  config file:", err)
		panic(fmt.Errorf("err parsing  config file:", err))
	}
}

// InitStorage 初始化数据库
func InitStorage() {
	if err := core.Orm.Init(viper.GetString("mysql.addr"), viper.GetString("mysql.prefix")); err != nil {
		logrus.Error("init orm client fail:", err)
		panic(fmt.Errorf("init orm client fail:", err))
	}
	install()
}

// install 注册
func install() {
	table := []interface{}{
		models.NewUser(),
		models.NewAccount(),
		models.NewAccountDetail(),
		models.NewBlockDetail(),
		models.NewCoin(),
		models.NewOrder(),
		models.NewCurrency(),
		models.NewBlockChain(),
		models.NewDepositAddr(),
		models.NewDepositDetail(),
		models.NewWithdrawalAddr(),
		models.NewWithdrawalDetail(),
		models.NewEmailCode(),
		models.NewCompanyStream(),
		models.NewCompanyAddr(),
		models.NewQuote(),
		models.NewQuoteHistory(),
	}

	core.Orm.Set("gorm:table_options", "ROW_FORMAT=DYNAMIC").AutoMigrate(table...)
}

// InitJwt 初始化jwt
func InitJwt() {
	middleware.Jwt = jwt.NewJWT(viper.GetString("jwt.signkey"))
	middleware.AdminJwt = jwt.NewJWT(viper.GetString("jwt.admin_signkey"))
}
