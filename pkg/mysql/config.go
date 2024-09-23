/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-11-06 17:24:22
 */

package mysql

// Config 配置
type Config struct {
	User          string `json:"user"`          // 用户名
	Password      string `json:"password"`      // 密码
	Host          string `json:"host"`          // 主机
	Port          int    `json:"port"`          // 端口
	Database      string `json:"database"`      // 数据库
	MaxOpenConns  int    `json:"maxOpenConns"`  // 最大打开连接
	MaxIdleConns  int    `json:"maxIdleConns"`  // 最大空闲连接
	MaxLifetime   int    `json:"maxLifetime"`   // 最长生命周期
	MaxIdleTime   int    `json:"maxIdleTime"`   // 最大连接数
	SlowThreshold int    `json:"slowThreshold"` // 慢查询时间
}
