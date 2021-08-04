package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTodos(t *testing.T) {
	DBInit()
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getTodos)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("got wrong status code: %v expected: %v",
			status, http.StatusOK)
	}
}

func TestGetTodo(t *testing.T) {
	DBInit()
	req, err := http.NewRequest("GET", "/todos/{taskId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getTodo)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("got wrong status code: %v expected: %v",
			status, http.StatusOK)
	}
}

func TestDeleteTodo(t *testing.T) {
	DBInit()
	req, err := http.NewRequest("DELETE", "/todos/{taskId}", nil)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteTodo)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusNoContent {
		t.Errorf("got wrong status code: %v expected: %v",
			status, http.StatusNoContent)
	}

}

func TestCreateTodo(t *testing.T) {
	DBInit()
	var jsonStr = []byte(`{"taskId":2,"taskType":"ass2","taskDesc":"create rest API documentation","deadline":"2021-07-26T00:00:00Z","assignedTo":[{"personID":3,"assigneeName":"nama3","assigneePos":"pos3"}]}`)

	req, err := http.NewRequest("POST", "/todo", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTodo)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("got wrong status code: %v expected: %v",
			status, http.StatusOK)
	}
}

func TestUpdateTodo(t *testing.T) {
	DBInit()
	var jsonStr = []byte(`{"taskId":2,"taskType":"ass2","taskDesc":"create rest API docs + test","deadline":"2021-07-26T00:00:00Z","assignedTo":[{"personID":4,"assigneeName":"nama4","assigneePos":"pos4"}]}`)

	req, err := http.NewRequest("PUT", "/todos/{taskId}", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTodo)
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		t.Errorf("got wrong status code: %v expected: %v",
			status, http.StatusOK)
	}
}
