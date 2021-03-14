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
func (p *Products) List(w http.ResponseWriter, r *http.Request) {

	p.Log.Println("Getting list of products")

	list, err := product.List(p.DB)
	if err != nil {
		p.Log.Printf("error: listing products: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Respond(w, list, http.StatusOK); err != nil {
		p.Log.Println("error marshalling result", err)
		return
	}
}

// List gets all products from the service layer and encodes them for the
// client response.
func (p *Products) Retrieve(w http.ResponseWriter, r *http.Request) {

	p.Log.Println("retrieving the single product")
	id := chi.URLParam(r, "id")
	product, err := product.Retrieve(p.DB, id)
	if err != nil {
		p.Log.Printf("error: retrieving product: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := web.Respond(w, product, http.StatusOK); err != nil {
		p.Log.Println("error marshalling result", err)
		return
	}
}

//Create decode JSON document from a POST request and create a new product
func (p *Products) Create(w http.ResponseWriter, r *http.Request) {

	var np product.NewProduct
	if err := web.Decode(r, np); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println(err)
		return
	}

	p.Log.Printf("creating a new product: %v", np)

	product, err := product.Create(p.DB, np, time.Now())
	if err != nil {
		p.Log.Printf("error: creating new product: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := web.Respond(w, product, http.StatusCreated); err != nil {
		p.Log.Println("error marshalling result", err)
		return
	}
}
