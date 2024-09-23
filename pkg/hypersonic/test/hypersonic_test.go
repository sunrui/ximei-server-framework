/*
 * Copyright © 2024 honeysense.com All rights reserved.
 * Author: sunrui
 * Date: 2024-07-13 18:02:55
 */

package test

import (
	"framework/pkg/hypersonic"
	"framework/pkg/hypersonic/test/user/user_public"
	"framework/pkg/mysql"
	"framework/pkg/redis"
	"testing"
)

func TestI18n(t *testing.T) {
	i18n, err := hypersonic.NewI18n("i18n")
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Log(i18n.Tf("en", "InternalError"))
	t.Log(i18n.T("zh-CN", "NoAuth"))
}

func TestHypersonic(t *testing.T) {
	db, err := mysql.New(mysql.Config{
		User:          "root",
		Password:      "root",
		Host:          "127.0.0.1",
		Port:          3306,
		Database:      "hypersonic_test",
		MaxOpenConns:  1,
		MaxIdleConns:  1,
		SlowThreshold: 50,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	cache, err := redis.New(redis.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
	}, 0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	i18n, err := hypersonic.NewI18n("i18n")
	if err != nil {
		t.Fatalf(err.Error())
	}

	h, err := hypersonic.New(hypersonic.Config{
		Listener: hypersonic.NewEchoListener(),
		I18n:     i18n,
		IsDev:    true,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	h.RegisterControllers("/public", []hypersonic.Controller{
		user_public.NewController(db, cache),
	})
	_ = h.Run(8080)
}
