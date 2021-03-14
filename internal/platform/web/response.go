package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

//Respond marhals a value and sends it to the client
func Respond(w http.ResponseWriter, val interface{}, statusCode int) error {

	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "marshalling to json")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(data); err != nil {
		
		return errors.Wrap(err, "writing to client")
	}
	return nil
}
