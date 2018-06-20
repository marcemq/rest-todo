package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/marcemq/rest-todo/models"
	"github.com/marcemq/rest-todo/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const SRVADDR = "0.0.0.0:8080"

type TodoRestApi interface {
	GetTodo(http.ResponseWriter, *http.Request, httprouter.Params)
	GetTodoList(http.ResponseWriter, *http.Request, httprouter.Params)
	CreateTodo(http.ResponseWriter, *http.Request, httprouter.Params)
	DeleteTodo(http.ResponseWriter, *http.Request, httprouter.Params)
}

type TodoController struct {
	session *mgo.Session
}

func NewTodoController(s *mgo.Session) (*TodoController, error) {
	if s == nil {
		return nil, errors.New(fmt.Sprintf("Can't work w/o valid mgo db session!"))
	}
	return &TodoController{s}, nil
}

func fmtOutput(w http.ResponseWriter, code int, todo ...models.Todo) {
	todoj, _ := json.Marshal(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s\n", todoj)
}

func (tc TodoController) GetTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(id)
	todo := models.Todo{}
	if err := tc.session.DB(utils.DBNAME).C(utils.COLLEC).FindId(oid).One(&todo); err != nil {
		w.WriteHeader(404)
		return
	}
	fmtOutput(w, 200, todo)
}
func (tc TodoController) GetTodoList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	todo := []models.Todo{}
	if err := tc.session.DB(utils.DBNAME).C(utils.COLLEC).Find(nil).All(&todo); err != nil {
		w.WriteHeader(404)
		return
	}
	fmtOutput(w, 200, todo...)
}

func (tc TodoController) CreateTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	todo := models.Todo{}
	json.NewDecoder(r.Body).Decode(&todo)
	todo.Id = bson.NewObjectId()
	if todo.Todo == "" {
		w.WriteHeader(400)
		return
	}
	tc.session.DB(utils.DBNAME).C(utils.COLLEC).Insert(todo)
	fmtOutput(w, 200, todo)
}

func (tc TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(id)
	if err := tc.session.DB(utils.DBNAME).C(utils.COLLEC).RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
}
