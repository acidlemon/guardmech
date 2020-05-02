package auth

import "net/http"

type httpError struct {
	code    int
	message string
	detail  error
}

type HttpError interface {
	Error() string
	StatusCode() int
	Detail() error
}

func NewHttpError(code int, message string, detail error) HttpError {
	return &httpError{
		code:    code,
		message: message,
		detail:  detail,
	}
}

func (e *httpError) Error() string {
	return e.message
}
func (e *httpError) StatusCode() int {
	return e.code
}
func (e *httpError) Detail() error {
	return e.detail
}

func WriteHttpError(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	httpErr, ok := err.(HttpError)
	if ok {
		statusCode = httpErr.StatusCode()
	}
	switch statusCode {
	case http.StatusUnauthorized:
		// hide failure reason
		w.WriteHeader(http.StatusUnauthorized)
	default:
		http.Error(w, err.Error(), statusCode)
	}
	return
}
