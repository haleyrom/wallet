package api

import (
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
)

// ReadListBlockChain 读取链列表
// @Tags  BlockChain 链
// @Summary 读取链列表接口
// @Description 读取链列表
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} resp.ReadBlockChainListResp
// @Router /chain/list [get]
func ReadListBlockChain(c *gin.Context) {
	data, err := models.NewBlockChain().GetAll(core.Orm.New())
	if err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// ReadListSymbolBlockChain 根据symbol获取链
// @Tags  BlockChain 链
// @Summary 根据symbol获取链接口
// @Description 根据symbol获取链
// @Produce json
// @Security ApiKeyAuth
// @Param symbol formData int true "链标识"
// @Success 200 {object} resp.ReadOrderSymbolByChainResp
// @Router /chain/symbol [get]
func ReadListSymbolBlockChain(c *gin.Context) {
	p := &params.ReadListSymbolBlockChainParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	coin := models.NewCoin()
	coin.Symbol = p.Symbol
	data, err := coin.GetOrderSymbolByChain(core.Orm.New())
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, data)
	return
}

// UpdateCoin 更新链
// @Tags  BlockChain 链
// @Summary 更新链接口
// @Description 更新链
// @Produce json
// @Security ApiKeyAuth
// @Param chain_id formData int true "链id"
// @Param chain_code formData string true "链标识"
// @Param name formData string true "名称"
// @Success 200
// @Router /chain/update [POST]
func UpdateBlockChain(c *gin.Context) {
	p := &params.UpdateBlockChainParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	chain := models.NewBlockChain()
	chain.ID = p.BlockChainId
	o := core.Orm.DB.New()
	if err := chain.IsExistBlockChain(o); err != nil {
		core.GResp.Failure(c, resp.CodeNotChain)
		return
	}

	chain.ChainCode = p.ChainCode
	chain.Name = p.Name
	if err := chain.UpdateBlockChain(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	core.GResp.Success(c, resp.EmptyData())
	return
}

// AddCoin 增加链
// @Tags  BlockChain 链
// @Summary 增加链接口
// @Description 增加链
// @Produce json
// @Security ApiKeyAuth
// @Param chain_code formData int true "链标识"
// @Param name formData string true "名称"
// @Success 200
// @Router /chain/add [post]
func AddBlockChain(c *gin.Context) {
	p := &params.AddBlockChainParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	chain := models.BlockChain{
		ChainCode: p.ChainCode,
		Name:      p.Name,
	}
	if err := chain.CreateBlockChain(core.Orm.New()); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// RemoveCoin 删除链
// @Tags  BlockChain 链
// @Summary 删除链接口
// @Description 删除链
// @Produce json
// @Security ApiKeyAuth
// @Param chain_id query int true "链id"
// @Success 200
// @Router /chain/remove [POST]
func RemoveBlockChain(c *gin.Context) {
	p := &params.RemoveBlockChainParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}

	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	}

	chain := models.NewBlockChain()
	chain.ID = p.BlockChainId
	o := core.Orm.DB.New()
	if err := chain.IsExistBlockChain(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}

	if err := chain.RmChain(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}
