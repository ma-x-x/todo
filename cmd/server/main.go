package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo/cmd/server/app"
	"todo/pkg/config"
)

// Swagger API 文档注解
// @title Todo API
// @version 1.0
// @description 这是一个待办事项管理系统的API服务
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @produce application/json
// @consume application/json

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请求头需要添加Bearer token

// main 函数是程序的入口点
// 它负责初始化和启动整个应用程序，包括：
// 1. 加载配置文件
// 2. 创建应用实例
// 3. 初始化所有组件
// 4. 处理优雅关闭
func main() {
	log.Println("=== Todo 应用程序启动 ===")

	// 1. 加载配置
	log.Println("正在加载配置文件...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}
	log.Println("✅ 配置加载成功")

	// 2. 创建应用实例
	log.Println("正在创建应用实例...")
	application := app.New(cfg)
	log.Println("✅ 应用实例创建成功")

	// 3. 初始化应用
	log.Println("正在初始化应用...")
	if err := application.Initialize(); err != nil {
		log.Fatalf("❌ 初始化应用失败: %v", err)
	}
	log.Println("✅ 应用初始化完成")
	log.Println("=== 服务已启动并开始接收请求 ===")

	// 4. 设置优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("收到系统信号 [%v]，开始执行优雅关闭...", sig)

	// 5. 执行优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Shutdown(ctx); err != nil {
		log.Fatalf("❌ 关闭应用失败: %v", err)
	}

	log.Println("✅ 应用已成功关闭")
	log.Println("=== Todo 应用程序退出 ===")
}
