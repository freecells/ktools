/*
 * @Author: Keven
 * @version: v1.0.1
 * @Date: 2021-09-28 13:18:10
 * @LastEditors: Keven
 * @LastEditTime: 2021-09-28 13:18:10
 */
package datadeal

import (
	"reflect"
	"strconv"
	"strings"
)

//Struct2Map 转换
func Struct2Map(obj interface{}) map[string]interface{} {

	t := reflect.TypeOf(obj)

	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {

		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

//Interface2String 接口类型转字符串
func Interface2String(inter interface{}) (res string) {

	switch inter.(type) {

	case string:
		res = inter.(string)
		break

	case int:
		res = strconv.Itoa(inter.(int))
		break

	default:
		res = ""
		break
	}

	return
}

//ReverseString 字符串反序
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)

}
