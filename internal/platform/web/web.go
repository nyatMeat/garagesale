package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

//App is entrypooint for all web applications
type App struct {
	mux *chi.Mux
	log *log.Logger
}

//Handlers is the signature that all application handlers will implement
type Handler func(http.ResponseWriter, *http.Request) error

//NewApp knows how to construct internal state for an App
func NewApp(logger *log.Logger) *App {

	return &App{mux: chi.NewMux(), log: logger}
}

//Handle connects a method and URL pattern to a particular application handler
func (app *App) Handle(method, pattern string, handler Handler) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			app.log.Printf("ERROR: %v", err)

			if err := RespondError(w, err); err != nil {
				app.log.Printf("ERROR: %v", err)
			}
		}
	}
	app.mux.MethodFunc(method, pattern, fn)
}

func (app *App) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	app.mux.ServeHTTP(writer, request)
}
