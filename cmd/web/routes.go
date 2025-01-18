package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.HandlerFunc(http.MethodGet, "/", app.Home)
	mux.HandlerFunc(http.MethodGet, "/fetch/:id", app.FetchTask)
	mux.HandlerFunc(http.MethodGet, "/create", app.CreateTaskForm)
	mux.HandlerFunc(http.MethodPost, "/create", app.CreateTask)
	mux.HandlerFunc(http.MethodPatch, "/update", app.UpdateTask)
	mux.HandlerFunc(http.MethodGet, "/update", app.UpdateTaskForm)
	return mux
}