package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/marcemq/rest-todo/controllers"
	"github.com/marcemq/rest-todo/utils"
)

func main() {
	r := httprouter.New()
	s := utils.GetSession(utils.DBurl)
	tc, err := controllers.NewTodoController(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error found: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	r.GET("/todo/:id", tc.GetTodo)
	r.GET("/list", tc.GetTodoList)
	r.POST("/todo", tc.CreateTodo)
	r.DELETE("/todo/:id", tc.DeleteTodo)

	err = http.ListenAndServe(controllers.SRVADDR, r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
