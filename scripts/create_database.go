package main

// 创建数据库的辅助脚本
// 使用方法: go run scripts/create_database.go

import (
	"TapTransit-backend/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 先连接到默认的postgres数据库
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		"postgres", // 使用默认的postgres数据库
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接PostgreSQL失败: %v", err)
	}

	// 检查数据库是否存在
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", cfg.Database.DBName).Scan(&exists).Error
	if err != nil {
		log.Fatalf("检查数据库是否存在失败: %v", err)
	}

	if exists {
		log.Printf("数据库 '%s' 已存在，跳过创建", cfg.Database.DBName)
	} else {
		// 创建数据库（注意：GORM不能直接创建数据库，需要使用原生SQL）
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("获取数据库实例失败: %v", err)
		}

		// 设置数据库为单用户模式（关闭其他连接）
		sqlDB.Exec(fmt.Sprintf("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()", cfg.Database.DBName))

		// 创建数据库
		createSQL := fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName)
		err = db.Exec(createSQL).Error
		if err != nil {
			log.Fatalf("创建数据库失败: %v\n提示: 请手动执行: CREATE DATABASE %s;", err, cfg.Database.DBName)
		}
		log.Printf("数据库 '%s' 创建成功！", cfg.Database.DBName)
	}

	log.Println("数据库准备完成！现在可以运行主程序了。")
}
