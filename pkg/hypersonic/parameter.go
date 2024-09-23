/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-03 12:06:03
 */

package hypersonic

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
	"sync"
)

// ParameterPool 参数池
type ParameterPool[T any] struct {
	Pool *sync.Pool // 池
}

// NewParameterPool 创建参数池
func NewParameterPool[T any]() ParameterPool[T] {
	return ParameterPool[T]{Pool: &sync.Pool{
		New: func() any {
			var t T
			return t
		},
	}}
}

// Acquire 获得
func (parameterPool ParameterPool[T]) Acquire() T {
	return parameterPool.Pool.Get().(T)
}

// Release 释放
func (parameterPool ParameterPool[T]) Release(t T) {
	parameterPool.Pool.Put(t)
}

// https://github.com/gookit/validate/blob/master/README.zh-CN.md

// 自定义验证器
type customValidator struct{}

// ValidateStruct 验证结构体
func (customValidator) ValidateStruct(ptr any) error {
	v := validate.Struct(ptr)
	v.Validate()

	if v.IsFail() {
		return v.Errors
	}

	return nil
}

// Engine 引擎
func (customValidator) Engine() any {
	return nil
}

// 初始化
func init() {
	// 中文认证
	zhcn.RegisterGlobal()

	// 更换验证器
	binding.Validator = &customValidator{}
}
