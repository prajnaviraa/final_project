package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "final_project/docs"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Todo struct {
	TaskId     uint       `json:"taskId" gorm:"primary_key" example:"1"`
	TaskType   string     `json:"taskType" example:"final project"`
	TaskDesc   string     `json:"taskDesc" example:"create rest API"`
	Deadline   time.Time  `json:"deadline" example:"2021-07-26T00:00:00Z"`
	AssignedTo []Assignee `json:"assignedTo" gorm:"foreignkey:TaskID"`
}

type Assignee struct {
	PersonID uint   `json:"personID" example:"1" gorm:"primary_key"`
	Name     string `json:"assigneeName" example:"nama1"`
	Position string `json:"assigneePos" example:"position1"`
	TaskId   uint   `json:"-"`
}

var db *gorm.DB

func DBInit() {
	var err error
	db, err = gorm.Open("mysql", "root:administrator@tcp(localhost:3306)/dbTodos?parseTime=True")
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{}, &Assignee{})
}

// @title Todo Application
// @version 1.0
// @description This is a todo list management application
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email apisupport@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	DBInit()
	router := mux.NewRouter()
	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todos", createTodo).Methods("POST")
	router.HandleFunc("/todos/{taskId}", getTodo).Methods("GET")
	router.HandleFunc("/todos/{taskId}", updateTodo).Methods("PUT")
	router.HandleFunc("/todos/{taskId}", deleteTodo).Methods("DELETE")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	fmt.Println("starting web server at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// GetTodos godoc
// @Summary Get all todos
// @Description Get all todos
// @Tags todos
// @Accept  json
// @Produce  json
// @Success 200 {array} Todo
// @Router /todos [get]
func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var todos []Todo
	db.Preload("AssignedTo").Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo godoc
// @Summary Create a new Todo
// @Description Create a new Todo
// @Tags todos
// @Accept  json
// @Produce  json
// @Param todo body Todo true "Create todo"
// @Success 200 {object} Todo
// @Router /todos [post]
func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	db.Create(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// GetTodo godoc
// @Summary Get all details for a task
// @Description Get all details of a task corresponding to taskId
// @Tags todos
// @Accept  json
// @Produce  json
// @Param taskId path int true "ID of the task"
// @Success 200 {object} Todo
// @Router /todos/{taskId} [get]
func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	inputID := params["taskId"]
	var todo Todo
	db.Preload("AssignedTo").First(&todo, inputID)
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo godoc
// @Summary Update a task
// @Description Update the task corresponding to the input taskId
// @Tags todos
// @Accept  json
// @Produce  json
// @Param taskId path int true "ID of the task to be updated"
// @Success 200 {object} Todo
// @Router /todos/{taskId} [post]
func updateTodo(w http.ResponseWriter, r *http.Request) {
	var todoUpdate Todo
	json.NewDecoder(r.Body).Decode(&todoUpdate)
	db.Save(&todoUpdate)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoUpdate)
}

// DeleteTodo godoc
// @Summary Delete a task
// @Description Delete the task identified by the input taskId
// @Tags todos
// @Accept  json
// @Produce  json
// @Param taskId path int true "ID of the task to be deleted"
// @Success 204 "No Content"
// @Router /todos/{taskId} [delete]
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	inputID := params["taskId"]
	id64, _ := strconv.ParseUint(inputID, 10, 64)
	deleteId := uint(id64)
	db.Where("task_id = ?", deleteId).Delete(&Todo{})
	db.Where("task_id = ?", deleteId).Delete(&Assignee{})
	w.WriteHeader(http.StatusNoContent)
}
