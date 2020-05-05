package main

import (
	"fmt"
	"reflect"
	"strings"

	DBModel "../DataBaseCenter/DataBaseModel"
)

var RPCfunc = map[string]interface{}{
	"User":     DBModel.User{},
	"Category": DBModel.Category{},
}

func RPC_getInfo(pattern string) {
	var obj interface{}
	for k, v := range RPCfunc {
		if strings.Contains(pattern, k) {
			obj = v
			break
		}
	}
	v1 := reflect.ValueOf(&obj).Elem().Interface()
	v3 := reflect.TypeOf(v1)
	oc := reflect.New(v3).Elem()
	fmt.Println(v3.NumField(), oc.CanSet())
}
