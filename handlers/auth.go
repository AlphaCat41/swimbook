package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"swimbook/database"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionCookie = "session_id"
	sessionTTL    = 7 * 24 * time.Hour
)

// สร้าง session ใหม่ใน DB แล้วส่ง cookie ให้เบราว์เซอร์
func createSession(c *gin.Context, userID int) error {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}
	token := hex.EncodeToString(b)

	_, err := database.Pool.Exec(c.Request.Context(),
		"INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)",
		token, userID, time.Now().Add(sessionTTL))
	if err != nil {
		return err
	}

	c.SetSameSite(http.SameSiteLaxMode)
	// name, value, maxAge, path, domain, secure, httpOnly
	// เปลี่ยน secure เป็น true เมื่อใช้ HTTPS
	c.SetCookie(SessionCookie, token, int(sessionTTL.Seconds()), "/", "", false, true)
	return nil
}

// render ช่วยลดโค้ดซ้ำ
func render(c *gin.Context, status int, page string, data gin.H) {
	tmpl := template.Must(template.ParseFiles(
		"templates/layouts/base.html",
		"templates/pages/"+page,
	))
	c.Status(status)
	if err := tmpl.ExecuteTemplate(c.Writer, "layouts/base.html", data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

func LoginPage(c *gin.Context) {
	alert := ""
	if c.DefaultQuery("pass", "0") == "1" {
		alert = "สมัครสมาชิกเรียบร้อย กรุณาเข้าสู่ระบบ"
	}
	render(c, http.StatusOK, "login.html", gin.H{
		"title": "เข้าสู่ระบบ",
		"alert": alert,
	})

}

func LoginPost(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	password := c.PostForm("password")

	var userID int
	var hash string
	err := database.Pool.QueryRow(c.Request.Context(),
		"SELECT id, password FROM users WHERE email = $1", email,
	).Scan(&userID, &hash)

	// ข้อความเดียวกันทั้งกรณีไม่พบอีเมลและรหัสผิด
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		render(c, http.StatusUnauthorized, "login.html", gin.H{
			"title": "เข้าสู่ระบบ",
			"error": "อีเมลหรือรหัสผ่านไม่ถูกต้อง",
		})
		return
	}

	if err := createSession(c, userID); err != nil {
		c.String(http.StatusInternalServerError, "สร้าง session ไม่สำเร็จ")
		return
	}
	c.Redirect(http.StatusFound, "/dashboard")
}

func RegisterPage(c *gin.Context) {
	render(c, http.StatusOK, "register.html", gin.H{"title": "ลงทะเบียน"})
}

func RegisterPost(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	password := c.PostForm("password")
	confirm := c.PostForm("confirm-password")

	fail := func(msg string) {
		render(c, http.StatusBadRequest, "register.html", gin.H{
			"title": "ลงทะเบียน",
			"error": msg,
			"email": email,
		})
	}

	if email == "" || !strings.Contains(email, "@") {
		fail("รูปแบบอีเมลไม่ถูกต้อง")
		return
	}
	if len(password) < 8 {
		fail("รหัสผ่านต้องมีอย่างน้อย 8 ตัวอักษร")
		return
	}
	if password != confirm {
		fail("รหัสผ่านและยืนยันรหัสผ่านไม่ตรงกัน")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "เข้ารหัสผ่านไม่สำเร็จ")
		return
	}

	_, err = database.Pool.Exec(c.Request.Context(),
		"INSERT INTO users (email, password) VALUES ($1, $2)", email, string(hash))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // email ซ้ำ
			fail("อีเมลนี้ถูกใช้งานแล้ว")
			return
		}
		c.String(http.StatusInternalServerError, "สมัครสมาชิกไม่สำเร็จ")
		return
	}

	c.Redirect(http.StatusFound, "/login?pass=1")
}

func Dashboard(c *gin.Context) {
	render(c, http.StatusOK, "dashboard.html", gin.H{
		"title": "แดชบอร์ด",
		"email": c.GetString("email"), // ถูกใส่โดย middleware
	})
}

func Logout(c *gin.Context) {
	if token, err := c.Cookie(SessionCookie); err == nil {
		database.Pool.Exec(c.Request.Context(),
			"DELETE FROM sessions WHERE id = $1", token)
	}
	c.SetCookie(SessionCookie, "", -1, "/", "", false, true) // ลบ cookie
	c.Redirect(http.StatusFound, "/login")
}
