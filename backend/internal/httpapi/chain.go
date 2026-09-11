package httpapi

import "net/http"

// Chain applies middleware in the order given: the first one listed is
// outermost, seeing the request before anything else.
func Chain(h http.Handler, mw ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
