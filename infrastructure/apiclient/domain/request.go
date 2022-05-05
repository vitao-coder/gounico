package domain

import (
	"io"
	"net/http"
)

type Request interface {
	Method() string
	URL() string
	Body() io.Reader
	ToHTTP() *http.Request
}
