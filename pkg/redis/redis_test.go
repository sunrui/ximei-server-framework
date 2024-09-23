/*
 * Copyright © 2022 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2022-11-28 15:31:21
 */

package redis

import (
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	rediz, err := New(Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
	}, 0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	n := rediz.SetIncr("hello")
	t.Log(n)
	n = rediz.SetIncr("hello")
	t.Log(n)
	n = rediz.SetDecr("hello")
	t.Log(n)
	n = rediz.SetDecr("hello")
	t.Log(n)
	n = rediz.SetDecr("hello")
	t.Log(n)

	rediz.Set("hello", []byte("world"), time.Duration(60)*time.Second)

	value, ok := rediz.Get("hello")
	t.Log(string(value), ok)

	value, ok = rediz.Get("hello-not-exist")
	t.Log(string(value), ok)

	if ok = rediz.Exists("hello"); !ok {
		t.Fatalf("hello exist")
	}

	if ok = rediz.Exists("hello-not-exist"); ok {
		t.Fatalf("hello exist")
	}

	ttl, ok := rediz.GetTtl("hello")
	t.Log(ttl, ok)

	ttl, ok = rediz.GetTtl("hello-not-exist")
	t.Log(ttl, ok)

	if ok = rediz.Del("hello"); !ok {
		t.Fatalf("delete failed")
	}

	if ok = rediz.Del("hello-not-exist"); ok {
		t.Fatalf("cannot delete")
	}

	rediz.Set("hello", []byte("world"), time.Duration(60)*time.Second)

	rediz.SetHash("hash", "hello", "world")
	value, ok = rediz.GetHash("hash", "hello")
	t.Log(string(value), ok)
}
