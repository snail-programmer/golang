package main

import (
	"html/template"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"../Utils"
)

//装载post数据->model对象
func LoadModelWithPostData(model interface{}, urlValue url.Values) {
	rType := reflect.TypeOf(model)
	if rType.Kind() != reflect.Ptr {
		panic("give a address of model")
	}
	rType = rType.Elem()
	rValue := reflect.ValueOf(model).Elem()
	for i := 0; i < rValue.NumField(); i++ {
		field := rType.Field(i).Name
		value := rValue.Field(i)
		if value.String() == "" {
			//首字母大写没得到值，转为小写重新请求一下
			if urlValue.Get(field) == "" {
				field = Utils.FirstTolower(field)
			}
			var fieldValue = urlValue.Get(field)
			//特殊字符转换
			fieldValue = Utils.SafeStrConvert(fieldValue)
			value.SetString(fieldValue)
		} else {
			value.SetString("")
		}
	}
}

/*
	查询str中第一个存在字段的索引
*/
func indexWithPropAtStr(model interface{}, str string) (string, int) {
	var min int = -1
	var find string = ""
	rType := reflect.TypeOf(model)
	if rType.Kind() != reflect.Ptr {
		panic("give a address of model")
	}
	rType = rType.Elem()
	rValue := reflect.ValueOf(model).Elem()
	var res int = -1
	for i := 0; i < rValue.NumField(); i++ {
		field := rType.Field(i).Name
		res = strings.Index(str, field+"=")
		if res == -1 {
			//首字母转小写
			field = Utils.FirstTolower(field)
			res = strings.Index(str, field+"=")
		}
		if res > -1 {
			if min == -1 {
				min = res
				find = field
			}
			if res < min {
				min = res
				find = field
			}
		}
	}
	return find, min
}

//text/plain -> model object
func LoadModelWithByte(model interface{}, vByte []byte) {
	var store = make(map[string]interface{})
	vD := string(vByte)
	//特殊字符转换
	vD = Utils.SafeStrConvert(vD)
	for vD != "" {
		/*
			在vByte中寻找第一个字段的索引值
		*/
		fD, oi := indexWithPropAtStr(model, vD)

		//没有字段存在，退出
		if oi < 0 {
			break
		}
		//取第一个字段名之后的记录
		vD = vD[oi+len(fD)+1 : len(vD)]
		//取下一个字段名的索引
		_, ooi := indexWithPropAtStr(model, vD)
		//不存在的话取剩余所有记录，退出
		if ooi < 0 {
			store[fD] = vD
			break
		}
		//第一个字段尾部和第二个字段名开始的区间，是第一个字段的记录值
		store[fD] = vD[0 : ooi-1]
		//取下一个字段名之后的记录
		vD = vD[ooi:len(vD)]
	}
	Utils.ModelOfMap(model, store)

}
func PackMsg(code string, msg string) map[string]interface{} {
	var ret = map[string]interface{}{"code": code, "msg": msg}
	return ret
}
func PackMsgAndSend(code string, msg string, w http.ResponseWriter) {
	w.Write(GenericPackJson(PackMsg(code, msg)))
}
func JumpErrorPage(w http.ResponseWriter, error string, jumpUrl string) {
	e := map[string]string{"error": error, "jumpUrl": jumpUrl}
	t, _ := template.ParseFiles("../view/error.html")
	t.Execute(w, e)
}
