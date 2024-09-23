/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-11-14 22:30:23
 */

package hypersonic

import (
	"github.com/gin-gonic/gin"
)

// InvokeFunc 路由处理回调
type InvokeFunc func(req *Request) (*Data, *Error)

// Router 路由
type Router struct {
	HttpMethod   string     // 方法类型 GET、POST、PUT、DELETE
	RelativePath string     // 相对路径
	Limits       []Limit    // 限制调用
	InvokeFunc   InvokeFunc // 路由处理回调
}

// 检查限制调用
func (router Router) checkLimit(req *Request, listener Listener) {
	for _, limit := range router.Limits {
		var key string

		switch limit.config.LimitType {
		case LimitIp:
			key = req.GetIp()
		case LimitUserId:
			key = req.GetUri() + "_" + req.Token.MustGetUserId()
		default:
			panic(NewError(CodeInternalError))
		}

		if !limit.Add(key) {
			listener.OnLimit(limit.config.LimitType, req)
			panic(NewErrorWithArgv(CodeRateLimit))
		}
	}
}

// 执行路由
func (router Router) run(ctx *gin.Context, i18n *I18n, adapter Listener) {
	req := newRequest(ctx)

	router.checkLimit(req, adapter)
	data, err := router.InvokeFunc(req)
	req.reply(data, err, i18n, adapter)
}

// Controller 控制器
type Controller struct {
	Path              string            // 路径
	RequestMiddleware RequestMiddleware // 中间件
	Routers           []Router          // 路由路径
}
