package DBCenter

import (
	"fmt"
	"reflect"

	"../Utils"
	_ "../lib/mysql"
)

func DeleteData(sctObj interface{}) bool {
	table := reflect.TypeOf(sctObj)
	store := Utils.MapOfModel(sctObj)

	sql := "delete from " + table.Name() + " where "
	hasValue := false
	for k, v := range store {
		vs := Utils.AnyTypeToString(v)
		if vs == "" {
			continue
		}
		hasValue = true
		sql += (k + " = '" + vs + "' and ")
	}
	if !hasValue {
		return false
	}
	sql = sql[0 : len(sql)-5]
	result, err := db.Exec(sql)
	row, _ := result.RowsAffected()
	if err != nil {
		fmt.Println("错误======!", row, err, "==========!")
		return false
	}
	fmt.Println(sql)
	return true
}
func UpdateTable(sctObj interface{}) bool {
	table := reflect.TypeOf(sctObj)
	store := Utils.MapOfModel(sctObj)
	sql := "update " + table.Name() + " set "
	/*
		取第一个字段为更新条件好吧,管它是jb啥
	*/
	id := table.Field(0).Name
	var id_value interface{}
	hasValue := false
	for k, v := range store {
		vs := Utils.AnyTypeToString(v)
		if id == k {
			id_value = v
			continue
		}
		if vs == "" {
			continue
		}
		hasValue = true
		sql += (k + " = '" + vs + "',")
	}
	if !hasValue {
		return false
	}
	//去掉逗号
	sql = sql[0 : len(sql)-1]
	sql += " where " + id + " = '" + Utils.AnyTypeToString(id_value) + "'"

	result, err := db.Exec(sql)

	fmt.Println(result, err)
	if err != nil {
		return false
	}
	fmt.Println(sql)
	return true
}
func InsertTable(sctObj interface{}) bool {
	table := reflect.TypeOf(sctObj)
	store := Utils.MapOfModel(sctObj)
	sql := "insert into " + table.Name() + "("
	sdl := ""
	hasValue := false
	for k, v := range store {
		vs := Utils.AnyTypeToString(v)
		if vs != "" {
			hasValue = true
			sql += k + ","
			sdl += "'" + vs + "',"
		}
	}
	if !hasValue {
		return false
	}
	sql = sql[0 : len(sql)-1]
	sdl = sdl[0 : len(sdl)-1]
	sql += ")values(" + sdl + ")"
	_, err := db.Exec(sql)

	if err != nil {
		return false
	}
	fmt.Println(sql)

	return true
}
