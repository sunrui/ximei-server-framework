/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-11-06 18:10:03
 */

package mysql

import (
	"framework/pkg/ip2region"
	"framework/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// ModelId id
type ModelId struct {
	Id string `json:"id" gorm:"primaryKey; not null; type:char(16); comment:主键 id"` // 主键 id
}

// BeforeCreate 创建前回调
func (modelId *ModelId) BeforeCreate(*gorm.DB) (err error) {
	modelId.Id = utils.NanoId(16)
	return nil
}

// ModelTenantId 租户 id
type ModelTenantId struct {
	TenantId string `json:"tenantId" gorm:"not null; type:char(32); comment:租户 id"` // 租户 id
}

// ModelUserId 用户 id
type ModelUserId struct {
	UserId string `json:"userId" gorm:"not null; type:char(32); comment:用户 id"` // 用户 id
}

// ModelRefer 来源
type ModelRefer struct {
	DeviceType DeviceType `json:"deviceType"  gorm:"type:varchar(32);not null;comment:设备类型" validate:"enum:ANDROID,IOS,WEB,H5,APPLET"` // 设备类型
	PhoneModel string     `json:"phoneModel" gorm:"type:varchar(32);comment:手机号型" validate:"max_len:32"`                               // 手机号型
	OsVersion  string     `json:"osVersion" gorm:"type:varchar(32);comment:系统版本" validate:"max_len:32"`                                // 系统版本
	AppPackage string     `json:"appPackage"  gorm:"type:varchar(32);comment:软件包名" validate:"max_len:32" `                             // 软件包名
	AppVersion string     `json:"appVersion" gorm:"type:varchar(32);comment:软件版本" validate:"max_len:32"`                               // 软件版本
}

// ModelHttp hypersonic
type ModelHttp struct {
	Ip                string  `json:"ip" gorm:"type:char(15);not null;comment:ip" validate:"required|max_len:15"` // ip
	ip2region.Address         // 位置
	UserAgent         *string `json:"userAgent" gorm:"varchar(2048);comment:用户代理" validate:"required|max_len:2048"` // 用户代理
}

// ModelEnable 启用
type ModelEnable struct {
	Enable bool `json:"enable" gorm:"type:tinyint(1); comment:启用"` // 启用
}

// ModelProcess 进度
type ModelProcess struct {
	Process ProcessType `json:"result" gorm:"type:varchar(8);not null;comment:结果" validate:"max_len:8"` // 进度
}

// ModelComment 备注
type ModelComment struct {
	Comment string `json:"comment" gorm:"type:varchar(256);comment:备注" validate:"max_len:256"` // 备注
}

// ModelCreatedAt 创建时间
type ModelCreatedAt struct {
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime:milli;comment:创建时间"` // 创建时间
}

// ModelTime 时间
type ModelTime struct {
	ModelCreatedAt

	UpdatedAt time.Time       `json:"updatedAt" gorm:"autoUpdateTime:milli;comment:更新时间"` // 更新时间
	DeletedAt *gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"comment:删除时间"`            // 删除时间
}
