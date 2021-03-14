package mid

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/nyatmeat/garagesale/internal/platform/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(ctxt context.Context, w http.ResponseWriter, r *http.Request) error {

			v, ok := ctxt.Value(web.KeyValues).(*web.Values)

			if !ok {
				return errors.New("web values missing from context")
			}

			start := time.Now()
			// Run the handler chain and catch any propagated error.
			err := before(ctxt, w, r)
			log.Printf("%d %s %s (%v)", v.StatusCode, r.Method, r.URL.Path, time.Since(start))

			// Return the error to be handled further up to chain
			return err
		}

		return h
	}

	return f
}
