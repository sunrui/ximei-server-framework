/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-03-30 10:40:10
 */

package hypersonic

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequestModel 请求模型
type RequestModel struct {
	Ip       string      `json:"ip"  gorm:"type:char(15); comment:ip 地址"`                // ip 地址
	Uri      string      `json:"uri"  gorm:"type:varchar(2083); comment:访问地址"`           // 访问地址
	Method   string      `json:"method"  gorm:"type:char(14); comment:请求方式"`             // 请求方式
	Header   http.Header `json:"header" gorm:"-"`                                        // 首部
	Body     *string     `json:"body,omitempty"  gorm:"type:longtext; comment:请求体"`      // 请求体
	Response *string     `json:"response,omitempty"  gorm:"type:longtext; comment:返回结果"` // 返回结果
	UserId   *string     `json:"userId,omitempty" gorm:"type:char(16); comment:用户 id"`   // 用户 id
	Elapsed  int64       `json:"elapsed" gorm:"type:int; comment:耗时"`                    // 耗时
}

// Filed 填充请求模型
func (requestModel *RequestModel) Filed(req *Request, userId *string, response *string, elapsed int64) {
	requestModel.Ip = req.GetIp()
	requestModel.Uri = req.GetUri()
	requestModel.Method = req.GetMethod()
	requestModel.Header = req.GetAllHeader()
	requestModel.Body = req.GetBody()
	requestModel.Response = response
	requestModel.UserId = userId
	requestModel.Elapsed = elapsed
}

// String 数据
func (requestModel *RequestModel) String() string {
	buffer := strings.Builder{} // 缓存

	// method http://host:port?query protocol
	buffer.WriteString(requestModel.Method + " " + requestModel.Uri)

	// 空一行
	buffer.WriteString("\n")

	// header
	for key, values := range requestModel.Header {
		for _, value := range values {
			buffer.WriteString(key + ": " + value + "\n")
		}
	}

	// 空一行
	buffer.WriteString("\n")

	// user 信息
	if requestModel.UserId == nil {
		buffer.WriteString("userId: <null>" + "\n")
	} else {
		buffer.WriteString("userId: " + *requestModel.UserId + "\n")
	}

	// 空一行
	buffer.WriteString("\n")

	// body
	if requestModel.Body == nil {
		buffer.WriteString("<null>" + "\n")
	} else {
		buffer.WriteString(*requestModel.Body + "\n")
	}

	// 空一行
	buffer.WriteString("\n")

	// 响应
	if requestModel.Response == nil {
		buffer.WriteString("<null>" + "\n")
	} else {
		buffer.WriteString(*requestModel.Response + "\n")
	}

	// 空一行
	buffer.WriteString("\n")

	return buffer.String()
}

type Lang string

const (
	LangEn   Lang = "en"
	LangZhCN Lang = "zh-CN"
)

// Request gin 请求
type Request struct {
	ctx   *gin.Context // 上下文
	Token Token        // 令牌
}

// 创建请求
func newRequest(ctx *gin.Context) *Request {
	return &Request{
		ctx:   ctx,
		Token: newToken(ctx),
	}
}

// GetMethod 获取请求方式
func (req *Request) GetMethod() string {
	return req.ctx.Request.Method
}

// GetUri 获取访问地址
func (req *Request) GetUri() string {
	return func(ctx *gin.Context) string {
		if ctx.Request.TLS != nil {
			return "https://"
		} else {
			return "http://"
		}
	}(req.ctx) + req.ctx.Request.Host + req.ctx.Request.RequestURI
}

// GetIp 获取 ip
func (req *Request) GetIp() (ip string) {
	return req.ctx.ClientIP()
}

// Bind 绑定
func (req *Request) Bind(param any) *Error {
	var bindingType binding.Binding
	if req.ctx.Request.Method == http.MethodGet {
		bindingType = binding.Query
	} else {
		bindingType = binding.JSON
	}

	// 解析
	if err := req.ctx.ShouldBindWith(param, bindingType); err != nil {
		return NewErrorWithMessage(CodeParameterError, err.Error())
	}

	return nil
}

// GetCookie 获取 cookie
func (req *Request) GetCookie(key string) *string {
	if cookie, err := req.ctx.Cookie(key); err != nil {
		return nil
	} else {
		return &cookie
	}
}

// GetParam 获取 param
func (req *Request) GetParam(key string) string {
	return req.ctx.Param(key)
}

// GetAllHeader 获取所有 header
func (req *Request) GetAllHeader() http.Header {
	return req.ctx.Request.Header
}

// GetHeader 获取 header
func (req *Request) GetHeader(key string) *string {
	if header := req.ctx.GetHeader(key); header == "" {
		return nil
	} else {
		return &header
	}
}

// GetBody 获取 body
func (req *Request) GetBody() *string {
	if b := getBody(req.ctx); len(b) != 0 {
		bStr := fmt.Sprintf("%s", b)
		return &bStr
	} else {
		return nil
	}
}

// GetLang 获取语言
func (req *Request) GetLang() Lang {
	if lang := req.ctx.GetHeader("Accept-Language"); lang == "" {
		return LangEn
	} else if strings.Index(lang, string(LangZhCN)) != -1 {
		return LangZhCN
	}

	return LangEn
}

// GetElapsed 获取耗时
func (req *Request) GetElapsed() int64 {
	return getElapsed(req.ctx)
}

// 获取响应
func (req *Request) reply(data *Data, err *Error, i18n *I18n, listener Listener) {
	if data != nil {
		req.ctx.AbortWithStatusJSON(http.StatusOK, *data)
	}

	if err != nil {
		err.i18n(string(req.GetLang()), i18n)
		req.ctx.AbortWithStatusJSON(http.StatusBadRequest, *err)
	}

	listener.OnLog(req, data, err)
}

// TokenRequestMiddleware token 请求中间件
func TokenRequestMiddleware(req *Request) {
	_ = req.Token.MustGetUserId()
}
