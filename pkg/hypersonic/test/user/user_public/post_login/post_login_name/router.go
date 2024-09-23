/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 19:47:03
 */

package post_login_name

import (
	"framework/pkg/hypersonic"
	"framework/pkg/hypersonic/test/user"
	"framework/pkg/mysql"
	"framework/pkg/redis"
	"net/http"
	"strings"
	"time"
)

// Router 路由
type Router struct {
	parameterPool  hypersonic.ParameterPool[parameter] // 参数池
	nameLimit      hypersonic.Limit                    // 用户名限流
	userRepository user.Repository                     // 用户仓库
}

// NewRouter 创建路由
func NewRouter(mysql *mysql.Mysql, redis *redis.Redis) hypersonic.Router {
	router := Router{
		parameterPool: hypersonic.NewParameterPool[parameter](),
		nameLimit: hypersonic.NewLimit(redis, hypersonic.LimitConfig{
			LimitType: hypersonic.LimitName,
			MaxTimes:  1,
			Interval:  time.Minute,
		}),
		userRepository: user.NewRepository(mysql),
	}

	return hypersonic.Router{
		HttpMethod:   http.MethodPost,
		RelativePath: "/login/name",
		Limits: []hypersonic.Limit{
			hypersonic.NewLimit(redis, hypersonic.LimitConfig{
				LimitType: hypersonic.LimitIp,
				MaxTimes:  1,
				Interval:  5 * time.Second,
			}),
		},
		InvokeFunc: router.invoke,
	}
}

// 处理
func (router Router) invoke(req *hypersonic.Request) (*hypersonic.Data, *hypersonic.Error) {
	// 参数
	param := router.parameterPool.Acquire()
	defer router.parameterPool.Release(param)

	if err := req.Bind(&param); err != nil {
		return nil, err
	}

	var userId string // userId

	// 大小写不敏感
	param.Name = strings.ToLower(param.Name)

	// 判断用户名称是否已经登录限制
	if router.nameLimit.Left(param.Name) == 0 {
		return nil, hypersonic.NewError(UserLoginRateLimit)
	}

	// 根据用户名查找用户
	if userOne := router.userRepository.FindOne("name = ?", param.Name); userOne != nil {
		// 判断用户名是否已经禁用
		if !userOne.Enable {
			return nil, hypersonic.NewError(UserLoginForbidden)
		}

		// 记录登录次数
		isAddOk := router.nameLimit.Add(param.Name)

		// 验证密码成功
		if userOne.IsValidatePassword(param.Password) {
			userId = userOne.Id
		} else {
			if isAddOk {
				return nil, hypersonic.NewError(UserLoginPasswordNotMatch)
			} else {
				return nil, hypersonic.NewError(UserLoginRateLimit)
			}
		}
	} else {
		// 没有当前用户
		return nil, hypersonic.NewError(UserLoginNotFound)
	}

	// 写入令牌
	req.Token.SetUserId(userId, 30*24*time.Hour)

	// 异步记录
	defer func() {

	}()

	return hypersonic.NewData(data{
		UserId: userId,
	}, nil), nil
}
