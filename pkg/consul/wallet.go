package consul

import (
	"fmt"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
)

// GetWalletAddress 获取钱包地址
func GetWalletAddress(chain string) (*tools.WalletResp, error) {
	service_url, err := ConsulGetServer("blockchain-pay.tfor")
	if len(service_url) == 0 {
		logrus.Errorf("request consul failure: %s,  %v", service_url, err)
		return nil, errors.New("deposit addr url is null")
	}

	if val, ok := viper.GetStringMap("deposit.addr")[chain]; ok && len(val.(string)) > 0 {
		url := fmt.Sprintf("%s%s", service_url, val.(string))
		if data, err := tools.RegisterWalletAddr(viper.GetString("appname"), url, viper.GetString("deposit.srekey")); err == nil {
			return data, nil
		} else {
			fmt.Println("err:", err)
		}
	}
	return nil, errors.New("deposit addr is null")
}

// IsWalletAddress 判断钱包地址
func IsWalletAddress(address string) error {
	service_url, err := ConsulGetServer("blockchain-pay.tfor")
	if len(service_url) == 0 {
		logrus.Errorf("request consul failure: %s,  %v", service_url, err)
		return err
	}

	data := map[string]interface{}{
		"address": address,
	}

	url := fmt.Sprintf("%s%s", service_url, "/api/v1/blockchain-pay/ethtereum/is-address")
	if data, err := tools.HttpPost(data, url, core.DefaultNilString); err == nil {
		if data.Code == http.StatusOK {
			return nil
		}
		return errors.Errorf("%s", data.Msg)
	} else {
		return err
	}
}

// GetUserInfo 获取用户信息
func GetUserInfo(uid, token string) (interface{}, error) {
	service_url, err := ConsulGetServer("user.tfor")
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"user_id": uid,
	}

	head := map[string]string{
		"Authorization": token,
	}

	url := fmt.Sprintf("%s%s", service_url, "/api/v1/user/inside/get-user")
	if data, err := tools.HttpGetBase(url, data, head); err == nil {
		if data.Code == http.StatusOK {
			return data.Data, nil
		}
		return nil, errors.Errorf("%s", data.Msg)
	} else {
		return nil, err
	}
}

// GetOrderEmailByInfo 根据email获取信息
func GetOrderEmailByInfo(email, token string) (interface{}, error) {
	service_url, err := ConsulGetServer("user.tfor")
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"email": email,
	}

	head := map[string]string{
		"Authorization": token,
	}

	url := fmt.Sprintf("%s%s", service_url, "/api/v1/user/inside/get-user-email")
	if data, err := tools.HttpGetBase(url, data, head); err == nil {
		if data.Code == http.StatusOK {
			return data.Data, nil
		}
		return nil, errors.Errorf("%s", data.Msg)
	} else {
		return nil, err
	}
}
