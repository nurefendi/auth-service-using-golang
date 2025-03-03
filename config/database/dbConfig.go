package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Global connection jika Prefork dimatikan
var DBConn *gorm.DB

// Function untuk membuat koneksi database baru
func CreateDBConnection() *gorm.DB {
	log.Info("Initiating database connection...")

	// Load environment variables
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	host := os.Getenv("DB_HOST")

	// Create DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, database)

	// Open connection
	conn, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Get sql.DB object
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalf("Cannot get sql.DB: %v", err)
	}

	// Set database connection pool settings
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(parseStringToInt(os.Getenv("DB_MAX_IDLE_CONNECTION")))
	sqlDB.SetMaxOpenConns(parseStringToInt(os.Getenv("DB_MAX_OPEN_CONNECTION")))
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("Database connection established")
	return conn
}

// Function untuk mendapatkan koneksi database dari Fiber Context
func GetDBConnection(c *fiber.Ctx) *gorm.DB {
	db, ok := c.Locals("db").(*gorm.DB)
	if ok {
		return db // Koneksi dari Prefork
	}
	return DBConn // Koneksi Global
}

// Middleware untuk setup DB di setiap request (hanya jika pakai Prefork)
func DBMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.App().Config().Prefork {
			db := CreateDBConnection()
			defer func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			}()
			c.Locals("db", db)
		}
		return c.Next()
	}
}

// Fungsi untuk inisialisasi koneksi global saat tidak menggunakan Prefork
func InitGlobalDB() {
	if DBConn == nil {
		DBConn = CreateDBConnection()
	}
}

// Function untuk menutup koneksi DB saat aplikasi berhenti
func CloseDBConnection() {
	if DBConn != nil {
		sqlDB, err := DBConn.DB()
		if err != nil {
			log.Errorf("Error while closing DB connection: %v", err)
			return
		}
		sqlDB.Close()
		log.Info("Global database connection closed")
	}
}

// Function untuk mengkonversi string ke int
func parseStringToInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Errorf("Error parsing int: %v", err)
		return 0
	}
	return i
}
