package main

import (
	"cloud.google.com/go/datastore"
	"github.com/julienschmidt/httprouter"
	"herrfennessey/book-recommender-api-v2/internal/response"
	"net/http"
	"strconv"
)

const userKind = "user-review-v1"

func (app *application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) getUserReviewCount(w http.ResponseWriter, req *http.Request) {

	bookId := httprouter.ParamsFromContext(req.Context()).ByName("bookId")
	limit := req.URL.Query().Get("limit")

	b, err := strconv.Atoi(bookId)
	if err != nil {
		app.serverError(w, req, err)
	}

	query := datastore.NewQuery(userKind).FilterField("book_id", "=", b).Namespace(app.config.env).KeysOnly()

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			app.serverError(w, req, err)
		}
		query = query.Limit(limitInt)
	}

	keys, err := app.dsClient.GetAll(req.Context(), query, nil)
	if err != nil {
		app.serverError(w, req, err)
	}

	type ApiResponse struct {
		UserCount int `json:"user_count"`
	}

	apiResponse := ApiResponse{len(keys)}

	err = response.JSON(w, http.StatusOK, apiResponse)
	if err != nil {
		app.serverError(w, req, err)
	}
}
