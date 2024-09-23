/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-04-23 00:33:14
 */

package area

import (
	_ "embed"
	"encoding/json"
)

//go:embed area.json
var areaJsonByte []byte

// 国家
var country Country

// 初始化
func init() {
	if err := json.Unmarshal(areaJsonByte, &country); err != nil {
		panic(err.Error())
	}
}
