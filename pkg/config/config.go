package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

// LogsConf 日志配置
type LogsConf struct {
	Path   string `yaml:"path"`
	Name   string `yaml:"name"`
	Suffix string `yaml:"suffix"`
}

// MongoConf mongo配置
type MongoConf struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// MysqlConf 数据库配置
type MysqlConf struct {
	Addr   string `yaml:"addr"`
	Prefix string `yaml:"prefix"`
}

// DepositConf 充值配置
type DepositConf struct {
	Addr     map[string]interface{} `yaml:"addr"`
	Ethereum string                 `yaml:"ethereum"`
	Srekey   string                 `yaml:"srekey"`
}

// LoadConfig 加载配置
type Configure struct{}

// Init 初始化
func (c *Configure) Init(path string) error {
	// 初始化配置文件
	if err := c.initConfig(path); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	c.watchConfig()
	return nil
}

// initConfig 初始化配置
func (c *Configure) initConfig(path string) error {
	if len(path) > 0 {
		viper.SetConfigFile(path) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("assets/config") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("conf")
	}
	viper.SetConfigType("yaml")  // 设置配置文件格式为YAML
	viper.AutomaticEnv()         // 读取匹配的环境变量
	viper.SetEnvPrefix("wallet") // 读取环境变量的前缀为APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

// 监控配置文件变化并热加载程序
func (c *Configure) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}
