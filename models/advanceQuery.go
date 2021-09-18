package models

import (
	"devcamper/utils"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

// advance result (sort, selct, limit , etc.)
type queryOption struct {
	Select string
	Sort   string
	Page   int
	Limit  int
}

// type ModelQuery interface {
// 	MallocOne()
// 	MallocMany()
// 	GetModel() *mongodm.Model
// }

// func (m *Model) MallocMany() []*Model {
// 	return []*Model{}
// }
// , m ModelQuery
func AdvanceQuery(urlQuery map[string][]string, Model *mongodm.Model) (*mongodm.Query, Pagination, error) {
	// init return data
	var pagination Pagination

	// extract data from url query
	rawQuery := utils.ExtractData(utils.ConvQuery(urlQuery))
	rawQuery = utils.CleanData(rawQuery)
	bs, err := json.Marshal(rawQuery)
	if err != nil {
		return nil, pagination, err
	}
	// load query to struct
	queryOption := &queryOption{}
	err = json.Unmarshal(bs, queryOption)
	if err != nil {
		return nil, pagination, err
	}

	// create query
	query := rawQuery
	// add deleted field
	query["deleted"] = false
	// delete options field
	delete(query, "select")
	delete(query, "sort")
	delete(query, "page")
	delete(query, "limit")

	// init query
	q := Model.Find(query)

	// select fields
	if queryOption.Select != "" {
		selects := strings.Split(queryOption.Select, ",")
		selectQuery := bson.M{}
		for _, v := range selects {
			selectQuery[v] = 1
		}
		q = q.Select(selectQuery)
	}

	// sort
	if queryOption.Sort != "" {
		sorts := strings.Split(queryOption.Sort, ",")
		q = q.Sort(sorts...)
	} else {
		q = q.Sort("-createdAt")
	}

	// pagination
	page := queryOption.Page
	limit := queryOption.Limit
	// set default if not provided
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 100
	}

	startIndex := (page - 1) * limit
	endIndex := page * limit
	total, _ := Model.Find(bson.M{"deleted": false}).Count()

	pagination.Fill(page, limit, startIndex, endIndex, total)

	q.Skip(startIndex).Limit(limit)

	return q, pagination, nil
}

func ExtractSelectField(models interface{}, selects []string) []map[string]interface{} {
	// access value of struct field name using reflect
	refVal := reflect.ValueOf(models)
	showFieldModels := make([]map[string]interface{}, refVal.Len())
	for i := 0; i < refVal.Len(); i++ {
		v := map[string]interface{}{}
		for _, fieldName := range selects {
			nameInStruc := strings.ToUpper(fieldName[:1]) + strings.ToLower(fieldName[1:])
			// check if the struct have the field name
			if val := refVal.Index(i).Elem().FieldByName(nameInStruc); val.IsValid() {
				v[fieldName] = val.Interface()
			}
		}
		showFieldModels[i] = v
	}
	return showFieldModels
}
