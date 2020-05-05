package DBCenter

import (
	"reflect"
	"unsafe"

	Utils "../Utils"
)

//根据条件获取主键列表
func DbgetIdentify(model interface{}, vague int) []string {
	var id string
	var cnt int
	tp := reflect.TypeOf(model).Elem()
	value := reflect.ValueOf(model).Elem()
	sql := "select "
	for i := 0; i < value.NumField(); i++ {
		v := Utils.ManyTypeToString(value.Field(i))
		if v == "identify" {
			value.Field(i).SetString("")
			cnt = DbgetCount(model)
			if vague == 2 {
				cnt = DbgetVagueCount(model)
			}
			id = tp.Field(i).Name
			sql += id + " from " + tp.Name() + " where "
		} else if !reflect.Value.IsZero(value.Field(i)) {
			if vague > 0 {
				sql += tp.Field(i).Name + " like '%" + Utils.ManyTypeToString((value.Field(i))) + "%'"
			} else {
				sql += tp.Field(i).Name + " = '" + Utils.ManyTypeToString((value.Field(i))) + "'"
			}
			break
		}

	}
	//fmt.Println(sql, cnt)
	var store = make([]interface{}, cnt)
	var ret = make([]string, 0)
	DbgetWithSql(sql, cnt, 0, model, store)
	for i := 0; i < len(store); i++ {
		v := reflect.ValueOf(store[i])
		ret = append(ret, v.FieldByName(id).String())
	}
	return ret
}
func DbgetSumWithModel(model interface{}, fieldName string, groupByName string) int {
	tp := reflect.TypeOf(model)
	value := reflect.ValueOf(model)
	sql := "select sum(" + fieldName + ") from " + tp.Name()
	for i := 0; i < value.NumField(); i++ {
		if !reflect.Value.IsZero(value.Field(i)) {
			sql += " where " + tp.Field(i).Name + "='" +
				Utils.ManyTypeToString((value.Field(i))) + "'"
			break
		}
	}
	sql += " group by " + groupByName
	//fmt.Println(sql)
	sum := DbgetCountWithSql(sql)
	return sum

}
func DbgetCountWithSql(sql string) (cnt int) {
	row, err := db.Query(sql)
	if err != nil {
		return 0
	}
	row.Next()
	er := row.Scan(&cnt)
	if er != nil {
		cnt = 0
	}
	return cnt
}
func DbgetVagueCount(model interface{}) (cnt int) {
	tType := reflect.TypeOf(model).Elem()
	sql := "select count(*) count from " + tType.Name() + " where "
	tValue := reflect.ValueOf(model).Elem()
	hasCondition := false
	for i := 0; i < tValue.NumField(); i++ {
		if tValue.Field(i).String() != "" {
			sql += tType.Field(i).Name + " like '%" + tValue.Field(i).String()
			sql += "%' and "
			hasCondition = true
		}
	}
	//增加了条件判断去掉最后的and,否则去掉where
	if hasCondition {
		sql = sql[0 : len(sql)-5]
	} else {
		sql = sql[0 : len(sql)-7]
	}
	row, _ := db.Query(sql)
	row.Next()
	er := row.Scan(&cnt)
	if er != nil {
		cnt = 0
	}
	defer row.Close()
	return cnt
}
func DbgetCount(model interface{}) (cnt int) {
	tType := reflect.TypeOf(model).Elem()
	sql := "select count(*) count from " + tType.Name() + " where "
	tValue := reflect.ValueOf(model).Elem()
	hasCondition := false
	for i := 0; i < tValue.NumField(); i++ {
		if tValue.Field(i).String() != "" {
			sql += tType.Field(i).Name + " = '" + tValue.Field(i).String()
			sql += "' and "
			hasCondition = true
		}
	}
	//增加了条件判断去掉最后的and,否则去掉where
	if hasCondition {
		sql = sql[0 : len(sql)-5]
	} else {
		sql = sql[0 : len(sql)-7]
	}
	row, _ := db.Query(sql)
	row.Next()
	er := row.Scan(&cnt)
	if er != nil {
		cnt = 0
	}
	defer row.Close()
	return cnt
}
func DbgetWithSql(sql string, rCount int, rOffset int, model interface{}, store interface{}) {
	if rCount != -1 {
		sql = sql + " limit " + Utils.IntToString(rCount)
		if rOffset > 0 {
			sql += " offset " + Utils.IntToString(rOffset)
		}
	}
	//fmt.Println("serach=====:", sql)
	rows, err := db.Query(sql)
	checkDbError(err)
	rowCount := 0
	for rows.Next() {
		columns, _ := rows.Columns()
		colNum := len(columns)

		vByte := make([][]byte, colNum)
		vArray := make([]interface{}, colNum)
		vData := make([]string, colNum)
		for i := 0; i < colNum; i++ {
			vArray[i] = &vByte[i]
		}
		//扫描一行记录
		err := rows.Scan(vArray...)
		checkDbError(err)
		var cvStr *string
		var cvByte *[]byte
		//*[]byte -> *string
		for i := 0; i < colNum; i++ {
			cvByte = vArray[i].(*[]byte)
			cvStr = (*string)(unsafe.Pointer(cvByte))
			//恢复转换的特殊字符
			*cvStr = Utils.SafeStrRecovery(*cvStr)
			vData[i] = *cvStr
		}
		//key[],value[] -> map
		tmpMap := Utils.MapOfArray(columns, vData)
		Utils.ModelOfMap(model, tmpMap)

		if store != nil {
			//如果传递了一个容器, storeType在为切片或数组时有效
			storeType := reflect.TypeOf(store).Kind()
			if storeType == reflect.Slice ||
				storeType == reflect.Array {
				arr := reflect.ValueOf(store)
				obj := reflect.ValueOf(model).Elem()
				arr.Index(rowCount).Set(obj)
				rowCount++
			}
		}
	}
	defer rows.Close()
}

//通过模型查找数据库,rOffset 指定偏移
func DbgetWithModel(structModel interface{}, store interface{}, rOffset int, vsort string) {
	stName := reflect.TypeOf(structModel).Elem()
	fields := reflect.ValueOf(structModel).Elem()
	var sql string = "select * from " + stName.Name() + " where "
	var isCondition = false
	for i := 0; i < fields.NumField(); i++ {
		kType := stName.Field(i).Type.String()
		kName := stName.Field(i).Name
		vValue := ""
		switch kType {
		case "int":
			intD := int(fields.Field(i).Int())
			vValue = Utils.IntToString(intD)
			fields.Field(i).SetInt(0)
		case "string":
			vValue = fields.Field(i).String()
			fields.Field(i).SetString("")
		}
		if vValue != "" {
			isCondition = true
			sql += Utils.FirstTolower(kName) + " = '" + vValue + "' and "
		}
	}
	/*
		有搜索条件:去掉多余的" and ";
		没有搜索条件:去掉多余的 " where "
	*/
	if isCondition {
		sql = sql[0 : len(sql)-5]
	} else {
		sql = sql[0 : len(sql)-7]
	}
	if vsort != "" {
		sql += " order by " + vsort
	}
	rc := 1

	if store != nil {
		rc = reflect.ValueOf(store).Len()
	}
	DbgetWithSql(sql, rc, rOffset, structModel, store)
	//fmt.Println("struct:", structModel)
}
func DbgetWithOneModel(structModel interface{}) {
	DbgetWithModel(structModel, nil, 0, "")
}

//查询所有,传递一个数组的地址
func DbgetAllModel(structModel interface{}, store interface{}) {
	seType := reflect.TypeOf(store).Kind()
	if seType == reflect.Ptr {
		cnt := DbgetCount(structModel)
		store = Utils.ExpandArray(store, cnt)
	}
	DbgetWithModel(structModel, store, 0, "")
}

func DbgetModelWithSql(model interface{}, store interface{}, sql string) {
	DbgetWithSql(sql, -1, 0, model, store)
}
