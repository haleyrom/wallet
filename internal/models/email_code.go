package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sort"
	"time"
)

// Account 帐号
type EmailCode struct {
	Uid     string `gorm:"size:50;index;column:uid;comment:'用户id';"`  // 用户id
	Code    string `gorm:"column:code;comment:'验证码';"`                // 验证码
	OutTime int64  `gorm:"column:out_time;default:0;comment:'过期时间';"` // 过期时间
	Email   string `gorm:"column:email"`
}

// Table 表
func GetEmailCodeTable() string {
	return viper.GetString("mysql.prefix") + "email_code"
}

// NewAccount 初始化
func NewEmailCode() *EmailCode {
	return &EmailCode{}
}

func (e *EmailCode) Insert(o *gorm.DB) error {
	err := o.Model(e).Create(e).Error
	return err
}

func (e *EmailCode) FindCode(o *gorm.DB) (error, EmailCode) {
	var codes codes
	err := o.Table(GetEmailCodeTable()).Where("email = ?", e.Email).Find(&codes).Error
	if err != nil {
		return err, EmailCode{}
	}
	if len(codes) == 0 {
		return errors.New("no code"), EmailCode{}
	}
	if len(codes) > 2 {
		sort.Stable(codes)
	}
	if codes[0].OutTime < time.Now().Unix() {
		return errors.New("code outtime"), EmailCode{}
	}
	return nil, codes[0]
}
func (e *EmailCode) Del(o *gorm.DB) error {
	err := o.Table(GetEmailCodeTable()).Where("email = ?", e.Email).Delete(e).Error
	return err
}

type codes []EmailCode

func (s codes) Len() int { return len(s) }

func (s codes) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s codes) Less(i, j int) bool { return s[i].OutTime < s[j].OutTime }
