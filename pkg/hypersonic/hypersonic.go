/*
 * Copyright © 2022 honeysense All rights reserved.
 * Author: sunrui
 * Date: 2022-1-1 00:00:01
 */

package hypersonic

import (
	"github.com/mattn/go-colorable"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Config 配置
type Config struct {
	Listener Listener // 适配器
	I18n     *I18n    // 国际化
	IsDev    bool     // 是否开发模式
}

// Hypersonic 服务
type Hypersonic struct {
	engine   *gin.Engine // gin 引擎
	listener Listener    // 适配器
	i18n     *I18n       // 国际化
}

// New 创建
func New(config Config) (*Hypersonic, error) {
	// 创建引擎
	engine := gin.New()

	// 创建服务
	hypersonic := &Hypersonic{
		engine:   engine,
		listener: config.Listener,
		i18n:     config.I18n,
	}

	// 注册中间件
	hypersonic.registerMiddleware(engine)

	if config.IsDev {
		// 注册文档中间件
		engine.GET("/doc/*any", swaggerMiddleware)
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 开启控制台颜色
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()

	return hypersonic, nil
}

func (hypersonic *Hypersonic) registerMiddleware(engine *gin.Engine) {
	// 注册耗时中件间
	engine.Use(elapsedMiddleware)

	// 注册 body 中间件
	engine.Use(bodyMiddleware)

	// 注册安全异常中间件
	engine.Use(func(ctx *gin.Context) {
		safeRecoverMiddleware(ctx)
	})

	// 注册异常中间件
	engine.Use(func(ctx *gin.Context) {
		if err := recoverMiddleware(ctx); err != nil {
			req := newRequest(ctx)
			req.reply(nil, err, hypersonic.i18n, hypersonic.listener)
		}
	})

	// 注册 404 回调
	engine.NoRoute(notFoundMiddleware)

	// 注册 405 回调
	engine.HandleMethodNotAllowed = true
	engine.NoMethod(methodNotAllowedMiddleware)

	// 注册全局限制中间件
	engine.Use(newRateLimit(hypersonic.listener).Filter)
}

// SetMiddleware 设置中间件
func (hypersonic *Hypersonic) SetMiddleware(requestMiddleware RequestMiddleware) {
	hypersonic.engine.Use(func(ctx *gin.Context) {
		requestMiddleware(newRequest(ctx))
	})
}

// RegisterControllers 注册控制器
func (hypersonic *Hypersonic) RegisterControllers(basePath string, controllers []Controller) {
	for _, controller := range controllers {
		group := hypersonic.engine.Group(basePath + controller.Path)

		// 启用中间件
		if controller.RequestMiddleware != nil {
			group.Use(func(ctx *gin.Context) {
				controller.RequestMiddleware(newRequest(ctx))
			})
		}

		// 路由函数
		routerFunc := func(router Router) gin.HandlerFunc {
			return func(ctx *gin.Context) {
				router.run(ctx, hypersonic.i18n, hypersonic.listener)
			}
		}

		// 注册路由处理回调
		for _, router := range controller.Routers {
			switch router.HttpMethod {
			case http.MethodGet:
				group.GET(router.RelativePath, routerFunc(router))
			case http.MethodPost:
				group.POST(router.RelativePath, routerFunc(router))
			case http.MethodPut:
				group.PUT(router.RelativePath, routerFunc(router))
			case http.MethodDelete:
				group.DELETE(router.RelativePath, routerFunc(router))
			default:
				panic("hypersonic method not supported")
			}
		}
	}
}

// Run 启动服务
func (hypersonic *Hypersonic) Run(port int) error {
	return hypersonic.engine.Run(":" + strconv.Itoa(port))
}
