package utils

import (
	"gopkg.in/mgo.v2"
	"log"
)

const (
	DBNAME = "rest_todo"
	COLLEC = "todos"
)

var DBurl = "mongodb://localhost"

func GetSession(url string) *mgo.Session {
	s, err := mgo.Dial(url)
	if err != nil {
		log.Println("Could not connect to mongo: ", err.Error())
		return nil
	}
	return s
}
