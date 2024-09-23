/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-11-06 17:24:22
 */

package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Mysql 数据库
type Mysql struct {
	DB *gorm.DB // mysql 对象
}

// New 创建
func New(config Config) (*Mysql, error) {
	sql := Mysql{}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: func() logger.Interface {
			return getLogger(config.SlowThreshold, sql)
		}(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tbl_", // 表名前缀
			SingularTable: true,   // 使用单数表名
		},
	}); err == nil {
		// 配置连接池
		sqlDb, _ := db.DB()
		sqlDb.SetMaxOpenConns(config.MaxOpenConns)
		sqlDb.SetMaxIdleConns(config.MaxIdleConns)
		sqlDb.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
		sqlDb.SetConnMaxIdleTime(time.Duration(config.MaxIdleTime) * time.Second)

		sql.DB = db
		return &sql, nil
	} else {
		return nil, err
	}
}

// AutoMigrate 创建表
func (mysql Mysql) AutoMigrate(dst ...any) {
	if err := mysql.DB.AutoMigrate(dst...); err != nil {
		panic(err.Error())
	}
}
