package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFound)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	router.HandlerFunc("GET", "/health", app.status)
	router.HandlerFunc("GET", "/users/book-popularity/:bookId", app.getUserReviewCount)

	return app.recoverPanic(router)
}
