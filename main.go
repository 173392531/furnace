package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

var db *pgx.Conn

func main() {
	var err error

	// 连接到 PostgreSQL 数据库
	db, err = pgx.Connect(context.Background(), "postgresql://myuser:mypassword@localhost:5432/mydatabase")
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	defer db.Close(context.Background())

	// 初始化 Gin 路由
	r := gin.Default()

	// 定义一个简单的 GET 路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 定义一个使用 PostgreSQL 的路由
	r.GET("/users", getUsers)

	// 启动服务器
	r.Run(":8080")
}

// 从数据库中获取用户的处理函数
func getUsers(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		users = append(users, gin.H{"id": id, "name": name, "email": email})
	}

	c.JSON(200, users)
}
