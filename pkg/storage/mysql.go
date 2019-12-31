package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// MysqlClient Mysql数据库
type MysqlClient struct {
	*gorm.DB
}

// Init Init
func (m *MysqlClient) Init(addr, prefix string) error {
	var err error

	//设置默认表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return prefix + defaultTableName
	}

	//连接数据库
	if m.DB, err = gorm.Open("mysql", addr); err != nil {
		logrus.Errorf("mysql client link failure : %s", err)
		return err
	}

	m.SingularTable(true)
	m.DB.DB().SetMaxIdleConns(100)
	m.DB.DB().SetMaxOpenConns(200)
	m.DB.LogMode(true)
	return nil
}
