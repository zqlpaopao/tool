package mysql

import (
	"errors"
	"github.com/zqlpaopao/tool/stringHelper/pkg"
	"strings"
)

// SQLHelperToJointWhere
//  简单拼接where条件助手
//  map的key为数据库字段的名字,驼峰命名则直接转换为蛇形
//  目前支持拼接:("=","in","not in","between and","not null","is null","like","gt","gte","lt","lte")
//
//  whereParamsDemo := map[string]interface{}{
// 		"id":         1,
// 		"name_str":   "str",
// 		"state_in_1": []string{"1", "2", "3"},
// 		"state_in_2": map[string]interface{}{
// 			"type":  "in",
// 			"value": []string{"1", "2", "3"},
// 		},
// 		"state_max": map[string]interface{}{
// 			"type":  "gt",
// 			"value": 1,
// 		},
//       "state_gte": map[string]interface{}{
// 			"type":  "gte",
// 			"value": 1,
// 		},
// 		"state_mix": map[string]interface{}{
// 			"type":  "lt",
// 			"value": 10,
// 		},
// 		"state_lte": map[string]interface{}{
// 			"type":  "lte",
// 			"value": 10,
// 		},
// 		"state_not_in": map[string]interface{}{
// 			"type":  "not in",
// 			"value": []string{"n-1", "n-2", "n-3"},
// 		},
// 		"created_at": map[string]interface{}{
// 			"type":  "between and",
// 			"value": map[string]string{"begin": "2019-10-10", "end": "2020-10-10"},
// 		},
// 		"not_null_column": map[string]interface{}{
// 			"type": "not null",
// 		},
// 		"is_null_column": map[string]interface{}{
// 			"type": "is null",
// 		},
// 		"name_like": map[string]interface{}{
// 			"type":  "like",
// 			"value": "%张三%", //或者"%张三"||"张三%" ,
// 		},
// 	}
//
func SQLHelperToJointWhere(whereParams map[string]interface{}) (whereSql string, whereCase []interface{}, err error) {
	var whereSqlSlice []string
	if whereParams == nil {
		return
	}
	//捕获异常panic
	defer func() {
		if errPanic := recover(); errPanic != nil {
			err = errors.New(errPanic.(string))
		}
	}()

	//遍历whereParams
	for key, value := range whereParams {
		//key驼峰转蛇形
		key = pkg.SnakeString(key)

		switch value1 := value.(type) {
		case string:
			whereSqlSlice = append(whereSqlSlice, key+" = ?")
			whereCase = append(whereCase, value1)
		case int64:
			whereSqlSlice = append(whereSqlSlice, key+" = ?")
			whereCase = append(whereCase, value1)
		case int:
			whereSqlSlice = append(whereSqlSlice, key+" = ?")
			whereCase = append(whereCase, value1)
		case []string:
			var inSlice []string
			var inCase []interface{}
			for _, caseOne := range value1 {
				inSlice = append(inSlice, "?")
				inCase = append(inCase, caseOne)
			}
			whereSqlSlice = append(whereSqlSlice, key+" in ("+strings.Join(inSlice, ",")+")")
			whereCase = append(whereCase, inCase...)
		case []int:
			var inSlice []string
			var inCase []interface{}
			for _, caseOne := range value1 {
				inSlice = append(inSlice, "?")
				inCase = append(inCase, caseOne)
			}
			whereSqlSlice = append(whereSqlSlice, key+" in ("+strings.Join(inSlice, ",")+")")
			whereCase = append(whereCase, inCase...)
		case []int64:
			var inSlice []string
			var inCase []interface{}
			for _, caseOne := range value1 {
				inSlice = append(inSlice, "?")
				inCase = append(inCase, caseOne)
			}
			whereSqlSlice = append(whereSqlSlice, key+" in ("+strings.Join(inSlice, ",")+")")
			whereCase = append(whereCase, inCase...)
		case map[string]interface{}:
			var iType string
			if mapTypeValue, exist := value1["type"]; exist {
				if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "BETWEEN AND" {
					iType = "BETWEEN AND"
					whereSqlSlice = append(whereSqlSlice, key+" BETWEEN ? AND ?")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "NOT IN" {
					iType = "NOT IN"
					var inSlice []string
					switch value1["value"].(type) {
					case []int64:
						for range value1["value"].([]int64) {
							inSlice = append(inSlice, "?")
						}
					case []string:
						for range value1["value"].([]string) {
							inSlice = append(inSlice, "?")
						}
					}
					whereSqlSlice = append(whereSqlSlice, key+" not in ("+strings.Join(inSlice, ",")+")")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "IN" {
					iType = "IN"
					var inSlice []string
					switch value1["value"].(type) {
					case []int64:
						for range value1["value"].([]int64) {
							inSlice = append(inSlice, "?")
						}
					case []string:
						for range value1["value"].([]string) {
							inSlice = append(inSlice, "?")
						}
					}
					whereSqlSlice = append(whereSqlSlice, key+" in ("+strings.Join(inSlice, ",")+")")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "GT" {
					iType = "GT"
					whereSqlSlice = append(whereSqlSlice, key+" > ?")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "GTE" {
					iType = "GTE"
					whereSqlSlice = append(whereSqlSlice, key+" >= ?")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "LTE" {
					iType = "LTE"
					whereSqlSlice = append(whereSqlSlice, key+" <= ?")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "LT" {
					iType = "LT"
					whereSqlSlice = append(whereSqlSlice, key+" < ?")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "NOT NULL" {
					iType = "NOT NULL"
					whereSqlSlice = append(whereSqlSlice, key+" IS NOT NULL")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "IS NULL" {
					iType = "IS NULL"
					whereSqlSlice = append(whereSqlSlice, key+" IS NULL")
				} else if strings.ToUpper(strings.Trim(mapTypeValue.(string), "")) == "LIKE" {
					iType = "LIKE"
					whereSqlSlice = append(whereSqlSlice, key+" LIKE ?")
				} else {
					whereSqlSlice = append(whereSqlSlice, key+" "+strings.ToUpper(strings.Trim(mapTypeValue.(string), ""))+" ?")
				}
			}
			if mapValue, exist := value1["value"]; exist {
				switch mapValue.(type) {
				case string:
					whereCase = append(whereCase, mapValue.(string))
				case int64:
					whereCase = append(whereCase, mapValue.(int64))
				case int:
					whereCase = append(whereCase, mapValue.(int))
				case []int64:
					if len(iType) > 0 {
						switch iType {
						case "NOT IN":
							var stateCase []interface{}
							for _, stateOne := range mapValue.([]int64) {
								stateCase = append(stateCase, pkg.StringFromAssertionFloat(stateOne))
							}
							whereCase = append(whereCase, stateCase...)
						case "IN":
							var stateCase []interface{}
							for _, stateOne := range mapValue.([]int64) {
								stateCase = append(stateCase, pkg.StringFromAssertionFloat(stateOne))
							}
							whereCase = append(whereCase, stateCase...)
						default:
						}
					}
				case []string:
					if len(iType) > 0 {
						switch iType {
						case "NOT IN":
							var stateCase []interface{}
							for _, stateOne := range mapValue.([]string) {
								stateCase = append(stateCase, pkg.StringFromAssertionFloat(stateOne))
							}
							whereCase = append(whereCase, stateCase...)
						case "IN":
							var stateCase []interface{}
							for _, stateOne := range mapValue.([]string) {
								stateCase = append(stateCase, pkg.StringFromAssertionFloat(stateOne))
							}
							whereCase = append(whereCase, stateCase...)
						default:
						}
					}
				case map[string]string:
					if iType == "BETWEEN AND" {
						if begin, exist := mapValue.(map[string]string)["begin"]; exist {
							whereCase = append(whereCase, begin)
						}
						if end, exist := mapValue.(map[string]string)["end"]; exist {
							whereCase = append(whereCase, end)
						}
					} else {
						//todo
					}
				case map[string]int64:
					if iType == "BETWEEN AND" {
						if begin, exist := mapValue.(map[string]int64)["begin"]; exist {
							whereCase = append(whereCase, begin)
						}
						if end, exist := mapValue.(map[string]int64)["end"]; exist {
							whereCase = append(whereCase, end)
						}
					} else {
						//todo
					}
				default:
					//todo
				}
			}
		default:
		}
	}
	whereSql = strings.Join(whereSqlSlice, " AND ")
	return
}
