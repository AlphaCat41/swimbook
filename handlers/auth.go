package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	// Define a custom template instance just for this request
	tmpl := template.Must(template.ParseFiles("templates/layouts/base.html", "templates/pages/login.html"))
	query := c.DefaultQuery("pass", "0")
	statusMessage := ""
	if query == "1" {
		statusMessage = "**สมัครสมาชิกเรียบร้อย**"
	}

	c.Status(http.StatusOK)
	tmpl.ExecuteTemplate(c.Writer, "layouts/base.html", gin.H{
		"title": "เข้าสู่ระบบ",
		"alert": statusMessage,
	})
}
func RegisterPage(c *gin.Context) {
	// Define a custom template instance just for this request
	tmpl := template.Must(template.ParseFiles("templates/layouts/base.html", "templates/pages/register.html"))

	c.Status(http.StatusOK)
	tmpl.ExecuteTemplate(c.Writer, "layouts/base.html", gin.H{
		"title": "ลงทะเบียน",
	})
}
func RegisterPost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	fmt.Println(email, password)
	c.Redirect(http.StatusFound, "/login?pass=1")
}
