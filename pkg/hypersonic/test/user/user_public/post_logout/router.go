/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 19:47:03
 */

package post_logout

import (
	"framework/pkg/hypersonic"
	"framework/pkg/mysql"
	"framework/pkg/redis"
	"net/http"
	"time"
)

// Router 路由
type Router struct {
	nameLimit hypersonic.Limit // 用户名限流
}

// NewRouter 创建路由
func NewRouter(mysql *mysql.Mysql, redis *redis.Redis) hypersonic.Router {
	router := Router{}

	return hypersonic.Router{
		HttpMethod:   http.MethodPost,
		RelativePath: "/logout",
		Limits: []hypersonic.Limit{
			hypersonic.NewLimit(redis, hypersonic.LimitConfig{
				LimitType: hypersonic.LimitUserId,
				MaxTimes:  1,
				Interval:  time.Minute,
			}),
		},
		InvokeFunc: router.invoke,
	}
}

// 处理
func (router Router) invoke(req *hypersonic.Request) (*hypersonic.Data, *hypersonic.Error) {
	if req.Token.GetUserId() == nil {
		return nil, hypersonic.NewError(hypersonic.CodeNoAuth)
	}

	req.Token.DeleteUserId()

	return hypersonic.NewData(data{}, nil), nil
}
