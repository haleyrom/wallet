package main

import (
	cmd "github.com/haleyrom/wallet/cmd/core"
)

// main main
// @title 测试
// @version 0.0.1
// @description  测试
// @BasePath /api/v1/wallet
// @in header
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// 仅用于初始化swag
	cmd.Init()
}
