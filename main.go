package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"swimbook/database"
	"swimbook/handlers"
	"swimbook/middleware"
)

func main() {
	database.Connect()
	defer database.Pool.Close()
	database.StartSessionCleanup()

	r := gin.Default()

	r.Static("/static", "./static") // เพิ่มบรรทัดนี้ (ปรับ path ให้ตรงโฟลเดอร์จริง)

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})
	r.GET("/login", handlers.LoginPage)
	r.POST("/login", handlers.LoginPost)
	r.GET("/register", handlers.RegisterPage)
	r.POST("/register", handlers.RegisterPost)
	r.GET("/logout", handlers.Logout)

	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	auth.GET("/dashboard", handlers.Dashboard)

	r.Run(":8080")
}
