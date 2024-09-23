/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-11-26 17:50:58
 */

package mysql

import (
	"golang.org/x/exp/slog"
	"time"

	"gorm.io/gorm/logger"
)

// 日志
type logWriter struct {
	mysql Mysql // mysql 数据库
}

// Printf 序列化
func (logWriter logWriter) Printf(format string, v ...any) {
	slog.Warn(format, v...)
}

// 获取日志
func getLogger(slowThreshold int, mysql Mysql) logger.Interface {
	return logger.New(
		&logWriter{
			mysql: mysql,
		},
		logger.Config{
			// 慢查询仅在 warn 级别时才会生效，默认 info 级别下全部输出
			SlowThreshold: time.Duration(slowThreshold) * time.Millisecond,
			LogLevel: func() logger.LogLevel {
				if slowThreshold <= 0 {
					return logger.Info
				} else {
					return logger.Warn
				}
			}(),
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
		},
	)
}
