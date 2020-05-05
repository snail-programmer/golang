package Utils

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func StringToInt(arg string) int {
	if arg == "" {
		arg = "0"
	}
	v, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	}
	return v
}
func IntToString(arg int) string {
	v := strconv.Itoa(arg)
	return v
}
func Float64ToString(arg float64) string {
	v := fmt.Sprintf("%.2f", arg)
	return v
}
func StringToFloat(arg string) float64 {
	v, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0
	}
	return v
}
func RdcNumString(params ...string) string {
	src := 0.0
	letNum := false
	for k, v := range params {
		if strings.Index(v, ".") > -1 {
			if k == 0 {
				letNum = true
				src = StringToFloat(v)
			} else {
				src -= StringToFloat(v)
			}
		} else {
			if k == 0 {
				src = float64(StringToInt(v))
			} else {
				src -= float64(StringToInt(v))
			}
		}
	}
	res := ""
	if letNum {
		res = fmt.Sprintf("%.2f", src)
	} else {
		res = fmt.Sprintf("%.0f", src)
	}
	return res
}
func AddNumString(params ...string) string {
	sum := 0.0
	letNum := false
	for _, v := range params {
		if strings.Index(v, ".") > -1 {
			letNum = true
			sum += StringToFloat(v)
		} else {
			sum += float64(StringToInt(v))
		}
	}
	res := ""
	if letNum {
		res = fmt.Sprintf("%.2f", sum)
	} else {
		res = fmt.Sprintf("%.0f", sum)
	}
	return res
}
func ManyTypeToString(arg interface{}) (res string) {
	vT := reflect.TypeOf(arg)
	isVt := false
	if vT.Name() == "Value" {
		vT = arg.(reflect.Value).Type()
		isVt = true
	}
	switch vT.Kind() {
	case reflect.Int:
		if isVt {
			res = IntToString(int(arg.(reflect.Value).Int()))
		} else {
			res = IntToString(arg.(int))
		}
		break
	case reflect.String:
		if isVt {
			res = arg.(reflect.Value).String()
		} else {
			res = arg.(string)
		}
		break
	case reflect.Float64:
		if isVt {
			res = Float64ToString(arg.(reflect.Value).Float())
		} else {
			res = Float64ToString(arg.(float64))
		}
		break
	}
	return res
}
func IsfileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
func runError(res *bool) {
	if err := recover(); err != nil {
		fmt.Println("Utils ------- runError", err, "--------------")
		*res = false
	}
}
func FirstToupper(arg string) string {
	/*首字母转大写*/
	if arg == "" {
		return arg
	}
	var newk *string = &arg
	var vkByte []byte = make([]byte, len(arg))
	copy(vkByte, arg)
	if vkByte[0] >= 97 && vkByte[0] <= 122 {
		vkByte[0] -= 32
		newk = (*string)(unsafe.Pointer(&vkByte))
	}
	return *newk
}
func FirstTolower(arg string) string {
	/*k 首字母转小写*/
	if arg == "" {
		return arg
	}
	var newk *string = &arg
	var vkByte []byte = make([]byte, len(arg))
	copy(vkByte, arg)
	if vkByte[0] >= 65 && vkByte[0] <= 90 {
		vkByte[0] += 32
		newk = (*string)(unsafe.Pointer(&vkByte))
	}
	return *newk
}

//传递一个结构体指针,数组转化为model,不支持model内嵌model的递归转化
func ModelOfArray(structName interface{}, data []string) (res bool) {
	res = true
	defer runError(&res)
	obj := reflect.ValueOf(structName)
	objType := reflect.TypeOf(structName)
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
		objType = objType.Elem()
	} else {
		panic("error: trans an address of object")
	}
	if objType.Name() == "Value" {
		obj = obj.Interface().(reflect.Value)
	}
	if obj.CanSet() {
		for i := 0; i < obj.NumField(); i++ {
			if i >= len(data) {
				res = i > 0
				return res //true || false
			}
			/*
			  如果obj为反射后创建的对象，不能再用反射获取成员属性的类型
			*/
			var oT string
			if objType.Name() == "Value" {
				oT = obj.Field(i).Type().Name()
			} else {
				oT = objType.Field(i).Type.String()
			}
			var setD = data[i]
			switch oT {
			case "string":
				obj.Field(i).SetString(setD)
			case "int":
				var intD int = StringToInt(setD)
				obj.Field(i).SetInt(int64(intD))
			}
		}
	} else {
		res = false
		panic("对象不可写")
	}
	return res
}
func ModelOfMap(structName interface{}, data map[string]interface{}) (res bool) {
	res = true
	defer runError(&res)
	obj := reflect.ValueOf(structName)
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	} else {
		panic("error: trans an address of object")
	}
	if obj.Type().Name() == "Value" {
		obj = obj.Interface().(reflect.Value)
	}
	if obj.CanSet() {
		for k, v := range data {
			/*k 首字母转大写*/
			k = FirstToupper(k)
			vType := reflect.TypeOf(v).Kind().String()
			switch vType {
			case "string":
				strD := v.(string)
				obj.FieldByName(k).SetString(strD)
			case "int":
				var intD int = v.(int)
				obj.FieldByName(k).SetInt(int64(intD))
			}
		}
	} else {
		res = false
	}
	return res
}
func MapOfArray(keys []string, value []string) map[string]interface{} {
	ret := map[string]interface{}{}
	if len(keys) == 0 {
		return ret
	}
	for i, key := range keys {
		if key != "" {
			ret[key] = value[i]
		}
	}
	return ret
}
func MapOfModel(sctObj interface{}) map[string]interface{} {
	store := make(map[string]interface{}, 0)
	fieldtype := reflect.TypeOf(sctObj)
	fieldValue := reflect.ValueOf(sctObj)
	if fieldValue.IsValid() {
		for i := 0; i < fieldValue.NumField(); i++ {
			key := fieldtype.Field(i).Name
			value := fieldValue.Field(i)
			store[key] = value.Interface()
		}
	}
	return store
}
func AnyTypeToString(conv interface{}) (res string) {
	vType := reflect.TypeOf(conv).Kind()
	switch vType {
	case reflect.Int:
		res = IntToString(conv.(int))
	case reflect.String:
		res = conv.(string)
	}
	return res
}

//传递一个数组地址,store => &[]interface{}
func ExpandArray(store interface{}, len int) []interface{} {
	var st []interface{}
	if reflect.TypeOf(store).Kind() == reflect.Ptr {
		rp := reflect.ValueOf(store).Elem()
		st = make([]interface{}, len)
		rp.Set(reflect.ValueOf(st))
	}
	return st
}
func SafeStrConvert(str string) string {
	if strings.Index(str, "\\") > -1 {
		str = strings.ReplaceAll(str, "\\", "\\\\")
	}
	if strings.Index(str, "'") > -1 {
		str = strings.ReplaceAll(str, "'", "\\'")
	}
	return str
}
func SafeStrRecovery(str string) string {
	if strings.Index(str, "\\'") > -1 {
		str = strings.ReplaceAll(str, "\\'", "'")
	}
	if strings.Index(str, "\\\\") > -1 {
		str = strings.ReplaceAll(str, "\\\\'", "\\")
	}
	return str
}
