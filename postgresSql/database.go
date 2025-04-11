package main

import "fmt"

func createProduct(product *Product) error {
	query := `INSERT INTO products (name, price) VALUES ($1, $2)`
	_, err := db.Exec(query, product.Name, product.Price)
	if err != nil {
		return err
	}

	fmt.Println("Product created successfully!")
	return nil
}

func getProduct(id int) (*Product, error) {
	query := `SELECT id, name, price FROM products WHERE id = $1`
	row := db.QueryRow(query, id)

	product := &Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price)

	if err != nil {
		return product, err
	}

	return product, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	var p Product
	query := `UPDATE products SET name = $1, price = $2 WHERE id = $3 RETURNING id, name, price`
	row := db.QueryRow(query, product.Name, product.Price, id)

	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, err
	}

	fmt.Println("Product updated successfully!")

	return p, err
}

func deleteProduct(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	fmt.Println("Product deleted successfully!")
	return nil
}

func getAllProducts() ([]Product, error) {
	query := `SELECT id, name, price FROM products`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
