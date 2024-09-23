/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-04-21 00:40:44
 */

package utils

import (
	_ "embed"
	"encoding/json"
	"testing"
)

// RedisConfig 缓存配置
type RedisConfig struct {
	Host     string // 主机
	Port     int    // 端口
	Password string // 密码
}

//go:embed config_test.json
var configTestJsonBytes []byte

// TestNewConfig 测试
func TestNewConfig(t *testing.T) {
	if config, err := NewConfig[RedisConfig](configTestJsonBytes); err != nil {
		t.Error(err.Error())
	} else {
		configBytes, _ := json.Marshal(config)
		t.Log(string(configBytes))
	}
}
