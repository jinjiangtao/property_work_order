package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func initDB() {
	// 数据库连接参数
	dsn := "root:root@tcp(localhost:3306)/property_work_order?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 测试数据库连接
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Database connected successfully")

	// 创建表结构
	createTables()
}

func createTables() {
	// 创建用户表
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		role ENUM('admin', 'owner') NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	// 创建保修表
	repairTable := `
	CREATE TABLE IF NOT EXISTS repairs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		location VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		image_url VARCHAR(255),
		status ENUM('pending', 'processing', 'completed') DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`

	// 执行创建表语句
	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(repairTable)
	if err != nil {
		log.Fatalf("Failed to create repairs table: %v", err)
	}

	fmt.Println("Tables created successfully")
}
