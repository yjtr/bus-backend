package main

// 测试数据库连接的脚本
// 使用方法: go run scripts/test_connection.go

import (
	"awesomeProject/config"
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

	fmt.Println("=== 数据库连接测试 ===")
	fmt.Printf("主机: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("用户: %s\n", cfg.Database.User)
	fmt.Printf("数据库: %s\n", cfg.Database.DBName)
	fmt.Println()

	// 1. 测试连接到默认的postgres数据库
	fmt.Println("步骤1: 测试连接到默认的postgres数据库...")
	defaultDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		"postgres",
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(defaultDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ 连接PostgreSQL失败: %v\n请检查:\n1. PostgreSQL服务是否运行\n2. 密码是否正确\n3. 端口是否正确", err)
	}
	fmt.Println("✅ 成功连接到PostgreSQL服务器")

	// 2. 检查目标数据库是否存在
	fmt.Printf("\n步骤2: 检查数据库 '%s' 是否存在...\n", cfg.Database.DBName)
	var exists bool
	err = db.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", cfg.Database.DBName).Scan(&exists).Error
	if err != nil {
		log.Fatalf("❌ 检查数据库失败: %v", err)
	}

	if exists {
		fmt.Printf("✅ 数据库 '%s' 已存在\n", cfg.Database.DBName)
	} else {
		fmt.Printf("❌ 数据库 '%s' 不存在\n\n", cfg.Database.DBName)
		fmt.Println("请执行以下命令创建数据库:")
		fmt.Println("---")
		fmt.Printf("方法1 (PowerShell):\n")
		fmt.Printf("  & \"C:\\Program Files\\PostgreSQL\\15\\bin\\psql.exe\" -U postgres -c \"CREATE DATABASE %s;\"\n\n", cfg.Database.DBName)
		fmt.Printf("方法2 (CMD):\n")
		fmt.Printf("  \"C:\\Program Files\\PostgreSQL\\15\\bin\\psql.exe\" -U postgres -c \"CREATE DATABASE %s;\"\n\n", cfg.Database.DBName)
		fmt.Println("方法3 (使用pgAdmin):")
		fmt.Println("  1. 打开pgAdmin")
		fmt.Println("  2. 右键 Databases -> Create -> Database...")
		fmt.Printf("  3. 输入数据库名称: %s\n", cfg.Database.DBName)
		fmt.Println("  4. 点击 Save")
		fmt.Println("---")

		// 询问是否自动创建
		fmt.Println("\n是否现在创建数据库? (y/n)")
		// 注意：在Go中自动创建数据库需要特殊权限，这里只提供提示
		fmt.Println("提示: 由于权限限制，建议使用上面的方法手动创建")
	}

	// 3. 列出所有数据库
	fmt.Println("\n步骤3: 列出所有数据库...")
	type Database struct {
		Datname string
	}
	var databases []Database
	db.Raw("SELECT datname FROM pg_database WHERE datistemplate = false ORDER BY datname").Scan(&databases)
	fmt.Println("当前数据库列表:")
	for _, db := range databases {
		marker := "  "
		if db.Datname == cfg.Database.DBName {
			marker = "✅ "
		}
		fmt.Printf("%s%s\n", marker, db.Datname)
	}

	fmt.Println("\n=== 测试完成 ===")
}
