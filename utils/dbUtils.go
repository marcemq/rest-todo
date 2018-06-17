package utils

import "gopkg.in/mgo.v2"

const (
	DBNAME = "rest_todo"
	COLLEC = "todos"
)

var dburl = "mongodb://localhost"

func GetSession() *mgo.Session {
	s, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	return s
}
