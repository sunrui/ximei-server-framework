/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-04-23 00:31:56
 */

package ip2region

import (
	_ "embed"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

// https://github.com/lionsoul2014/ip2region/tree/master/binding/golang

// 搜索
var searcher *xdb.Searcher

//go:embed ip2region.xdb
var ip2regionXdbByte []byte

// 初始化
func init() {
	var err error
	if searcher, err = xdb.NewWithBuffer(ip2regionXdbByte); err != nil {
		panic(err.Error())
	}
}
