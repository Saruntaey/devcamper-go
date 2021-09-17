package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func ConvQuery(q map[string][]string) map[string]interface{} {
	tmp := map[string]interface{}{}
	for k, v := range q {
		tmp[k] = v
	}
	return tmp
}

func ExtractData(data map[string]interface{}) map[string]interface{} {
	dataFormated := map[string]interface{}{}
	operator := []string{"gt", "gte", "lt", "lte", "in"}
	for k, v := range data {
		// extract data
		var key string
		var val interface{}
		ks := strings.Split(k, "[")
		if len(ks) == 1 {
			// this mean key is unique and not have nested
			key = ks[0]
			// insert operator
			for _, o := range operator {
				if key == o {
					key = fmt.Sprintf("$%s", key)
				}
			}

			val = v

		} else {
			// this mean key have nested (name[in], price[lt], etc.)
			key = ks[0]
			nested := map[string]interface{}{
				strings.TrimRight(ks[1], "]"): v,
			}
			val = ExtractData(nested)

			// check if key anlready push to dataFormated
			// the query sting look like this (price[gt]=1000&price[lt]=2000)
			if v, ok := dataFormated[key]; ok {
				comb := map[string]interface{}{}
				// copy pushed data
				refVal := reflect.ValueOf(v)
				for _, refValKey := range refVal.MapKeys() {
					comb[refValKey.String()] = refVal.MapIndex(refValKey).Interface()
				}
				// push new data
				refVal = reflect.ValueOf(val)
				for _, refValKey := range refVal.MapKeys() {
					// check if nested key exist
					// the query sting look like this (role[in]=user&role[in]=admin)
					if pushedVal, ok := comb[refValKey.String()]; ok {
						n := reflect.ValueOf(pushedVal).Len()
						valDupKey := make([]interface{}, n+1)
						// copy pushed data
						for i := 0; i < n; i++ {
							valDupKey[i] = reflect.ValueOf(pushedVal).Index(i).String()
						}
						// append new data
						valDupKey = append(valDupKey, refVal.MapIndex(refValKey).String())

						comb[refValKey.String()] = valDupKey

					} else {
						comb[refValKey.String()] = refVal.MapIndex(refValKey).Interface()
					}

				}
				val = comb
			}
		}
		dataFormated[key] = val
	}
	return dataFormated
}

func CleanData(m map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	var key string
	var val interface{}
	for k, v := range m {
		key = k

		switch refVal := reflect.ValueOf(v); refVal.Kind() {
		case reflect.Slice:

			tmps := make([]interface{}, refVal.Len())
			for i := 0; i < refVal.Len(); i++ {
				tmps[i] = StringToType(refVal.Index(i).String())
			}
			val = tmps

			// remove [] for only one element
			if refVal := reflect.ValueOf(val); refVal.Len() == 1 {
				val = refVal.Index(0).Interface()
			}

		case reflect.Map:
			val = CleanData(v.(map[string]interface{}))
		}
		result[key] = val
	}

	return result
}

func StringToType(s string) interface{} {
	var result interface{}
	if d, err := strconv.ParseInt(s, 0, 64); err == nil {
		result = d
	} else if d, err := strconv.ParseFloat(s, 64); err == nil {
		result = d
	} else if d, err := strconv.ParseBool(s); err == nil {
		result = d
	} else {
		result = s
	}

	return result
}
