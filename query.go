package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

//QueryDeal 处理查询参数 生成查询字串
func QueryDeal(querry string) (whereStr string) {

	querryArr := strings.Split(querry, "&")

	querryMap := make(map[string]string)

	cpMap := make(map[string]string)

	for _, item := range querryArr {

		keyVal := strings.Split(item, "=")

		if keyVal[1] != "" {

			if strings.HasSuffix(keyVal[0], "_cp") {

				cpMap[keyVal[0]] = keyVal[1]

			} else {
				if keyVal[0] != "page" {

					querryMap[keyVal[0]] = "\"" + keyVal[1] + "\""
				}
			}
		}

	}

	cp := map[string]string{
		"gt": ">",
		"ge": ">=",
		"lt": "<",
		"le": "<=",
		"eq": "=",
	}

	querryCount := len(querryMap)

	qand := " and "

	for key, item := range querryMap {

		if querryCount <= 1 {

			qand = ""
		}

		querryCount--

		whereStr += key
		item, _ := url.QueryUnescape(item)
		if cpMap[key+"_cp"] != "" {

			whereStr += cp[cpMap[key+"_cp"]] + item + qand
		} else {

			whereStr += "=" + item + qand
		}
	}

	return
}

//多表关联 where条件生成
func RelationQuery(c *gin.Context, mainTable string) (querryStr string) {

	querry := c.Request.URL.RawQuery

	querryArr := strings.Split(querry, "&")

	querryMap := make(map[string]string)

	cpMap := make(map[string]string)

	for _, item := range querryArr {

		keyVal := strings.Split(item, "=")

		if keyVal[1] != "" {

			if strings.HasSuffix(keyVal[0], "_cp") {

				cpMap[keyVal[0]] = keyVal[1]

			} else {
				if keyVal[0] != "page" {

					querryMap[keyVal[0]] = "\"" + keyVal[1] + "\""
				}
			}
		}

	}

	// fmt.Println(querryMap)
	if len(querryMap) == 0 {
		return
	}

	cp := map[string]string{
		"eq": " = ",
		"gt": " > ",
		"ge": " >= ",
		"lt": " < ",
		"le": " <= ",
	}

	qand := " and "

	whereStrs := make(map[string]string)

	for key, item := range querryMap {

		tableKey := strings.Split(key, ".")

		if len(tableKey) == 2 {

			tableName := tableKey[0]
			tableColumn := tableKey[1]

			whereStrs[tableName] += tableColumn

			if cpMap[key+"_cp"] != "" {

				whereStrs[tableName] += cp[cpMap[key+"_cp"]] + item + qand

			} else {

				whereStrs[tableName] += "=" + item + qand
			}

		} else {

			whereStrs["main"] += key

			if cpMap[key+"_cp"] != "" {

				whereStrs["main"] += cp[cpMap[key+"_cp"]] + item + qand

			} else {

				whereStrs["main"] += "=" + item + qand
			}

		}

	}

	for i, val := range whereStrs {

		whereStrs[i] = strings.TrimSuffix(val, qand)
	}

	querryStr += whereStrs["main"]

	for key, val := range whereStrs {

		ex := " and EXISTS( SELECT * FROM %ss WHERE %s.%s_id=%ss.id AND %s)"

		if key != "main" {

			str := fmt.Sprintf(ex, key, mainTable, key, key, val)

			querryStr += str
		}
	}

	querryStr = strings.TrimPrefix(querryStr, " and ")

	return
}

//QueryDeal2 mssql where querry generate
func QueryDeal2(querryMap url.Values) (whereStr string) {

	fmt.Printf("%+v \n", querryMap)

	cp := map[string]string{
		"eq": "=",
		"gt": ">",
		"ge": ">=",
		"lt": "<",
		"le": "<=",
	}

	whereArr := []string{}

	for key, item := range querryMap {

		//如果是比较条件 或 无数据填充 则返回
		if strings.HasSuffix(key, "_cp") || item[0] == "" {
			// fmt.Printf("item length %d\n", len(item))
			continue
		}
		cpStr := "="

		if _, ok := querryMap[key+"_cp"]; ok {

			val := querryMap[key+"_cp"][0]
			cpStr = cp[val]

		}

		if strings.HasSuffix(key, ":time") {

			key = strings.TrimSuffix(key, ":time")

			dateFmt := "datediff(day,cast('%s' as datetime),%s) %s 0"
			ws := fmt.Sprintf(dateFmt, item[0], key, cpStr)
			whereArr = append(whereArr, ws)

		} else {

			whereArr = append(whereArr, key+cpStr+"'"+item[0]+"'")
		}

	}

	whereStr = strings.Join(whereArr, " and ")

	return
}

func QuerryMulity(c *gin.Context, mainTable string) (querryStr string) {

	querryMap := c.Request.URL.Query()

	delete(querryMap, "page")

	if len(querryMap) == 0 {
		return
	}

	cp := map[string]string{
		"eq": " = ",
		"gt": " > ",
		"ge": " >= ",
		"lt": " < ",
		"le": " <= ",
	}

	qand := " and "

	whereStrs := make(map[string]string)

	for key, item := range querryMap {

		//跳过字段
		if strings.HasSuffix(key, "_cp") || item[0] == "" {
			continue
		}

		tableKey := strings.Split(key, ".")

		if len(tableKey) == 2 {

			tableName := tableKey[0]
			tableColumn := tableKey[1]

			//循环 拉出同一字段的多个条件
			for i, val := range item {

				valCP := querryMap[key+"_cp"]

				if len(valCP) > 0 && val != "" {

					whereStrs[tableName] += tableColumn + cp[valCP[i]] + "\"" + val + "\"" + qand

				} else if val != "" {

					whereStrs[tableName] += tableColumn + "=" + "\"" + val + "\"" + qand
				}

			}

		} else {

			for i, val := range item {

				valCP := querryMap[key+"_cp"]

				if len(valCP) > 0 && val != "" {

					whereStrs["main"] += key + cp[valCP[i]] + "\"" + val + "\"" + qand

				} else if val != "" {

					whereStrs["main"] += key + "=" + "\"" + val + "\"" + qand
				}

			}

		}

	}

	for i, val := range whereStrs {

		whereStrs[i] = strings.TrimSuffix(val, qand)
	}

	querryStr += whereStrs["main"]

	for key, val := range whereStrs {

		ex := " and EXISTS( SELECT * FROM %ss WHERE %s.%s_id=%ss.id AND %s)"

		if key != "main" {

			str := fmt.Sprintf(ex, key, mainTable, key, key, val)

			querryStr += str
		}
	}

	querryStr = strings.TrimPrefix(querryStr, " and ")

	return
}
