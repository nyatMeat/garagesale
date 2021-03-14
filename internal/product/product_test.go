package product_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/nyatmeat/garagesale/internal/product"
	"github.com/nyatmeat/garagesale/internal/schema"
	"github.com/nyatmeat/garagesale/internal/tests"
)

func TestProducts(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	ctxt := context.Background()

	newP := product.NewProduct{
		Name:     "Comic Book",
		Cost:     10,
		Quantity: 55,
	}
	now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	p0, err := product.Create(ctxt, db, newP, now)
	if err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	p1, err := product.Retrieve(ctxt, db, p0.ID)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()
	ctxt := context.Background()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	ps, err := product.List(ctxt, db)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 2, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
