package handler

import (
	"fmt"
	"html/template"
	"net/http"
)

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
	message := err.Error()
	detail := ""
	if ok {
		statusCode = httpErr.StatusCode()
		if httpErr.Detail() != nil {
			detail = httpErr.Detail().Error()
		}
	}

	tmpl := template.Must(template.New("errorHTML").Parse(errorHTML))

	switch statusCode {
	case http.StatusUnauthorized:
		// hide failure reason
		w.WriteHeader(http.StatusUnauthorized)
		break
	default:
		w.WriteHeader(statusCode)
		tmpl.Execute(w, map[string]string{
			"StatusCode":    fmt.Sprintf("%d", statusCode),
			"StatusMessage": http.StatusText(statusCode),
			"ErrorMessage":  message,
			"ErrorDetail":   detail,
		})
	}
	return
}

const errorHTML = `<!doctype html>
<html>
<head>
  <title>{{ .StatusCode }} {{ .StatusMessage }} - guardmech</title>
</head>
<body>
<h1>{{ .StatusMessage }}</h1>
<p>Error: {{ .ErrorMessage }}</p>
<p>Detail: {{ .ErrorDetail }}</p>

<p><a href="/">back</a></p>
</body>
</html>
`
