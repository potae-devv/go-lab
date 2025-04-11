package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	dbname   = "mydatabase"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var db *sql.DB

func main() {
	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	dbs, err := sql.Open("postgres", connStr)
	db = dbs

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database!")

	app := fiber.New()

	app.Get("/products", getProductsHandler)
	app.Get("/products/:id", getProductHandler)
	app.Post("/products", createProductHandler)
	app.Put("/products/:id", updateProductHandler)
	app.Delete("/products/:id", deleteProductHandler)

	log.Fatal(app.Listen(":3000"))
}

func getProductHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	product, err := getProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

func createProductHandler(c *fiber.Ctx) error {
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	err := createProduct(product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(product)
}

func updateProductHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	product := new(Product)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	updatedProduct, err := updateProduct(id, product)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(updatedProduct)
}

func deleteProductHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID format"})
	}
	err = deleteProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func getProductsHandler(c *fiber.Ctx) error {
	products, err := getAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(products)
}
