/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-04-21 00:40:44
 */

package utils

import (
	"encoding/json"
)

// 环境
type env[T any] struct {
	Environment string `json:"environment"` // 当前环境
	Dev         T      `json:"dev"`         // 开发环境
	Test        T      `json:"test"`        // 测试环境
	Prod        T      `json:"prod"`        // 生产环境
}

// NewConfig 创建配置
func NewConfig[T any](jsonByte []byte) (*T, error) {
	var e env[T]
	if err := json.Unmarshal(jsonByte, &e); err != nil {
		return nil, err
	}

	switch e.Environment {
	case "dev":
		return &e.Dev, nil
	case "test":
		return &e.Test, nil
	case "prod":
		return &e.Prod, nil
	default:
		panic("environment err")
	}
}
