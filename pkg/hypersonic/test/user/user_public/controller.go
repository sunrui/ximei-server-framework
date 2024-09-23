/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 19:40:55
 */

package user_public

import (
	"framework/pkg/hypersonic"
	"framework/pkg/hypersonic/test/user/user_public/post_login/post_login_name"
	"framework/pkg/hypersonic/test/user/user_public/post_logout"
	"framework/pkg/mysql"
	"framework/pkg/redis"
)

// NewController 创建控制器
func NewController(mysql *mysql.Mysql, redis *redis.Redis) hypersonic.Controller {
	return hypersonic.Controller{
		Path:              "/user",
		RequestMiddleware: nil,
		Routers: []hypersonic.Router{
			post_login_name.NewRouter(mysql, redis),
			post_logout.NewRouter(mysql, redis),
		},
	}
}
