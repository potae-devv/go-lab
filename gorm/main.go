package main

import (
	"fmt"
	"strconv"

	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

func authRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return c.Next()
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	if err != nil {
		panic("failed to connect to database")
	}
	fmt.Println("Connected to the database!")

	// Migrate the schema
	err = db.AutoMigrate(&Book{}, &User{})
	if err != nil {
		panic("failed to migrate database")
	}
	fmt.Println("Database migration completed!")

	app := fiber.New()

	app.Use("/books", authRequired)

	app.Get("/books", getBooksHandler(db))
	app.Get("/books/:id", getBookByIDHandler(db))
	app.Post("/books", createBookHandler(db))
	app.Put("/books/:id", updateBookHandler(db))
	app.Delete("/books/:id", deleteBookHandler(db))

	// User API
	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		if err := CreateUser(db, user); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(201)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		token, err := Login(db, user)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 72),
			HTTPOnly: true,
		})

		return c.SendString("login successful")
	})

	app.Listen(":3000")
}

func createBookHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var book Book
		if err := c.BodyParser(&book); err != nil {
			return err
		}
		CreateBook(db, &book)
		return c.JSON(book)
	}
}

func getBooksHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(GetBooks(db))
	}
}

func getBookByIDHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("Invalid ID format")
		}

		book, err := GetBook(db, uint(id))
		if err != nil {
			return c.Status(404).SendString("Book not found")
		}
		return c.JSON(book)
	}
}

func updateBookHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("Invalid ID format")
		}

		var book Book
		if err := c.BodyParser(&book); err != nil {
			return err
		}
		book.ID = uint(id)
		UpdateBook(db, &book)
		return c.JSON(book)
	}
}

func deleteBookHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("Invalid ID format")
		}

		DeleteBook(db, uint(id))
		return c.SendStatus(204)
	}
}
