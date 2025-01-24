package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"web.taswiya-todo.cc/models"
)

type data struct {
	Todo  *models.ToDo
	ToDos []*models.ToDo
}

func (app *application) Home(response http.ResponseWriter, request *http.Request) {
	ts, err := template.ParseFiles("./ui/home.html")
	if err != nil {
		app.logger.Printf("failed to parse the home template: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	todos, err := app.models.FetchAll()
	if err != nil {
		app.logger.Printf("failed to fetch the todos from the database: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	dt := &data{
		ToDos: todos,
	}

	if err := ts.Execute(response, dt); err != nil {
		app.logger.Printf("failed to execute the home template: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func (app *application) CreateTaskForm(response http.ResponseWriter, request *http.Request) {
	ts, err := template.ParseFiles("./ui/create.html")
	if err != nil {
		app.logger.Printf("failed to parse the create template: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if err := ts.Execute(response, nil); err != nil {
		app.logger.Printf("failed to execute the create template: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func (app *application) CreateTask(response http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(response, "Bad Request", http.StatusBadRequest)
		return
	}

	title := request.PostForm.Get("title")
	description := request.PostForm.Get("description")
	priority := request.PostForm.Get("priority")
	temp := request.PostForm.Get("deadline")
	// app.logger.Println(temp)
	parts := strings.Split(temp, "T")
	// app.logger.Println(parts)
	deadline, err := time.Parse("2006-01-02 15:04", parts[0]+" "+parts[1])
	if err != nil {
		http.Error(response, "Bad Request", http.StatusBadRequest)
		return
	}

	todo := &models.ToDo{
		Title:       title,
		Description: description,
		Priority:    priority,
		Deadline:    deadline,
		Done:        false,
		UpdatedAt:   time.Now(),
	}

	if err := app.models.Create(todo); err != nil {
		app.logger.Printf("failed to create a new task: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	fmt.Fprintln(response, "Congratulations!")
}

func (app *application) UpdateTask(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Task has been updated...\n")
}

func (app *application) UpdateTaskForm(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Task has been updated...\n")
}

func (app *application) FetchTask(response http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(response, "Task Not Found", http.StatusNotFound)
		return
	}

	todo, err := app.models.Fetch(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrTaskNotFound):
			http.Error(response, "Task Not Found", http.StatusNotFound)
		default:
			app.logger.Printf("failed to fetch the task: %s", err)
			http.Error(response, "Internal Server Error", http.StatusInternalServerError)
		}

		return
	}

	ts, err := template.ParseFiles("./ui/task.html")
	if err != nil {
		app.logger.Printf("failed to parse the task html: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	if err := ts.ExecuteTemplate(response, "task", todo); err != nil {
		app.logger.Printf("failed to execute the task template: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func (app *application) DeleteTask(response http.ResponseWriter, request *http.Request) {
	params := httprouter.ParamsFromContext(request.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.Error(response, "Task Not Found", http.StatusNotFound)
		return
	}

	if err := app.models.Delete(id); err != nil {
		app.logger.Printf("failed to delete the task from the database: %s", err)
		http.Error(response, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	fmt.Fprintf(response, "Success")
}
