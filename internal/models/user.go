package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"time"
)

// User 用户
type User struct {
	gorm.Model
	Uid         string `gorm:"size:50;index;column:uid;unique_index;comment:'用户id';"` // 用户id
	Email       string `gorm:"size:50;column:email;comment:'email'"`                  // 邮箱
	Name        string `gorm:"size:255;column:name;comment:'帐号名称'"`                   // 帐号名称
	PayPassword string `gorm:"size:255;column:pay_password;comment:'支付密码'"`           // 支付密码
}

// Table 表
func GetUserTable() string {
	return viper.GetString("mysql.prefix") + "user"
}

// NewAccount 初始化
func NewUser() *User {
	return &User{}
}

// IsExistUser 判断用户是否存在
func (u *User) IsExistUser(o *gorm.DB) error {
	if err := o.Table(GetUserTable()).Where("uid = ?", u.Uid).First(u).Error; err == gorm.ErrRecordNotFound {
		return fmt.Errorf("%s", "user not exist")
	}
	return nil
}

// CreateUser 创建用户
func (u *User) CreateUser(o *gorm.DB) error {
	return o.Create(u).Error
}

// GetAllByUid 获取用户id
func (u *User) GetAllByUid(o *gorm.DB) ([]uint, error) {
	data := make([]uint, 0)
	rows, err := o.Table(GetUserTable()).Select("id").Rows()
	defer rows.Close()

	if err == nil {
		var id uint
		for rows.Next() {
			if err = rows.Scan(&id); err == nil {
				data = append(data, id)
			}
		}
	}
	return data, err
}

// UpdatePayPassword  更新支付密码
func (u *User) UpdatePayPassword(o *gorm.DB) error {
	if err := o.Model(u).Where("id = ? ", u.ID).Update(map[string]interface{}{
		"pay_password": u.PayPassword,
		"updated_at":   time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}

// GetInfo 获取信息
func (u *User) GetInfo(o *gorm.DB) error {
	return o.Table(GetUserTable()).Where("id = ?", u.ID).Find(u).Error
}

// IsSetPayPassword 判断是否设置支付密码
func (u *User) IsSetPayPassword(o *gorm.DB) bool {
	if o.Table(GetUserTable()).First(u); len(u.PayPassword) > 0 {
		return true
	}
	return false
}

// UpdateInfo  更新信息
func (u *User) UpdateInfo(o *gorm.DB) error {
	if err := o.Model(u).Where("uid = ? ", u.Uid).Update(map[string]interface{}{
		"email":      u.Email,
		"name":       u.Name,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	return nil
}
