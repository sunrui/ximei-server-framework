/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-08-30 23:02:38
 */

package user

import (
	"framework/pkg/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Model 用户模型
type Model struct {
	Name     string `json:"name" gorm:"type:varchar(32);index;comment:用户名" validate:"ascii|min_len:2|max_len:32"` // 用户名
	IdCard   string `json:"idCard" gorm:"type:char(18);index;comment:身份证" validate:"len:18"`                      // 身份证
	Phone    string `json:"phone" gorm:"type:varchar(11);index;comment:手机号" validate:"numeric|len:11"`            // 手机号
	Password string `json:"password" gorm:"binary(60);comment:密码" validate:"len:60"`                              // 密码
}

// User 用户
type User struct {
	mysql.ModelId
	mysql.ModelTenantId

	Model

	mysql.ModelEnable
	mysql.ModelTime
}

// BeforeSave 创建前
func (user *User) BeforeSave(tx *gorm.DB) error {
	if pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0); err != nil {
		tx.Statement.SetColumn("password", pw)
		return nil
	} else {
		return err
	}
}

// IsValidPassword 验证密码
func (user *User) IsValidPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
