package main

import (
	"awesomeProject/config"
	"awesomeProject/routes"
	"awesomeProject/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := utils.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 验证数据库连接（db变量已存储在utils.DB中，供routes使用）
	if db == nil {
		log.Fatalf("数据库连接为空")
	}
	log.Println("数据库连接成功")

	// 初始化Redis
	_, err = utils.InitRedis(cfg)
	if err != nil {
		log.Printf("初始化Redis失败（将使用数据库替代）: %v", err)
	} else {
		log.Println("Redis连接成功")
	}

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 健康检查接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "ok",
		})
	})

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("服务器启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
