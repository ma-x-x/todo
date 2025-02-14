package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"todo/docs"
	routes "todo/internal/router"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// initHTTPServer 初始化HTTP服务器
func (a *App) initHTTPServer() error {
	// 设置 Gin 模式
	if err := a.setGinMode(); err != nil {
		return err
	}

	// 创建 Gin 引擎
	r := gin.New()

	// 初始化路由
	routes.InitRouter(a.cfg, a.services, a.rdb, r)

	// 设置 Swagger
	a.setupSwagger(r)

	// 创建 HTTP 服务器
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("服务器启动在 http://localhost:%d", a.cfg.Server.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听失败: %v", err)
		}
	}()

	return nil
}

// setGinMode 设置Gin的运行模式
func (a *App) setGinMode() error {
	log.Printf("设置 Gin 模式之前: %s", a.cfg.Server.Mode)

	if a.cfg.Server.Mode != "debug" && a.cfg.Server.Mode != "release" && a.cfg.Server.Mode != "test" {
		log.Printf("警告: 未知的服务器模式 '%s'，使用默认的 'release' 模式", a.cfg.Server.Mode)
		a.cfg.Server.Mode = "release"
	}

	log.Printf("最终使用的 Gin 模式: %s", a.cfg.Server.Mode)
	gin.SetMode(a.cfg.Server.Mode)
	return nil
}

// setupSwagger 配置Swagger文档
func (a *App) setupSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "Todo 应用后端 API 文档"
	docs.SwaggerInfo.Version = "1.0"

	if a.cfg.Server.Mode == "release" && a.cfg.Server.SwaggerHost != "" {
		if !strings.Contains(a.cfg.Server.SwaggerHost, ":") {
			docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", a.cfg.Server.SwaggerHost, a.cfg.Server.Port)
		} else {
			docs.SwaggerInfo.Host = a.cfg.Server.SwaggerHost
		}
	} else {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", a.cfg.Server.Port)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
