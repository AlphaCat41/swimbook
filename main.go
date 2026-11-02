package main

import (
	"fmt"
	"swimbook/handlers"

	"github.com/gin-gonic/gin"
)

func handleRequest() {
	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/", handlers.DashboardPage)
	router.GET("/login", handlers.LoginPage)

	fmt.Println("[/] Server start at http://localhost:8080/")
	router.Run(":8080")
}

func main() {
	handleRequest()
}
