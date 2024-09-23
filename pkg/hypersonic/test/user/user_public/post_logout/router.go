/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 19:47:03
 */

package post_logout

import (
	hypersonic2 "framework/pkg/hypersonic"
	"framework/pkg/mysql"
	"framework/pkg/redis"
	"net/http"
	"time"
)

// Router 路由
type Router struct {
	nameLimit hypersonic2.Limit // 用户名限流
}

// NewRouter 创建路由
func NewRouter(mysql *mysql.Mysql, redis *redis.Redis) hypersonic2.Router {
	router := Router{}

	return hypersonic2.Router{
		HttpMethod:   http.MethodPost,
		RelativePath: "/logout",
		Limits: []hypersonic2.Limit{
			hypersonic2.NewLimit(redis, hypersonic2.LimitConfig{
				LimitType: hypersonic2.LimitUserId,
				MaxTimes:  1,
				Interval:  time.Minute,
			}),
		},
		InvokeFunc: router.invoke,
	}
}

// 处理
func (router Router) invoke(req *hypersonic2.Request) (*hypersonic2.Data, *hypersonic2.Error) {
	if req.Token.GetUserId() == nil {
		return nil, hypersonic2.NewError(hypersonic2.CodeNoAuth)
	}

	req.Token.DeleteUserId()

	return hypersonic2.NewData(data{}, nil), nil
}
