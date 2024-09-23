/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-03 12:00:28
 */

package hypersonic

import "encoding/json"

// Page 页
type Page struct {
	Page     int `json:"page" form:"page" validate:"required|min:1|max:9999"`         // 分页，从 1 开始
	PageSize int `json:"pageSize" form:"pageSize" validate:"required|min:1|max:9999"` // 分页大小，最大 100
}

// Pagination 分页
type Pagination struct {
	Page
	TotalPage int64 `json:"totalPage" validate:"required"` // 总页数
	TotalSize int64 `json:"totalSize" validate:"required"` // 总大小
}

// M 键值
type M map[string]any

// Data 数据
type Data struct {
	Data       any         `json:"data,omitempty"`       // 数据
	Pagination *Pagination `json:"pagination,omitempty"` // 分页
}

// NewData 创建数据
func NewData(data any, pagination *Pagination) *Data {
	return &Data{
		Data:       data,
		Pagination: pagination,
	}
}

// String 字符串
func (data Data) String() string {
	dataBytes, _ := json.Marshal(data)
	return string(dataBytes)
}

// Code 错误码
type Code string

const (
	CodeOK               Code = "OK"               // 成功
	CodeNoContent        Code = "NoContent"        // 无内容
	CodeNotFound         Code = "NotFound"         // 未找到
	CodeNotMatch         Code = "NotMatch"         // 不匹配
	CodeNotImplemented   Code = "NotImplemented"   // 尚未实现
	CodeParameterError   Code = "ParameterError"   // 参数错误
	CodeConflict         Code = "Conflict"         // 数据冲突
	CodeThirdPartyError  Code = "ThirdPartyError"  // 第三方错误
	CodeInternalError    Code = "InternalError"    // 内部错误
	CodeMethodNotAllowed Code = "MethodNotAllowed" // 方法不允许
	CodeRateLimit        Code = "RateLimit"        // 速率限制
	CodeForbidden        Code = "Forbidden"        // 禁止访问
	CodeNoAuth           Code = "NoAuth"           // 未授权
)

// Error 错误
type Error struct {
	Code    Code   `json:"code,omitempty"`    // 错误码
	Message string `json:"message,omitempty"` // 错误信息
	Argv    []any  `json:"argv,omitempty"`    // 参数值
}

// NewError 创建错误
func NewError(code Code) *Error {
	return &Error{
		Code: code,
	}
}

// NewErrorWithMessage 创建错误并附加消息
func NewErrorWithMessage(code Code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithArgv 创建错误并附加参数值
func NewErrorWithArgv(code Code, argv ...any) *Error {
	return &Error{
		Code: code,
		Argv: argv,
	}
}

// String 字符串
func (err *Error) String() string {
	errBytes, _ := json.Marshal(err)
	return string(errBytes)
}

// 国际化
func (err *Error) i18n(lang string, i18n *I18n) {
	if len(err.Argv) > 0 {
		err.Message = i18n.Tf(lang, string(err.Code), err.Argv)
	} else if err.Message == "" {
		err.Message = i18n.T(lang, string(err.Code))
	}
}
