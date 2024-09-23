/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-03-21 21:58:43
 */

package ip2region

import (
	_ "embed"
	"encoding/binary"
	"encoding/json"
	"strings"
)

// Address 位置
type Address struct {
	Country  string `json:"country" gorm:"type:varchar(32);not null;comment:国家" validate:"required|max_len:32"` // 国家
	Province string `json:"province" gorm:"type:varchar(32);not null;comment:省" validate:"required|max_len:32"` // 省
	City     string `json:"city" gorm:"type:varchar(32);not null;comment:市" validate:"required|max_len:32"`     // 市
	ISP      string `json:"isp" gorm:"type:varchar(32);not null;comment:运营商" validate:"required|max_len:32"`    // 运营商
}

// String 字符串
func (address Address) String() string {
	addressBytes, _ := json.Marshal(address)
	return string(addressBytes)
}

// 结果
func searchResult(value string) Address {
	r := strings.Split(value, "|")

	return Address{
		Country:  r[0],
		Province: r[2],
		City:     r[3],
		ISP:      r[4],
	}
}

// Str2Uint32 字符串转 ip
func Str2Uint32(netIp string) uint32 {
	return binary.BigEndian.Uint32([]byte(netIp))
}

// SearchIp 搜索 ip
func SearchIp(ip uint32) Address {
	if value, err := searcher.Search(ip); err != nil {
		panic(err.Error())
	} else {
		return searchResult(value)
	}
}

// SearchIpStr 搜索 ip 字符串
func SearchIpStr(ip string) Address {
	if value, err := searcher.SearchByStr(ip); err != nil {
		panic(err.Error())
	} else {
		return searchResult(value)
	}
}
