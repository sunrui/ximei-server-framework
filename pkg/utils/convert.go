/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 12:25:38
 */

package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

// IntToBytes 整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// BytesToInt 字节转换成整形
func BytesToInt(b []byte) int {
	var x int32
	bytesBuffer := bytes.NewBuffer(b)
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

// Trim 裁减
func Trim(str string) string {
	for _, old := range []string{
		" ", "\n", "\r", "\t",
	} {
		str = strings.Replace(str, old, "", -1)
	}

	return str
}

// TimeFormat 时间格式化
func TimeFormat(seconds int) (day, hour, minute, second int) {
	day = seconds / (24 * 3600)
	hour = (seconds - day*3600*24) / 3600
	minute = (seconds - day*24*3600 - hour*3600) / 60
	second = seconds - day*24*3600 - hour*3600 - minute*60
	return
}

// LocalTimeFormat 时间格式化字符串
func LocalTimeFormat(seconds int, isChinese bool) string {
	day, hour, minute, second := TimeFormat(seconds)

	dayString, hourString, minuteString, secondString := "天", "小时", "分钟", "秒"
	if !isChinese {
		dayString, hourString, minuteString, secondString = "day", "hour", "minute", "second"
	}

	// 格式化函数
	formatFunc := func(value int, text string) string {
		if value == 0 {
			return ""
		}

		return fmt.Sprintf(" %d %s", value, text)
	}

	return formatFunc(day, dayString) +
		formatFunc(hour, hourString) +
		formatFunc(minute, minuteString) +
		formatFunc(second, secondString)
}
