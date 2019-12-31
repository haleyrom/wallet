package storage

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient mongo数据库
type MongoClient struct {
	Client *mongo.Client
}

// Init Init
func (m *MongoClient) Init(addr string) error {
	var err error
	opt := options.Client().ApplyURI(m.parseMongodb(addr)) // mongo client
	opt.SetMaxPoolSize(200)                                // 使用最大的连接数
	//opt.SetLocalThreshold(3 * time.Second)                 // 只使用与mongo操作耗时小于3秒的
	//opt.SetMaxConnIdleTime(5 * time.Second)                // 指定连接可以保持空闲的最大毫秒数
	if m.Client, err = mongo.Connect(context.TODO(), opt); err != nil {
		logrus.Errorf("mongo client link failure : %s", err)
		return err
	}
	return nil
}

// ParseMongodb 解析mongodb配置
func (m *MongoClient) parseMongodb(addr string) string {
	return fmt.Sprintf("mongodb://%s", addr)
}
