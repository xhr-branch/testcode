package main

import (
	"github.com/gin-gonic/gin"
	"inventory-system/config"
	"inventory-system/controllers"
	"net/http"
)

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 创建Gin路由
	r := gin.Default()

	// 配置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 根路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Inventory System API",
			"version": "1.0.0",
			"status": "running",
		})
	})

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"database": "connected",
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 物品相关路由
		api.POST("/items", controllers.CreateItem)
		api.GET("/items", controllers.GetItems)
		api.GET("/items/:id", controllers.GetItem)
		api.PUT("/items/:id", controllers.UpdateItem)
		api.DELETE("/items/:id", controllers.DeleteItem)

		// 订单相关路由
		api.POST("/orders", controllers.CreateOrder)
		api.GET("/orders", controllers.GetOrders)
		api.GET("/orders/:id", controllers.GetOrder)
		api.PUT("/orders/:id", controllers.UpdateOrder)
		api.DELETE("/orders/:id", controllers.DeleteOrder)
	}

	// 启动服务器
	r.Run(":8080")
} 