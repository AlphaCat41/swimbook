package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://myuser:mypassword@localhost:5432/swimbook_db?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("เชื่อมต่อฐานข้อมูลไม่ได้: ", err)
	}
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("ping ฐานข้อมูลไม่ได้: ", err)
	}
	Pool = pool
	log.Println("เชื่อมต่อ PostgreSQL สำเร็จ")
}

// ลบ session ที่หมดอายุทุก 1 ชั่วโมง
func StartSessionCleanup() {
	go func() { //คล้าย thread เบา ๆ ของ Go
		ticker := time.NewTicker(time.Hour)
		for range ticker.C {
			_, err := Pool.Exec(context.Background(),
				"DELETE FROM sessions WHERE expires_at < NOW()")
			if err != nil {
				log.Println("ลบ session หมดอายุไม่สำเร็จ:", err)
			}
		}
	}()
}
