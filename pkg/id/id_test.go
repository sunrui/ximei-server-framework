/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-05 12:23:52
 */

package id

import "testing"

func TestId(t *testing.T) {
	InitSnowflake(0)

	t.Log(SixNumber())
	t.Log(NanoId(8))
	t.Log(SnowflakeId())
	t.Log(Ulid())
	t.Log(Uuid())
}
