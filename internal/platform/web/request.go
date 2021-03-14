package web

import (
	"encoding/json"
	"net/http"
)

//Decode looks for a JSON document body in the request body and unmarhals it into that
func Decode(request *http.Request, value interface{}) error {
	if err := json.NewDecoder(request.Body).Decode(value); err != nil {

		return NewRequestError(err, http.StatusBadRequest)
	}

	return nil
}
