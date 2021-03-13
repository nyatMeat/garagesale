package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nyatMeat/garagesale/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	//App starting
	log.Println("main: started")
	defer log.Println("main: completed")

	//Setup dependencies

	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
	case "seed":
	}

	//Start api service
	api := http.Server{
		Addr:         "localhost:20000",
		Handler:      http.HandlerFunc(ListProducts),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	//make a channel to listen for errors coming from the listener
	//Use a buffered channel so the gorutine can exit if we don't collect this error
	serverErrors := make(chan error, 1)

	//start the service listening for requests
	go func() {
		log.Printf("main: Api listering on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	//make a channel to listen for an interupt or terminate signal from OS
	//Use a buffered channel because the signal package requires it
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	//Shutdown

	//Blocking a main and waiting for shutdown
	select {
	case err := <-serverErrors:
		log.Fatalf("error: Listening and serving: %s", err)
	case <-shutdown:
		log.Println("main: Start shutdown")

		const timeout = 5 * time.Second

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		//Asking listener to shutdown and load shed
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}
		if err != nil {
			log.Fatalf("main: could not stop server gracefully: %v", err)
		}
	}
}

func openDB() (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword("postgres", "postgres"),
		Host:     "localhost",
		Path:     "postgres",
		RawQuery: q.Encode(),
	}
	return sqlx.Open("postgres", u.String())
}

type Product struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

// Echo is a basic HTTP Handler.
// If you open localhost:20000 in your browser, you may notice
// double requets being made. This happens because the browser
// sends a request in the background for a website favicon.
func ListProducts(w http.ResponseWriter, r *http.Request) {
	list := []Product{}

	if true {
		list = append(list, Product{Name: "Comic book", Cost: 85, Quantity: 50})
		list = append(list, Product{Name: "Mcdonald toys", Cost: 50, Quantity: 30})
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error marhalling", err)
		return
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("Error writing", err)
	}
}
