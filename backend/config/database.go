package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"inventory-system/models"
)

var DB *gorm.DB

func InitDB() {
	// 修改这里的用户名和密码为你的MySQL配置
	dsn := "xhr:xhr@tcp(127.0.0.1:3306)/inventory?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// 自动迁移数据库表
	db.AutoMigrate(&models.Item{}, &models.Order{})

	DB = db
	fmt.Println("Database connected successfully")
} 