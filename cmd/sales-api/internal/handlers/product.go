package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/nyatmeat/garagesale/internal/platform/web"
	"github.com/nyatmeat/garagesale/internal/product"
	"github.com/pkg/errors"
)

// Products defines all of the handlers related to products. It holds the
// application state needed by the handler methods.
type Products struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List gets all products from the service layer and encodes them for the
// client response.
func (productService *Products) List(w http.ResponseWriter, r *http.Request) error {

	list, err := product.List(productService.DB)
	if err != nil {
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

// List gets all products from the service layer and encodes them for the
// client response.
func (productService *Products) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	productElement, err := product.Retrieve(productService.DB, id)

	if err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidId:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for product %q", id)
		}
	}
	return web.Respond(w, productElement, http.StatusOK)
}

//Create decode JSON document from a POST request and create a new product
func (productService *Products) Create(w http.ResponseWriter, r *http.Request) error {

	var np product.NewProduct
	if err := web.Decode(r, np); err != nil {
		return err
	}

	product, err := product.Create(productService.DB, np, time.Now())
	if err != nil {
		return err
	}
	return web.Respond(w, product, http.StatusCreated)
}
