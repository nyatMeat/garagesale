package product

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// List gets all Products from the database.
func List(db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT product_id, quantity, cost, name, date_created, date_updated FROM products`

	if err := db.Select(&products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Retrieve the single product
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	var p Product

	const q = `SELECT product_id, quantity, cost, name, date_created, date_updated 
		FROM products
		WHERE product_id = $1`

	if err := db.Get(&p, q, id); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return &p, nil
}

