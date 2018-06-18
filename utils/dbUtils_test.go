package utils

import (
	"reflect"
	"testing"
)

func TestGetSession(t *testing.T) {
	tt := []struct {
		name    string
		url     string
		expType interface{}
		err     error
	}{
		{name: "Successful session", url: DBurl},
		//{name: "Invalid url", url: "MyInvalidUrl", expType: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := GetSession(tc.url)
			t.Log(reflect.TypeOf(s))
			defer s.Close()
		})
	}
}
