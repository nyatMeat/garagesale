package product

import (
	"time"

	"github.com/google/uuid"
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

//Create a new Product
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {
	
	p := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Quantity:    np.Quantity,
		Cost:        np.Cost,
		DateCreated: now,
		DateUpdated: now,
	}

	const q = `INSERT INTO products 
	(product_id, name, cost, quantity, date_created, date_updated) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(q, p.ID, p.Name, p.Cost, p.Quantity, p.DateCreated, p.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "Inserting product: %v", np)
	}

	return &p, nil
}
