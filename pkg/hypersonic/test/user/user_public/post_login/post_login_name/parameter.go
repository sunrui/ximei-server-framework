/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 19:47:03
 */

package post_login_name

import (
	"framework/pkg/mysql"
)

// 参数
type parameter struct {
	Refer mysql.ModelRefer // 来源

	Name     string `json:"name" validate:"required|ascii|min_len:2|max_len:32"`     // 用户名
	Password string `json:"password" validate:"required|ascii|min_len:6|max_len:32"` // 密码
}
