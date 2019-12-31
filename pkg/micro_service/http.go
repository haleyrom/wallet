package micro_service

import (
	"fmt"
	"gopkg.in/resty.v1"
)

var (
	// defaultNilString defaultNilString
	defaultNilString string = ""
)

// MicroServiceHttp 微服务请求
type MicroServiceHttp struct {
	Token string `json:"token"`
	Url   string `json:"url"`
}

// NewMicroServiceHttp 初始化微服务请求
func NewMicroServiceHttp() *MicroServiceHttp {
	return &MicroServiceHttp{}
}

// HttpPost post请求
func (m *MicroServiceHttp) HttpPost(body map[string]string) (string, error) {
	client := resty.New()
	req := client.R()
	req.SetFormData(body)
	req.SetHeader(`content-type`, "application/json")
	req.SetHeader(`Authorization`, m.Token)
	resp, err := req.Post(m.Url)
	client.SetCloseConnection(true)
	if err != nil {
		return defaultNilString, err
	}
	fmt.Println(resp)
	return string("result"), nil
}
