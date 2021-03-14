package product

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	ErrNotFound  = errors.New("Product not found")
	ErrInvalidId = errors.New("id provided was not a valid UUID")
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
	var product Product

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidId
	}

	const q = `SELECT product_id, quantity, cost, name, date_created, date_updated 
		FROM products
		WHERE product_id = $1`

	if err := db.Get(&product, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting products")
	}

	return &product, nil
}

//Create a new Product
func Create(db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {

	newProduct := Product{
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

	if _, err := db.Exec(q, newProduct.ID, newProduct.Name, newProduct.Cost, newProduct.Quantity, newProduct.DateCreated, newProduct.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "Inserting product: %v", np)
	}

	return &newProduct, nil
}
