package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/marcemq/rest-todo/models"
	"github.com/marcemq/rest-todo/utils"
)

func TestCreateTodo(t *testing.T) {
	tt := []struct {
		name        string
		todo        string
		expectedMsg string
		err         error
	}{
		{name: "Secret TODO", todo: "My secret TODO", expectedMsg: "My secret TODO"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewTodoController(utils.GetSession())
			router := httprouter.New()
			router.POST("/todo", handler.CreateTodo)

			data := url.Values{}
			data.Set("todo", tc.todo)
			srv := "http://localhost:8080/"
			t.Log(strings.NewReader(data.Encode()))
			req, err := http.NewRequest("POST", srv+"todo", strings.NewReader(data.Encode()))
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Fatalf("Could not create POST request: %v", err)
			}
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req) //We migth need to create a httprouter.params
			if status := rec.Code; status != http.StatusOK {
				t.Fatalf("Wrong status  %v", status)
			}
			resp := rec.Result()
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Could not read response: %v", err)
			}
			//Assert content of inserted TODO
			strbody := string(bytes.TrimSpace(body))
			t.Log(strbody)
			// output [{"id":"5b25fe6aa0d88c61c8dcc533","todo":""}]
			// Todo is empty, dunno why
			out := models.Todo{}
			json.NewDecoder(resp.Body).Decode(&out)
			outj, _ := json.Marshal(out)
			t.Log(out)
			t.Log(string(outj))
			t.Log(out.Todo)
		})
	}
}
