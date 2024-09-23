/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 13:47:37
 */

package swagger

import (
	_ "embed"
	"os"
)

//go:embed swagger_redoc.js
var swaggerRedocJsByte []byte

//go:embed swagger.html
var swaggerHtmlByte []byte

// Redoc 文档
func Redoc(suffix string) []byte {
	if suffix == "doc.json" {
		data, _ := os.ReadFile("docs/swagger.json")
		return data
	}

	if suffix == "redoc.js" {
		return swaggerRedocJsByte
	} else {
		return swaggerHtmlByte
	}
}
