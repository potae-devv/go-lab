package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Read database configuration from .env file
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT")) // Convert port to int
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Configure your PostgreSQL database details here
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Book{})

	// Setup Fiber
	app := fiber.New()

	// CRUD routes
	app.Get("/books", func(c *fiber.Ctx) error {
		return getBooks(db, c)
	})
	app.Get("/books/:id", func(c *fiber.Ctx) error {
		return getBook(db, c)
	})
	app.Post("/books", func(c *fiber.Ctx) error {
		return createBook(db, c)
	})
	app.Put("/books/:id", func(c *fiber.Ctx) error {
		return updateBook(db, c)
	})
	app.Delete("/books/:id", func(c *fiber.Ctx) error {
		return deleteBook(db, c)
	})

	// Start server
	log.Fatal(app.Listen(":8000"))
}
