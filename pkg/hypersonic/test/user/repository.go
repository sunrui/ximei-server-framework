/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-08-30 23:02:38
 */

package user

import (
	"framework/pkg/mysql"
	"sync"
)

// Repository 仓库
type Repository struct {
	mysql.Repository[User] // 通用仓库
}

// 迁移一次
var migrateOnce sync.Once

// NewRepository 创建仓库
func NewRepository(sql *mysql.Mysql) Repository {
	migrateOnce.Do(func() {
		sql.AutoMigrate(&User{})
	})

	return Repository{
		Repository: mysql.NewRepository[User](sql),
	}
}
