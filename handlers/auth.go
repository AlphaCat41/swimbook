package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginPage(c *gin.Context) {
	// Define a custom template instance just for this request
	tmpl := template.Must(template.ParseFiles("templates/layouts/base.html", "templates/pages/login.html"))

	c.Status(http.StatusOK)
	tmpl.ExecuteTemplate(c.Writer, "layouts/base.html", gin.H{
		"title": "เข้าสู่ระบบ",
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
