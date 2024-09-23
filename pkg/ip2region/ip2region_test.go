/*
 * Copyright © 2023 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2023-03-21 22:48:36
 */

package ip2region

import (
	"testing"
)

func Test(t *testing.T) {
	t.Logf("%d", Str2Uint32("117.136.111.22"))

	ips := []string{
		"127.0.0.1",
		"183.137.26.235",
		"117.136.111.22",
	}

	for _, ip := range ips {
		address := SearchIpStr(ip)
		t.Logf(address.String())
	}
}
