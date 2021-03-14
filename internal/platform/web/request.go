package web

import (
	"encoding/json"
	"net/http"
)

//Decode looks for a JSON document body in the request body and unmarhals it into that
func Decode(r *http.Request, val interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {

		return NewRequestError(err, http.StatusBadRequest)
	}

	return nil
}
