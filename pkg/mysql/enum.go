/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-06-17 18:05:24
 */

package mysql

// AuthType 认证类型
type AuthType string

const (
	AuthName   AuthType = "NAME"   // 用户名
	AuthPhone  AuthType = "PHONE"  // 手机号
	AuthWechat AuthType = "WECHAT" // 微信
	AuthAlipay AuthType = "ALIPAY" // 支付宝
	AuthLogout AuthType = "LOGOUT" // 退出
)

// DeviceType 设备类型
type DeviceType string

const (
	DeviceAndroid      DeviceType = "ANDROID"       // 安卓
	DeviceIOS          DeviceType = "IOS"           // 苹果
	DeviceWeb          DeviceType = "WEB"           // 网页
	DeviceHtml5        DeviceType = "HTML5"         // 移动页
	DeviceWechatApplet DeviceType = "WECHAT_APPLET" // 微信小程序
)

// GenderType 性别类型
type GenderType string

const (
	GenderMale   GenderType = "MALE"   // 男性
	GenderFemale GenderType = "FEMALE" // 女性
)

// ProcessType 进度类型
type ProcessType string

const (
	ProcessWaiting ProcessType = "WAITING" // 等待
	ProcessOk      ProcessType = "OK"      // 成功
	ProcessFailed  ProcessType = "FAILED"  // 失败
)
