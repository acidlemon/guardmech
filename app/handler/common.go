package handler

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/acidlemon/guardmech/app/usecase"
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

func NewHttpErrorFromErr(err error) HttpError {
	if uerr, ok := err.(usecase.Error); ok {
		code := http.StatusInternalServerError
		switch uerr.Type() {
		case usecase.AuthError:
			code = http.StatusUnauthorized
			break
		case usecase.SecurityError:
			code = http.StatusForbidden
			break
		case usecase.VerificationError:
			code = http.StatusForbidden
			break
		}
		return &httpError{
			code:    code,
			message: uerr.Error(),
			detail:  uerr.Detail(),
		}
	}

	return &httpError{
		code:    http.StatusInternalServerError,
		message: err.Error(),
		detail:  nil,
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
  <link rel="stylesheet" href="https://newcss.net/new.min.css">
</head>
<body>
<header>
  <h1>Authorization Required</h1>
</header>
<main>
  <h2>{{ .StatusCode }} {{ .StatusMessage }}</h2>
  <blockquote>
    <dl>
      <dt>Reason</dt><dd>{{ .ErrorMessage }}</dd>
    </dl>
    {{ if .ErrorDetail }}
    <dl>
      <dt>Detail</dt><dd>{{ .ErrorDetail }}</dd>
    </dl>
    {{ end }}
  </blockquote>
  <p><a href="/">back</a></p>
</main>
</body>
</html>
`
