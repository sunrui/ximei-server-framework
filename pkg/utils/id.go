/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 12:23:52
 */

package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/oklog/ulid"
	"github.com/sony/sonyflake"
)

// SixNumber 6 位数字
func SixNumber() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

// NanoId nanoid
func NanoId(size int) string {
	const dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if id, err := gonanoid.Generate(dictionary, size); err != nil {
		panic(err.Error())
	} else {
		return id
	}
}

var sf *sonyflake.Sonyflake // 雪花算法

// InitSnowflake 初始化雪花
func InitSnowflake(machineId uint16) {
	var st sonyflake.Settings

	st.MachineID = func() (uint16, error) {
		return machineId, nil
	}

	sf = sonyflake.NewSonyflake(st)
}

// SnowflakeId 雪花 id
func SnowflakeId() uint64 {
	nextId, _ := sf.NextID()
	return nextId
}

// Ulid ulid
func Ulid() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	id, _ := ulid.New(ms, entropy)
	return id.String()
}

// Uuid uuid
func Uuid() string {
	id := uuid.NewString()
	id = strings.ToUpper(id)
	id = strings.ReplaceAll(id, "-", "")
	return id
}
