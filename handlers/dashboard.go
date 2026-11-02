package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DashboardPage(c *gin.Context) {
	// Define a custom template instance just for this request
	tmpl := template.Must(template.ParseFiles("templates/layouts/base.html", "templates/pages/dashboard.html"))

	c.Status(http.StatusOK)
	tmpl.ExecuteTemplate(c.Writer, "layouts/base.html", gin.H{
		"title": "Dashboard",
	})
}
