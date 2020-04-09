package sqlBuild

import (
	"fmt"
	"reflect"
	"strings"
	"alleyFunAdmin/component/sql-build/debug"
)

func CheckInjection(val string) (injection bool) {
	if val != "" {
		val = strings.ToLower(val)
		injection = strings.Contains(val, "select ") ||
			strings.Contains(val, "update ") ||
			strings.Contains(val, "delete ") ||
			strings.Contains(val, "insert ") ||
			strings.Contains(val, "declare ") ||
			strings.Contains(val, "drop ") ||
			strings.Contains(val, "create ") ||
			strings.Contains(val, "alter ")
		if injection {
			debug.Error("Injection <" + val + ">")
		}
	}
	return
}
func GetInValues(inValues interface{}) (strs []string, err error) {
	switch values := inValues.(type) {
	case []int:
		for _, v := range []int(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []int8:
		for _, v := range []int8(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []int16:
		for _, v := range []int16(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []int32:
		for _, v := range []int32(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []int64:
		for _, v := range []int64(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []uint:
		for _, v := range []uint(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []uint8:
		for _, v := range []uint8(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []uint16:
		for _, v := range []uint16(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []uint32:
		for _, v := range []uint32(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []uint64:
		for _, v := range []uint64(values) {
			str := fmt.Sprintf("%d", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []float64:
		for _, v := range []float64(values) {
			str := fmt.Sprintf("%g", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []float32:
		for _, v := range []float32(values) {
			str := fmt.Sprintf("%g", v)
			if str != "" {
				strs = append(strs, str)
			}
		}
	case []string:
		for _, v := range []string(values) {
			if v != "" {
				if CheckInjection(v) {
					err = ErrInjection
					return
				}
				strs = append(strs, strings.Join([]string{"'", "'"}, v))
			}
		}
	default:
		err = ErrValueType
	}
	return
}

func GetWhereSetFuncValues(values interface{}, rule Rule) (value string,
	err error) {
	return getWhereSetValues(values, rule, func(value string) string {
		return string(value)
	})
}

func GetWhereSetValues(values interface{}, rule Rule) (value string,
	err error) {
	return getWhereSetValues(values, rule, func(value string) string {
		return strings.Join([]string{"'", "'"}, string(value))
	})
}

func getWhereSetValues(values interface{}, rule Rule, f func(value string) string) (value string,
	err error) {
	switch value := values.(type) {
	case int:
		if int(value) > rule.IntValue {
			return fmt.Sprintf("%d", int(value)), nil
		}
	case int8:
		if int8(value) > rule.Int8Value {
			return fmt.Sprintf("%d", int8(value)), nil
		}
	case int16:
		if int16(value) > rule.Int16Value {
			return fmt.Sprintf("%d", int16(value)), nil
		}
	case int32:
		if int32(value) > rule.Int32Value {
			return fmt.Sprintf("%d", int32(value)), nil
		}
	case int64:
		if int64(value) > rule.Int64Value {
			return fmt.Sprintf("%d", int64(value)), nil
		}
	case uint:
		if uint(value) > rule.UintValue {
			return fmt.Sprintf("%d", uint(value)), nil
		}
	case uint8:
		if uint8(value) > rule.Uint8Value {
			return fmt.Sprintf("%d", uint8(value)), nil
		}
	case uint16:
		if uint16(value) > rule.Uint16Value {
			return fmt.Sprintf("%d", uint16(value)), nil
		}
	case uint32:
		if uint32(value) > rule.Uint32Value {
			return fmt.Sprintf("%d", uint32(value)), nil
		}
	case uint64:
		if uint64(value) > rule.Uint64Value {
			return fmt.Sprintf("%d", uint64(value)), nil
		}
	case float64:
		if float64(value) > rule.Float64Value {
			return fmt.Sprintf("%f", float64(value)), nil
		}
	case float32:
		if float32(value) > rule.Float32Value {
			return fmt.Sprintf("%f", float32(value)), nil
		}
	case string:
		if string(value) != rule.StringValue {
			if CheckInjection(string(value)) {
				return "", ErrInjection
			}
			return f(string(value)), nil
		}
	default:
		err = ErrValueType
	}
	return
}

//得到数据类型
func GetValue(value reflect.Value, rule Rule) (string, error) {
	switch value.Kind() {
	case reflect.Int:
		temp := value.Int()
		if int(temp) > rule.IntValue {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Int8:
		temp := value.Int()
		if int8(temp) > rule.Int8Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Int16:
		temp := value.Int()
		if int16(temp) > rule.Int16Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Int32:
		temp := value.Int()
		if int32(temp) > rule.Int32Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Int64:
		temp := value.Int()
		if temp > rule.Int64Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Uint:
		temp := value.Uint()
		if uint(temp) > rule.UintValue {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Uint8:
		temp := value.Uint()
		if uint8(temp) > rule.Uint8Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Uint16:
		temp := value.Uint()
		if uint16(temp) > rule.Uint16Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Uint32:
		temp := value.Uint()
		if uint32(temp) > rule.Uint32Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Uint64:
		temp := value.Uint()
		if uint64(temp) > rule.Uint64Value {
			return fmt.Sprintf("%d", temp), nil
		}
	case reflect.Float32:
		temp := value.Float()
		if float32(temp) > rule.Float32Value {
			return fmt.Sprintf("%f", temp), nil
		}
	case reflect.Float64:
		temp := value.Float()
		if float64(temp) > rule.Float64Value {
			return fmt.Sprintf("%f", temp), nil
		}
	case reflect.String:
		temp := value.String()
		if temp != rule.StringValue {
			if CheckInjection(temp) {
				return "", ErrInjection
			}
			return strings.Join([]string{"'", "'"}, temp), nil
		}
	}
	return "DEFAULT", nil
}

func MapToString(mapdata map[string]interface{}) string{
	reStr:=""
	for k, v := range mapdata {
		kArr:=strings.Split(k, "__")
		kArr[0]=strings.ToUpper(kArr[0])
		if kArr[0] == ""{
			kArr[0] = "AND"
		}
		if kArr[2] == ""{
			kArr[2] = "="
		}
		if kArr[2] == "NULL"{
			kArr[2] = ""
		}
		reStr=reStr+kArr[0]+" "+kArr[1]+""+kArr[2]+v.(string)+" "
	}
	reStr=strings.Trim(reStr, " ")
	reStr=strings.Trim(reStr, "AND")
	reStr=strings.Trim(reStr, "OR")
	reStr=strings.Trim(reStr, " ")
	return reStr
}

func OrderByArrString(sort string) string  {
	if sort=="" {
		return  ""
	}
	sortRank:="DESC"
	sortFirst:=sort[0:1]
	sortString:=""
	if sortFirst=="-" {
		sortRank="DESC"
	}else {
		if sortFirst=="+" {
			sortRank="ASC"
		}
	}
	sortString=sort[1:len(sort)]+" "+sortRank
	return  sortString
}

func DisposeOffset(page,limit int) int {
	if page==0 {
		page=1
	}
	offset:=(page-1)*limit
	return offset
}
