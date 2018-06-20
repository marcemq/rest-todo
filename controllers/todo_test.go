package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/marcemq/rest-todo/models"
	"github.com/marcemq/rest-todo/utils"
	"gopkg.in/mgo.v2/bson"
)

func TestNewTodoController(t *testing.T) {
	tt := []struct {
		name  string
		dburl string
		err   error
	}{
		{name: "New controller success", dburl: utils.DBurl},
		{name: "New controller failure", dburl: "myDBurl", err: errors.New(fmt.Sprintf("Can't work w/o valid mgo db session!"))},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTodoController(utils.GetSession(tc.dburl))
			if err != tc.err {
				t.Fatalf("Error failure-> exp: %v, got: %v", tc.err, err)
			}
		})
	}
}

func TestCreateTodo(t *testing.T) {
	tt := []struct {
		name        string
		todo        string
		expHttpCode int
	}{
		{name: "Create success TODO", todo: "My secret TODO to be created", expHttpCode: 200},
		{name: "Create failure-empty TODO", todo: "", expHttpCode: 400},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tctrl, _ := NewTodoController(utils.GetSession(utils.DBurl))
			router := httprouter.New()
			router.POST("/todo", tctrl.CreateTodo)

			data := map[string]string{"todo": tc.todo}
			dataj, _ := json.Marshal(data)
			todourl := "http://" + SRVADDR + "/todo"
			req, err := http.NewRequest("POST", todourl, bytes.NewBuffer(dataj))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create POST request: %v", err)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if status := rec.Code; status != tc.expHttpCode {
				t.Fatalf("Wrong request status, expected %v:got %v", tc.expHttpCode, status)
			}

			resp := rec.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Could not read body response: %v", err)
			}
			strbody := string(bytes.TrimSpace(body))
			t.Log(strbody)
			if !strings.Contains(strbody, tc.todo) {
				t.Fatalf("Task:%v not in server response, got %v", tc.todo, strbody)
			}
		})
	}
}

func TestGetTodoList(t *testing.T) {
	tt := []struct {
		name        string
		expHttpCode int
	}{
		{name: "List all content", expHttpCode: 200},
		//{name: "List err", expHttpCode: 404},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tctrl, _ := NewTodoController(utils.GetSession(utils.DBurl))
			router := httprouter.New()
			router.GET("/list", tctrl.GetTodoList)

			todourl := "http://" + SRVADDR + "/list"
			req, err := http.NewRequest("GET", todourl, nil)
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create GET request: %v", err)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if status := rec.Code; status != tc.expHttpCode {
				t.Fatalf("Wrong request status, expected %v:got %v", tc.expHttpCode, status)
			}

			resp := rec.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Could not read body response: %v", err)
			}
			strbody := string(bytes.TrimSpace(body))
			if len(strbody) == 0 {
				t.Fatalf("Could not get TODO list")
			}
		})
	}
}

func TestGetTodo(t *testing.T) {
	tt := []struct {
		name        string
		newTodo     string
		expHttpCode int
	}{
		{name: "Get success TODO", newTodo: "My secret TODO to be created", expHttpCode: 200},
		{name: "Get failure-empty TODO", newTodo: "", expHttpCode: 400},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tctrl, _ := NewTodoController(utils.GetSession(utils.DBurl))
			router := httprouter.New()
			router.POST("/todo", tctrl.CreateTodo)

			data := map[string]string{"todo": tc.newTodo}
			dataj, _ := json.Marshal(data)
			todourl := "http://" + SRVADDR + "/todo"
			req, err := http.NewRequest("POST", todourl, bytes.NewBuffer(dataj))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create POST request: %v", err)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			tododb := models.Todo{}
			json.NewDecoder(rec.Body).Decode(&tododb)
			t.Log(tododb.Id)
			todoId := bson.ObjectId(tododb.Id).Hex()
			router.GET("/todo/:id", tctrl.GetTodo)
			todourl = "http://" + SRVADDR + "/todo/:id"
			req, err = http.NewRequest("GET", todourl, strings.NewReader(todoId))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create GET request: %v", err)
			}

			rec = httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if status := rec.Code; status != tc.expHttpCode {
				t.Fatalf("Wrong request status, expected %v:got %v", tc.expHttpCode, status)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	tt := []struct {
		name        string
		newTodo     string
		expHttpCode int
	}{
		{name: "Delete success TODO", newTodo: "My secret TODO to be delete", expHttpCode: 200},
		{name: "Delete failure-empty TODO", newTodo: "", expHttpCode: 400},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tctrl, _ := NewTodoController(utils.GetSession(utils.DBurl))
			router := httprouter.New()
			router.POST("/todo", tctrl.CreateTodo)

			data := map[string]string{"todo": tc.newTodo}
			dataj, _ := json.Marshal(data)
			todourl := "http://" + SRVADDR + "/todo"
			req, err := http.NewRequest("POST", todourl, bytes.NewBuffer(dataj))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create POST request: %v", err)
			}

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			tododb := models.Todo{}
			json.NewDecoder(rec.Body).Decode(&tododb)
			t.Log(tododb.Id)
			todoId := bson.ObjectId(tododb.Id).Hex()
			router.DELETE("/todo/:id", tctrl.DeleteTodo)
			todourl = "http://" + SRVADDR + "/todo/:id"
			req, err = http.NewRequest("DELETE", todourl, strings.NewReader(todoId))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create DELETE request: %v", err)
			}

			rec = httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if status := rec.Code; status != tc.expHttpCode {
				t.Fatalf("Wrong request status, expected %v:got %v", tc.expHttpCode, status)
			}
		})
	}
}
