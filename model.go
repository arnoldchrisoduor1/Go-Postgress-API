package main

import (
	"database/sql"
)

// This struct will represent the product.
type product struct {
	// We will specify how these fields should be encoded when working with json.
	ID int `json:"id"`	// Integer field to represent the product ID.
	Name string `json:"name"`	// String field to represent the product name.
	Price float64 `json:"price"`	// Float string representing the product price.
}

// These functions will deal with a single product as methods on this struct.
func (p *product) getProduct(db *sql.DB) error { // Defining a method in the product struct
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
		// Returns an error if the query or scan fails otherwise returns 'nil'
}

// Function to update a product in the database.
func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
	p.Name, p.Price, p.ID) // Executes the query and assigns any error to `err`

	return err	// returns an error if any.
}

// Defines a method to delete a product from the databse.
func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID) // Executes the query ans assugns any error to 'err'.

	return err	// Returns the error.
}

// Defines a method to create a new product in th database.
func (p *product) createProduct(db *sql.DB) error {		// 
	err := db.QueryRow( 	// Executes the query and assign any error to `err`
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err // Returns an error if the query or scan fails.
	}

	return nil
}

// A standalone function that fetches a list of products.
func getProducts(db *sql.DB, start ,count int) ([]product, error) {
	// fetching products from the products table.
	rows, err := db.Query(
		// Limits the number of recprds based on the count parameter.
		// Start parameter determines how many records are skipped at the beginning.
		"SELECT id, name, price FROM products LIMIT $1 OFFSET $2", count, start)

		if err != nil {
			return nil, err
		}

		defer rows.Close() // Ensures tha the result set is closed when the function exists.

		products := []product{}	// Initializes an empty slice of `product` structs.

		for rows.Next() { // Iterating over each row in the result.
			var p product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil { // Scanning the current row into p.
				return nil, err
			}
			products = append(products, p) // Adds the products to the products slice.
		}
		return products, nil
}