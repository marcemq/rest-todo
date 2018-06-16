package utils

import "gopkg.in/mgo.v2"

const (
	DBNAME = "rest_todo"
	COLLEC = "todos"
)

func GetSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	return s
}
