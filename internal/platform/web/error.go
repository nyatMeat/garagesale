package web

//ErrorResponse how we respond to clients when something goes wrong
type ErrorResponse struct {
	Error string `json:"error"`
}

//Error is used to add web information to a request error
type Error struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}

//Returning the error string
func (e *Error) Error() string {
	return e.Err.Error()
}
