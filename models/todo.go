package models

import "gopkg.in/mgo.v2/bson"

type (
	Todo struct {
		Id   bson.ObjectId `json:"id" bson:"_id"`
		Todo string        `json:"todo" bson:todo`
	}
)
