package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/marcemq/rest-todo/controllers"
	"github.com/marcemq/rest-todo/utils"
)

func main() {
	r := httprouter.New()
	tc := controllers.NewTodoController(utils.GetSession())
	r.GET("/todo/:id", tc.GetTodo)
	r.GET("/list", tc.GetTodoList)
	r.POST("/todo", tc.CreateTodo)
	r.DELETE("/todo/:id", tc.DeleteTodo)

	err := http.ListenAndServe("0.0.0.0:8080", r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
