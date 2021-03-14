package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

//Decode looks for a JSON document body in the request body and unmarhals it into that
func Decode(r *http.Request, val interface{}) error{
	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
		
		return errors.Wrap(err, "decoding request body")
	}

	return nil
}