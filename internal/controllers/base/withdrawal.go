package base

import (
	"fmt"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/http"
)

// WithdrawalAudioOK 提现处理
func WithdrawalAudioOK(o *gorm.DB, detail *models.WithdrawalDetail) (string, string, error) {
	// TODO: 等待调试提币接口
	consul_service, err := consul.ConsulGetServer("blockchain-pay.tfor")
	if err != nil {
		return core.DefaultNilString, core.DefaultNilString, err
	}

	coin := models.NewCoin()
	coin.Symbol = detail.Symbol
	contract_address, err := coin.GetConTractAddress(o)
	if err != nil {
		return core.DefaultNilString, core.DefaultNilString, err
	}

	company_addr := models.NewCompanyAddr()
	company_addr.Symbol, company_addr.Code = detail.Symbol, models.CodeWithdrawal
	address, err := company_addr.GetOrderSymbolByAddress(o)
	if err != nil {
		return core.DefaultNilString, fmt.Sprintf("%s", resp.CodeNotCompanyAddress), resp.CodeNotCompanyAddress
	}

	url := fmt.Sprintf("%s%s", consul_service, "/api/v1/blockchain-pay/ethtereum/withdrawal")

	data := map[string]interface{}{
		"app_id":           viper.GetString("appname"),
		"order_id":         detail.OrderId,
		"symbol":           detail.Symbol,
		"contract_address": contract_address,
		"from_address":     address,
		"to_address":       detail.Address,
		"value":            fmt.Sprintf("%.2f", detail.Value),
	}

	result, err := tools.WithdrawalAudio(data, url, viper.GetString("deposit.Srekey"))

	if result == nil || err != nil || result.Code != http.StatusOK {
		return core.DefaultNilString, core.DefaultNilString, errors.Errorf("%s", result.Msg)
	}

	return address, result.Msg, nil
}
