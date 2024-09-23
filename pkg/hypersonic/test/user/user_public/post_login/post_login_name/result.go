/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 19:47:03
 */

package post_login_name

import "framework/pkg/hypersonic"

type data struct {
	UserId string // 用户id
}

const (
	UserLoginRateLimit        hypersonic.Code = "UserLoginRateLimit"
	UserLoginForbidden                        = "UserLoginForbidden"
	UserLoginPasswordNotMatch                 = "UserLoginPasswordNotMatch"
	UserLoginNotFound                         = "UserLoginNotFound"
)
