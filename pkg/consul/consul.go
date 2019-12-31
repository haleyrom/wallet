package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"log"
)

// Service 服务
type Service struct {
	Id   string `json:"id"`
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

// newClient 初始化
func newClient() (*consulapi.Client, error) {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = viper.GetString("consul.config.address")
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
		return nil, err
	}
	return client, nil
}

//注册服务到consul
func ConsulRegister() error {
	client, err := newClient()
	if err == nil {
		// 创建注册到consul的服务到
		registration := new(consulapi.AgentServiceRegistration)
		registration.ID = viper.GetString("consul.register.id")
		registration.Name = viper.GetString("consul.register.name")
		registration.Port = viper.GetInt("consul.register.port")
		registration.Tags = viper.GetStringSlice("consul.register.tags")
		registration.Address = viper.GetString("consul.register.address")

		// 增加consul健康检查回调函数
		check := new(consulapi.AgentServiceCheck)
		check.HTTP = fmt.Sprintf("http://%s:%d", registration.Address, registration.Port)
		check.Timeout = "5s"
		check.Interval = "5s"
		check.DeregisterCriticalServiceAfter = "30s" // 故障检查失败30s后 consul自动将注册服务删除
		registration.Check = check

		// 注册服务到consul
		err = client.Agent().ServiceRegister(registration)
	}
	return err
}

// 取消consul注册的服务
func ConsulDeRegister() {
	// 创建连接consul服务配置
	client, err := newClient()
	if err == nil {
		_ = client.Agent().ServiceDeregister("111")
	}
}

// ConsulFindServer 从consul中发现服务
func ConsulFindServer() (map[string]Service, error) {
	// 创建连接consul服务配置
	client, err := newClient()
	if err == nil {
		data := make(map[string]Service, 0)
		// 获取所有service
		services, _ := client.Agent().Services()
		for _, value := range services {
			data[value.ID] = Service{
				Id:   value.ID,
				Addr: value.Address,
				Port: value.Port,
			}
		}
		return data, nil
	}
	return nil, err
}

// ConsulGetServer 获取服务
func ConsulGetServer(name string) (string, error) {
	// 创建连接consul服务配置
	client, err := newClient()
	if err == nil {
		// 获取所有service
		if services, _, err := client.Agent().Service(name, nil); err == nil {
			return fmt.Sprintf("http://%s:%d", services.Address, services.Port), nil
		}
		return "", err
	}
	return "", err
}
