package utils

import (
	"devcamper/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

func AdvanceQuery(urlQuery map[string][]string, Model *mongodm.Model) (map[string]interface{}, int, error) {
	// extract data from url query
	rawQuery := ExtractData(ConvQuery(urlQuery))
	rawQuery = CleanData(rawQuery)
	bs, err := json.Marshal(rawQuery)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("bad request data")
	}
	// load query to struct
	queryOption := &queryOption{}
	err = json.Unmarshal(bs, queryOption)
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("server error")
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

	// // this not work
	// models := []*interface{}{}

	// // this not work
	// models := []*mongodm.IDocumentBase{}

	// this work
	models := []*models.Bootcamp{}

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

	var pagination Pagination
	pagination.Fill(page, limit, startIndex, endIndex, total)

	q.Skip(startIndex).Limit(limit)

	// execute query
	err = q.Exec(&models)
	if err != nil {
		log.Println(err)
		return nil, http.StatusInternalServerError, errors.New("server error")
	}

	//  response data
	respData := map[string]interface{}{
		"success":    true,
		"count":      len(models),
		"pagination": pagination,
	}

	// remove other field that not selected
	if queryOption.Select != "" {
		// list all select field in slice
		selects := strings.Split(queryOption.Select, ",")
		// access value of struct field name using reflect
		refVal := reflect.ValueOf(models)
		showFieldModels := make([]map[string]interface{}, len(models))
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
		respData["data"] = showFieldModels
	} else {
		respData["data"] = models
	}

	return respData, http.StatusOK, nil
}
