package handler

import (
	"fmt"
	"reflect"
	"strings"
	"unicode/utf8"
)

type NullType byte

const (
	_ NullType = iota
	// IsNull the same as `is null`
	IsNull
	// IsNotNull the same as `is not null`
	IsNotNull
)

func WhereBuild(where map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for k, v := range where {
		ks := strings.Split(k, " ")
		if len(ks) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", k)
		}

		if whereSQL != "" {
			whereSQL += " AND "
		}
		strings.Join(ks, ",")
		switch len(ks) {
		case 1:
			//fmt.Println(reflect.TypeOf(v))
			switch v := v.(type) {
			case NullType:
				if v == IsNotNull {
					whereSQL += fmt.Sprint(k, " IS NOT NULL")
				} else {
					whereSQL += fmt.Sprint(k, " IS NULL")
				}
			default:
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
			}
			break
		case 2:
			k = ks[0]
			switch ks[1] {
			case "=":
				whereSQL += fmt.Sprint(k, "=?")
				vals = append(vals, v)
				break
			case ">":
				whereSQL += fmt.Sprint(k, ">?")
				vals = append(vals, v)
				break
			case ">=":
				whereSQL += fmt.Sprint(k, ">=?")
				vals = append(vals, v)
				break
			case "<":
				whereSQL += fmt.Sprint(k, "<?")
				vals = append(vals, v)
				break
			case "<=":
				whereSQL += fmt.Sprint(k, "<=?")
				vals = append(vals, v)
				break
			case "!=":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "<>":
				whereSQL += fmt.Sprint(k, "!=?")
				vals = append(vals, v)
				break
			case "in":
				whereSQL += fmt.Sprint(k, " in (?)")
				vals = append(vals, v)
				break
			case "like":
				whereSQL += fmt.Sprint(k, " like ?")
				vals = append(vals, v)
			}
			break
		}
	}
	return
}

func StructSliceToMap(source interface{}, filedName string) map[string][]interface{} {
	filedIndex := 0
	v := reflect.ValueOf(source) // 判断，interface转为[]interface{}
	resMap := make(map[string][]interface{})
	if v.Kind() != reflect.Slice {
		return resMap
	}
	l := v.Len()
	retList := make([]interface{}, l)
	for i := 0; i < l; i++ {
		retList[i] = v.Index(i).Interface()
	}
	if len(retList) > 0 {
		firstObj := retList[0]
		objT := reflect.TypeOf(firstObj)
		for i := 0; i < objT.NumField(); i++ {
			if objT.Field(i).Name == filedName {
				filedIndex = i
			}
		}
	}
	for _, elem := range retList {
		key := reflect.ValueOf(elem).Field(filedIndex).Interface()
		value := make([]interface{}, 0)
		resMap[fmt.Sprint(key)] = value
	}

	for _, elem := range retList {
		key := reflect.ValueOf(elem).Field(filedIndex).Interface()
		resMap[fmt.Sprint(key)] = append(resMap[fmt.Sprint(key)], elem)
	}
	return resMap
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
