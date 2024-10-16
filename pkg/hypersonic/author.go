/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-03 13:09:26
 */

package hypersonic

import (
	"encoding/json"
	"errors"
	"framework/pkg/redis"
	"framework/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

// Cookie cookie
type Cookie struct {
	ctx *gin.Context
}

// 创建 cookie
func newCookie(ctx *gin.Context) Cookie {
	return Cookie{ctx: ctx}
}

// Set 设置
func (cookie Cookie) Set(key string, value string, maxAge int) {
	cookie.ctx.SetCookie(key, value, maxAge, "/", "", false, true)
}

// Get 获取
func (cookie Cookie) Get(key string) string {
	// 从 cookie 中获取令牌
	value, err := cookie.ctx.Cookie(key)

	if err != nil {
		// 从 header 中获取令牌
		if value = cookie.ctx.GetHeader(key); value == "" {
			// 从 Authorization 中获取令牌
			if value = cookie.ctx.GetHeader("Authorization"); value != "" {
				prefix := "Bearer "
				if strings.Index(value, prefix) == 0 {
					value = value[len(prefix):]
				}
			}
		}
	}

	return value
}

// Delete 删除
func (cookie Cookie) Delete(key string) {
	cookie.ctx.SetCookie(key, "", -1, "/", "", false, true)
}

// 令牌存储
type tokenStorage interface {
	Set(payload any, maxAge time.Duration) (value string, err error) // 设置
	Get(key string) (payload any, ttl time.Duration, err error)      // 获取
}

// jwt 负荷
type jwtPayload struct {
	jwt.StandardClaims     // 标准
	payload            any // 负荷
}

// jwt 令牌存储
type jwtTokenStorage struct {
	ctx    *gin.Context // gin 上下文
	Secret []byte       // jwt 私钥
}

// 创建 jwt 令牌存储
func newJwtTokenStorage(secret []byte) jwtTokenStorage {
	return jwtTokenStorage{
		Secret: secret,
	}
}

// Set 设置
func (jwtTokenStorage jwtTokenStorage) Set(payload any, maxAge time.Duration) (value string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + maxAge.Milliseconds(),
		},
		payload: payload,
	}).SignedString(jwtTokenStorage.Secret)
}

// Get 设置
func (jwtTokenStorage jwtTokenStorage) Get(key string) (payload any, ttl time.Duration, err error) {
	token, err := jwt.ParseWithClaims(key, &jwtPayload{}, func(token *jwt.Token) (any, error) {
		return jwtTokenStorage.Secret, nil
	})

	if token != nil {
		if claims, ok := token.Claims.(*jwtPayload); ok && token.Valid {
			payload = &claims.payload
			ttl = time.Duration(claims.ExpiresAt-time.Now().Unix()) * time.Millisecond
			err = nil
			return
		}
	}

	return nil, 0, err
}

// redis 令牌存储
type redisTokenStorage struct {
	redis *redis.Redis // 缓存
}

// 创建 redis 令牌存储
func newRedisTokenStorage(redis *redis.Redis) redisTokenStorage {
	return redisTokenStorage{
		redis: redis,
	}
}

// Set 设置
func (redisTokenStorage redisTokenStorage) Set(payload any, maxAge time.Duration) (value string, err error) {
	key := utils.NanoId(12)
	payloadBytes, _ := json.MarshalIndent(payload, "", "\t")
	redisTokenStorage.redis.Set(key, payloadBytes, maxAge)
	return key, nil
}

// Get 获取
func (redisTokenStorage redisTokenStorage) Get(value string) (payload any, ttl time.Duration, err error) {
	if b, ok := redisTokenStorage.redis.Get(value); !ok {
		return nil, 0, errors.New("empty value")
	} else {
		if ttl, ok = redisTokenStorage.redis.GetTtl(value); !ok {
			return nil, 0, errors.New("key expired")
		}

		if err = json.Unmarshal(b, &payload); err != nil {
			return nil, 0, err
		} else {
			return payload, ttl, nil
		}
	}
}

// 令牌负荷
type tokenPayload struct {
	UserId string `json:"userId"` // 用户 id
}

// Token 令牌
type Token struct {
	ctx    *gin.Context
	cookie Cookie
	saver  tokenStorage
}

const tokenTag = "token"                                 // 令牌标志
const jwtSecret = "9B95050D-B322-49E2-BDF2-3FC9A834F779" // jwt 密钥

// 创建令牌
func newToken(ctx *gin.Context) Token {
	return Token{
		ctx:    ctx,
		cookie: newCookie(ctx),
		saver:  newJwtTokenStorage([]byte(jwtSecret)),
	}
}

// SetUserId 设置用户 id
func (token Token) SetUserId(userId string, maxAge time.Duration) {
	if value, err := token.saver.Set(tokenPayload{UserId: userId}, maxAge); err != nil {
		panic(err.Error())
	} else {
		token.cookie.Set(tokenTag, value, int(maxAge.Milliseconds()))
		token.ctx.Set(tokenTag, value)
	}
}

// GetUserId 获取用户 id
func (token Token) GetUserId() *string {
	userId := token.ctx.GetString(tokenTag)

	if userId == "" {
		if value := token.cookie.Get(tokenTag); value != "" {
			if payload, ttl, err := token.saver.Get(value); err != nil {
				panic(err.Error())
			} else if ttl > 0 {
				if p, ok := payload.(*tokenPayload); ok {
					userId = p.UserId
				}
			}
		}
	}

	if userId == "" {
		return nil
	}

	return &userId
}

// MustGetUserId 强制获取用户 id
func (token Token) MustGetUserId() string {
	if userId := token.GetUserId(); userId == nil {
		panic(NewError(CodeNoAuth))
	} else {
		return *userId
	}
}

// DeleteUserId 删除用户 id
func (token Token) DeleteUserId() {
	token.cookie.Delete(tokenTag)
}

// Filter 认证中间件
func (token Token) Filter(_ *gin.Context) {
	_ = token.MustGetUserId()
}
