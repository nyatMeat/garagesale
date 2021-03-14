package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

//ctxKey represents the type of value for the context key
type ctxKey int

//KeyValues is how request values or stored/retrieved
const KeyValues ctxKey = 1

// Values carries information about each request.
type Values struct {
	StatusCode int
	Start      time.Time
}

//App is entrypooint for all web applications
type App struct {
	mux           *chi.Mux
	log           *log.Logger
	midlewareList []Middleware
}

//Handlers is the signature that all application handlers will implement
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

//NewApp knows how to construct internal state for an App
func NewApp(logger *log.Logger, mw ...Middleware) *App {

	return &App{mux: chi.NewMux(), log: logger, midlewareList: mw}
}

//Handle connects a method and URL pattern to a particular application handler
func (app *App) Handle(method, pattern string, handler Handler) {

	handler = wrapMiddleware(app.midlewareList, handler)

	fn := func(w http.ResponseWriter, r *http.Request) {

		v := Values{
			Start: time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyValues, &v)

		// Run the handler chain and catch any propagated error.
		if err := handler(ctx, w, r); err != nil {
			app.log.Printf("Unhandled error: %+v", err)
		}
	}
	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	app.mux.ServeHTTP(writer, request)
}
