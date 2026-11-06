package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"swimbook/database"
	"swimbook/handlers"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(handlers.SessionCookie)
		if err != nil || token == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var userID int
		var email string
		err = database.Pool.QueryRow(c.Request.Context(), `
			SELECT u.id, u.email
			FROM sessions s
			JOIN users u ON u.id = s.user_id
			WHERE s.id = $1 AND s.expires_at > NOW()`, token,
		).Scan(&userID, &email)

		if err != nil { // ไม่พบ หรือหมดอายุ
			c.SetCookie(handlers.SessionCookie, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Next()
	}
}
