package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/marcemq/rest-todo/utils"
)

func TestCreateTodo(t *testing.T) {
	tt := []struct {
		name        string
		todo        string
		expHttpCode int
	}{
		{name: "Secret TODO", todo: "My secret TODO", expHttpCode: 200},
		{name: "Empty TODO", todo: "", expHttpCode: 400},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tctrl := NewTodoController(utils.GetSession())
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
