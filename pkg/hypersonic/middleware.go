/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 13:30:09
 */

package hypersonic

import (
	"bytes"
	"fmt"
	"framework/pkg/hypersonic/swagger"
	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// RequestMiddleware 请求中间件
type RequestMiddleware func(req *Request)

const bodyTag = "BODY" // body tag

// 设置 body
func setBody(ctx *gin.Context) {
	if data, err := ctx.GetRawData(); err != nil {
		panic(err.Error())
	} else if len(data) != 0 {
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		ctx.Set(bodyTag, data)
	}
}

// 内容中间件
func bodyMiddleware(ctx *gin.Context) {
	setBody(ctx)
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Next()
}

// 获取 body
func getBody(ctx *gin.Context) []byte {
	if body, ok := ctx.Get(bodyTag); ok {
		return body.([]byte)
	} else {
		return nil
	}
}

const elapsedTag = "ELAPSED" // 耗时 tag

// 耗时中间件
func elapsedMiddleware(ctx *gin.Context) {
	ctx.Set(elapsedTag, time.Now().UnixMilli())
	ctx.Next()
}

// 获取耗时
func getElapsed(ctx *gin.Context) int64 {
	elapsed, _ := ctx.Get(elapsedTag)
	elapsed = time.Now().UnixMilli() - elapsed.(int64)
	return elapsed.(int64)
}

// 405 中间件
func methodNotAllowedMiddleware(ctx *gin.Context) {
	ctx.Abort()

	panic(NewErrorWithArgv(CodeMethodNotAllowed, ctx.Request.URL.RequestURI(), ctx.Request.Method))
}

// 404 中间件
func notFoundMiddleware(ctx *gin.Context) {
	ctx.Abort()

	panic(NewErrorWithArgv(CodeNotFound, ctx.Request.URL.RequestURI()))
}

// 全局限流
type rateLimit struct {
	lmt      *limiter.Limiter // 对象
	listener Listener         // 适配器
}

// 创建全局限流
func newRateLimit(adapter Listener) rateLimit {
	const ttl int = 5 * 60              // 活跃时间（秒）
	const capacity float64 = 5 * 60 * 5 // 容量

	lmt := tollbooth.NewLimiter(capacity, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Duration(ttl) * time.Second})
	lmt.SetIPLookups([]string{"X-Forwarded-For", "X-Real-IP", "RemoteAddr"})

	return rateLimit{
		lmt:      lmt,
		listener: adapter,
	}
}

// Filter 限流中间件
func (rateLimit rateLimit) Filter(ctx *gin.Context) {
	if err := tollbooth.LimitByRequest(rateLimit.lmt, ctx.Writer, ctx.Request); err != nil {
		ctx.Abort()

		req := newRequest(ctx)
		rateLimit.listener.OnLimit(LimitIp, req)

		panic(NewErrorWithArgv(CodeRateLimit, req.GetUri(), req.GetIp()))
	} else {
		ctx.Next()
	}
}

// 获取堆栈
func getStack(skip int, level int) []string {
	stacks := make([]string, 0)

	pc := make([]uintptr, level)
	runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc)

	for frame, ok := frames.Next(); ok; frame, ok = frames.Next() {
		stacks = append(stacks, fmt.Sprintf("%s:%d", frame.File, frame.Line))
	}

	return stacks
}

// 异常捕获中间件
func recoverMiddleware(ctx *gin.Context) (err *Error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if e, ok := recovered.(*Error); ok {
				err = e
			} else {
				err = NewErrorWithArgv(CodeInternalError, fmt.Sprintf("%+v", recovered), getStack(0, 10))
			}
		}
	}()

	ctx.Next()

	return err
}

// 安全异常捕获中间件，用于在抛出异常时触发了一个异常。
func safeRecoverMiddleware(ctx *gin.Context) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if e, ok := recovered.(*Error); ok {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, *e)
			} else {
				e = NewErrorWithArgv(CodeInternalError, fmt.Sprintf("%+v", recovered), getStack(0, 10))
				ctx.AbortWithStatusJSON(http.StatusBadRequest, *e)
			}
		}
	}()

	ctx.Next()
}

// 文档中间件
func swaggerMiddleware(ctx *gin.Context) {
	path := ctx.Request.URL.Path

	// 非 /doc 开头不是文档
	if !strings.HasPrefix(path, "/doc/") {
		return
	}

	// 过滤掉非法的 /doc/? 路径
	suffix := filepath.Base(path)
	if suffix != "doc" && suffix != "doc.json" && suffix != "redoc.js" {
		ctx.Redirect(http.StatusFound, "/doc")
		return
	}

	_, _ = ctx.Writer.Write(swagger.Redoc(suffix))
}
