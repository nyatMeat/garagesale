package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/nyatmeat/garagesale/internal/platform/web"
	"github.com/nyatmeat/garagesale/internal/product"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) List(w http.ResponseWriter, r *http.Request) error {

	p.Log.Println("Getting list of products")

	list, err := product.List(p.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) error {

	p.Log.Println("retrieving the single product")
	id := chi.URLParam(r, "id")
	product, err := product.Retrieve(p.DB, id)
	
	if err != nil {
		return err
	}

	return web.Respond(w, product, http.StatusOK)
}

//Create decode JSON document from a POST request and create a new product
func (p *Products) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, np); err != nil {
		return err
	}

	p.Log.Printf("creating a new product: %v", np)

	product, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		return err
	}
	return web.Respond(w, product, http.StatusCreated)
}
