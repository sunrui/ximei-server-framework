/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-03 13:08:50
 */

package hypersonic

import (
	"fmt"
	"framework/pkg/redis"
	"strconv"
	"time"
)

// LimitType 限制类型
type LimitType string

const (
	LimitIp     LimitType = "Ip"     // ip
	LimitUserId LimitType = "UserId" // 用户 id
	LimitName   LimitType = "Name"   // 名称
	LimitPhone  LimitType = "Phone"  // 手机号
)

// LimitConfig 限制配置
type LimitConfig struct {
	LimitType LimitType     // 限制类型
	MaxTimes  int64         // 最大次数
	Interval  time.Duration // 限制间隔
}

// Limit 限制
type Limit struct {
	redis  *redis.Redis // redis
	config LimitConfig  // 限制配置
}

// NewLimit 创建限制
func NewLimit(redis *redis.Redis, config LimitConfig) Limit {
	return Limit{
		redis:  redis,
		config: config,
	}
}

// 获取格式化 Key
func (limit Limit) getFormatKey(key string) string {
	return fmt.Sprintf("Limit:%s:%s", limit.config.LimitType, key)
}

// Add 增加
func (limit Limit) Add(key string) bool {
	formatKey := limit.getFormatKey(key)
	exist := limit.redis.Exists(formatKey)
	times := limit.redis.SetIncr(formatKey)

	if !exist {
		limit.redis.SetTtl(formatKey, limit.config.Interval)
	}

	if times > limit.config.MaxTimes {
		return false
	}

	return true
}

// Left 剩余次数
func (limit Limit) Left(key string) int64 {
	formatKey := limit.getFormatKey(key)

	if timesByte, ok := limit.redis.Get(formatKey); ok {
		times, _ := strconv.Atoi(string(timesByte))
		if limit.config.MaxTimes-int64(times) < 0 {
			return 0
		}

		return limit.config.MaxTimes - int64(times)
	} else {
		return limit.config.MaxTimes
	}
}

// Reset 重置
func (limit Limit) Reset(key string) {
	formatKey := limit.getFormatKey(key)

	limit.redis.Del(formatKey)
}
