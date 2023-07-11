package middleware

import (
	"net/http"
)

func Chain(h http.Handler, ms ...func(http.Handler) http.Handler) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}

	return h
}
