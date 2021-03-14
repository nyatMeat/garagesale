package product

import (
	"context"
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
func List(ctxt context.Context, db *sqlx.DB) ([]Product, error) {
	products := []Product{}

	const q = `SELECT
		p.product_id, p.quantity, p.cost, p.name,
		COALESCE(SUM(s.paid), 0) as revenue,
		COALESCE(SUM(s.quantity), 0) as sold,
		p.date_created, p.date_updated 
	 
	 FROM products AS p
	 LEFT JOIN sales AS s ON p.product_id = s.product_id 
	 GROUP BY p.product_id
	`

	if err := db.SelectContext(ctxt, &products, q); err != nil {
		return nil, errors.Wrap(err, "selecting products")
	}

	return products, nil
}

// Retrieve the single product
func Retrieve(ctxt context.Context, db *sqlx.DB, id string) (*Product, error) {
	var product Product

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidId
	}

	const q = `SELECT 
			p.product_id, p.quantity, p.cost, p.name,
			COALESCE(SUM(s.paid), 0) as revenue,
			COALESCE(SUM(s.quantity), 0) as sold,
			p.date_created, p.date_updated 

		FROM products AS p 
		LEFT JOIN sales AS s ON s.product_id = p.product_id
		WHERE p.product_id = $1
		GROUP BY p.product_id
		`

	if err := db.GetContext(ctxt, &product, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "selecting products")
	}

	return &product, nil
}

//Create a new Product
func Create(ctxt context.Context, db *sqlx.DB, np NewProduct, now time.Time) (*Product, error) {

	newProduct := Product{
		ID:          uuid.New().String(),
		Name:        np.Name,
		Quantity:    np.Quantity,
		Cost:        np.Cost,
		DateCreated: now.UTC(),
		DateUpdated: now.UTC(),
	}

	const q = `INSERT INTO products 
	(product_id, name, cost, quantity, date_created, date_updated) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.ExecContext(ctxt, q, newProduct.ID, newProduct.Name, newProduct.Cost, newProduct.Quantity, newProduct.DateCreated, newProduct.DateUpdated); err != nil {
		return nil, errors.Wrapf(err, "Inserting product: %v", np)
	}

	return &newProduct, nil
}
