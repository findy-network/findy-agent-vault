package model

import (
	"encoding/base64"
	"reflect"
	"strconv"
)

func CreateCursor(created int64, object interface{}) string {
	typeName := reflect.TypeOf(object).Name()
	return base64.StdEncoding.EncodeToString([]byte(typeName + ":" + strconv.FormatInt(created, 10)))
}

type APIObject interface {
	Identifier() string
	Created() int64
	Pairwise() *InternalPairwise
	Event() *InternalEvent
	Job() *InternalJob
}
