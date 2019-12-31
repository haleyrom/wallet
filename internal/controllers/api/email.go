package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/models"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/consul"
	"github.com/haleyrom/wallet/pkg/tools"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

//生成盐的盐
const salt = 1

// @Tags 支付密码
// @Summary 判断是否设置支付密码接口
// @Description 使用jwt命令牌
// @Produce json
// @Success 200
// @Router /user/pay/is_init [get]
func IsSetPassWord(c *gin.Context) {
	p := &params.BaseBindParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	o := core.Orm.DB.New()
	user := models.NewUser()
	user.ID = p.Base.Uid
	// 判断支付密码是否设置
	if ok := user.IsSetPayPassword(o); ok == true {
		core.GResp.Failure(c, resp.CodeSetPayPassword)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// @Tags 支付密码
// @Summary 设置支付密码接口
// @Description 使用jwt命令牌
// @Produce json
// @Param password query string true "用户支付密码"
// @Success 200
// @Router /user/pay/set-paypwd [post]
func SetPayPassWordHandler(c *gin.Context) {
	p := &params.SetPayPassWordHandlerParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
		return
	} else if len(p.Password) != 6 {
		core.GResp.Failure(c, resp.CodeIllegalPassword, err)
		return
	}

	o := core.Orm.New()
	user := models.NewUser()
	user.ID = p.Base.Uid
	// 判断支付密码
	if ok := user.IsSetPayPassword(o); ok == true {
		core.GResp.Failure(c, resp.CodeSetPayPassword)
		return
	}

	user.Uid = p.Base.Claims.UserID
	user.PayPassword = tools.Hash256(p.Password, tools.NewPwdSalt(p.Base.Claims.UserID, salt))
	if err := user.UpdatePayPassword(o); err != nil {
		core.GResp.Failure(c, err)
		return
	}
	core.GResp.Success(c, resp.EmptyData())
	return
}

// @Tags 支付密码
// @Summary 重置支付密码
// @Description 使用jwt命令牌
// @Produce json
// @Param email    query string true "用户邮箱"
// @Param password query string true "新密码"
// @Param code query string true "验证码"
// @Success 200
// @Router /user/pay/reset-paypwd [post]
func ReSetPayWordHandler(c *gin.Context) {
	var err error
	p := &params.ReSetPayWordHandlerParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err = c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
	}
	//查找验证码是否正确
	emailCode := models.NewEmailCode()
	emailCode.Uid = p.Base.Claims.UserID
	emailCode.Email = p.Email
	o := core.Orm.DB.New()
	err, ecode := emailCode.FindCode(o)
	if err != nil {
		core.GResp.Failure(c, err)
		return
	}
	if ecode.Code == p.Code {
		//更新密码
		pay := models.User{
			Model: gorm.Model{
				ID: p.Base.Uid,
			},
			PayPassword: tools.Hash256(p.Password, tools.NewPwdSalt(p.Base.Claims.UserID, salt)),
		}

		if err = pay.UpdatePayPassword(o); err != nil {
			core.GResp.Failure(c, err)
			return
		}

		go func() {
			_ = emailCode.Del(o)
		}()

		core.GResp.Success(c, resp.EmptyData())
	} else {
		core.GResp.Failure(c, resp.CodeCodeError)
	}
	return
}

// @Tags 支付密码
// @Summary 修改支付密码时发送邮箱验证码接口
// @Description 使用jwt命令牌
// @Param email  query string true "用户邮箱"
// @Produce json
// @Success 200
// @Router /user/pay/send-email_util [post]
func SendEmailPayHandler(c *gin.Context) {
	email := c.Request.FormValue("email")
	if !tools.VerifyEmailFormat(email) {
		core.GResp.Failure(c, errors.New("邮件格式错误"))
		return
	}

	p := &params.UserIdParam{
		Base: core.UserInfoPool.Get().(*params.BaseParam),
	}
	// 绑定参数
	if err := c.ShouldBind(p); err != nil {
		core.GResp.Failure(c, resp.CodeIllegalParam, err)
	}

	code := tools.RandStr()
	duration, _ := time.ParseDuration("5m")
	sms := models.EmailCode{
		Uid:     p.Base.Claims.UserID,
		Email:   email,
		Code:    code,
		OutTime: time.Now().Add(duration).Unix(),
	}
	o := core.Orm.New().Begin()
	if err := sms.Del(o); err != nil && err != gorm.ErrRecordNotFound {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	} else if err := sms.Insert(o); err != nil {
		o.Rollback()
		core.GResp.Failure(c, err)
		return
	}
	// 服务发现，发送邮件
	if service, err := consul.ConsulGetServer("auth.tfor"); err == nil {
		if err, responseData := tools.SendHttpPost(service, "/api/v1/auth/send-emailcode", map[string]string{
			"code":  code,
			"email": email,
		}, c.Request.Header.Get(core.HttpHeadToken)); err != nil {
			fmt.Println(err, "错误")
			logrus.Errorf("SendEmailPayHandler send failure : %s data : %v", err.Error(), responseData)
		}
	} else {
		fmt.Println(err, "错误")
		logrus.Errorf("consul failure : %s", err.Error())
	}
	o.Commit()
	core.GResp.Success(c, resp.EmptyData())
	return
}
