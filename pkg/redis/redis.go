/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-04-16 22:43:52
 */

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis 缓存
type Redis struct {
	context context.Context // 上下文
	client  *redis.Client   // 客户端
}

// New 创建
func New(config Config, database int) (*Redis, error) {
	rdb := &Redis{
		context: context.Background(),
		client: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
			Password: config.Password,
			DB:       database,
		}),
	}

	if cmd := rdb.client.Ping(rdb.context); cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return rdb, nil
}

// SetTtl 设置过期时间
func (redis *Redis) SetTtl(key string, ttl time.Duration) (ok bool) {
	cmd := redis.client.Expire(redis.context, key, ttl)
	return cmd.Err() != nil
}

// GetTtl 获取过期时间
func (redis *Redis) GetTtl(key string) (ttl time.Duration, ok bool) {
	if cmd := redis.client.TTL(redis.context, key).Val(); cmd.Nanoseconds() > 0 {
		return cmd, true
	} else {
		return 0, false
	}
}

// Set 设置字符串
func (redis *Redis) Set(key string, value []byte, expired time.Duration) {
	if cmd := redis.client.Set(redis.context, key, value, expired); cmd.Err() != nil {
		panic(cmd.Err())
	}
}

// Get 获取字符串
func (redis *Redis) Get(key string) (value []byte, ok bool) {
	if cmd := redis.client.Get(redis.context, key); cmd.Err() != nil {
		return nil, false
	} else {
		bytes, _ := cmd.Bytes()
		return bytes, true
	}
}

// SetJson 设置 json
func (redis *Redis) SetJson(key string, value any, ttl time.Duration) {
	bytes, _ := json.Marshal(value)
	redis.Set(key, bytes, ttl)
}

// GetJson 获取 json
func (redis *Redis) GetJson(key string, value any) (ok bool) {
	var bytes []byte

	if bytes, ok = redis.Get(key); ok {
		if err := json.Unmarshal(bytes, &value); err != nil {
			panic(err.Error())
		} else {
			return true
		}
	}

	return false
}

// SetIncr 设置增加
func (redis *Redis) SetIncr(key string) int64 {
	if cmd := redis.client.Incr(redis.context, key); cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		return cmd.Val()
	}
}

// SetDecr 设置减少
func (redis *Redis) SetDecr(key string) int64 {
	if cmd := redis.client.Decr(redis.context, key); cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		return cmd.Val()
	}
}

// SetHash 设置 hash
func (redis *Redis) SetHash(hash string, key string, value any) {
	if cmd := redis.client.HSet(redis.context, hash, key, value); cmd.Err() != nil {
		panic(cmd.Err())
	}
}

// GetHash 获取 hash
func (redis *Redis) GetHash(hash string, key string) (value []byte, ok bool) {
	if cmd := redis.client.HGet(redis.context, hash, key); cmd.Err() != nil {
		panic(cmd.Err())
	} else if bytes, err := cmd.Bytes(); err != nil {
		panic(err)
	} else {
		return bytes, true
	}
}

// Exists 是否存在
func (redis *Redis) Exists(key string) bool {
	if cmd := redis.client.Exists(redis.context, key); cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		return cmd.Val() == 1
	}
}

// Del 删除
func (redis *Redis) Del(key string) bool {
	if cmd := redis.client.Del(redis.context, key); cmd.Err() != nil {
		panic(cmd.Err())
	} else {
		return cmd.Val() == 1
	}
}
