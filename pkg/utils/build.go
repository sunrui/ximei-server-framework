/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 12:27:30
 */

package utils

import (
	"flag"
	"testing"
)

// 当前环境
var build *string

// IsDev 是否为开发环境
func IsDev() bool {
	return build != nil && *build != "prod"
}

// 初始化
func init() {
	testing.Init()

	// 解析参数，如 -build prod
	flag.Parse()
	build = flag.String("build", "dev", "编译类型")
}
