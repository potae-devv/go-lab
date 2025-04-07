package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	_ "github.com/potae-dev/fiber-test/docs"
)

// loggingMiddleware logs the processing time for each request
func loggingMiddleware(c *fiber.Ctx) error {
	// Start timer
	start := time.Now()

	// Process request
	err := c.Next()

	// Calculate processing time
	duration := time.Since(start)

	// Log the information
	fmt.Printf("Request URL: %s - Method: %s - Duration: %s\n", c.OriginalURL(), c.Method(), duration)

	return err
}

func isAdmin(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)

	claims := user.Claims.(jwt.MapClaims)
	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

// Book struct to hold book data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book // Slice to store books

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Initialize in-memory data
	books = append(books, Book{ID: 1, Title: "1984", Author: "George Orwell"})
	books = append(books, Book{ID: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"})

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Login route
	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}))

	app.Use(isAdmin)
	app.Use(loggingMiddleware)

	app.Get("/book", getBooks)
	app.Get("/book/:id", getBook)
	app.Post("/book", createBook)
	app.Put("/book/:id", updateBook)
	app.Delete("/book/:id", deleteBook)

	app.Post("/upload", uploadImage)

	app.Get("/test-html", renderTemplate)

	app.Get("/api/config", getConfig)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app.Listen(":" + port)
}

// Dummy user for example
var user = struct {
	Email    string
	Password string
}{
	Email:    "user@example.com",
	Password: "password123",
}

func login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var request LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	// Check credentials - In real world, you should check against a database
	if request.Email != user.Email || request.Password != user.Password {
		return fiber.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = request.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func uploadImage(c *fiber.Ctx) error {
	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Save the file to the server
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File uploaded successfully: " + file.Filename)
}

func renderTemplate(c *fiber.Ctx) error {
	// Render the template with variable data
	return c.Render("index", fiber.Map{
		"Name": "Potae",
	})
}

func getConfig(c *fiber.Ctx) error {
	// Example: Return a configuration value from environment variable
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "defaultSecret" // Default value if not specified
	}

	return c.JSON(fiber.Map{
		"secret_key": secretKey,
	})
}
