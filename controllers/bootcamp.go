package controllers

import (
	"devcamper/models"
	"devcamper/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type Bootcamp struct {
	connection *mongodm.Connection
}

func NewBootcamp(conn *mongodm.Connection) *Bootcamp {
	return &Bootcamp{conn}
}

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps
// @access  Public
func (bc *Bootcamp) GetBootcamps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// advance result (sort, selct, limit , etc.)
	// type query struct {
	// 	Select string
	// 	Sort   string
	// 	Page   int
	// 	Limit  int
	// }
	err := r.ParseForm()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad request data"))
		return
	}
	tmp := map[string]interface{}{}
	for k, v := range r.Form {
		tmp[k] = v
	}
	q := extractData(tmp)
	q = cleanData(q)
	// q["deleted"] = false
	utils.SendJSON(w, http.StatusOK, q)
	return
	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamps := []*models.Bootcamp{}

	err = Bootcamp.Find(q).Exec(&bootcamps)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"q":       q,
		"count":   len(bootcamps),
		"data":    bootcamps,
	})
}

// @desc    Get single bootcamp
// @route   GET /api/v1/bootcamps/:id
// @access  Public
func (bc *Bootcamp) GetBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}
	err := Bootcamp.FindId(bson.ObjectIdHex(id)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found bootcamp with id of %s", id))
		return
	} else if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("server error"))
		return
	} else if bootcamp.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this bootcamp was deleted"))
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bootcamp,
	})
}

// @desc    Create bootcamp
// @route   POST /api/v1/bootcamps
// @access  Private
func (bc *Bootcamp) CreateBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	Bootcamp.New(bootcamp)
	err := json.NewDecoder(r.Body).Decode(bootcamp)
	if err != nil {
		log.Println("bad data")
		utils.SendJSON(w, http.StatusBadRequest, errors.New("bad data"))
		return
	}
	if valid, issues := bootcamp.Validate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}
	err = bootcamp.Save()
	if v, ok := err.(*mongodm.ValidationError); ok {
		log.Println(err)
		utils.SendJSON(w, http.StatusBadRequest, v)
		return
	} else if err != nil {
		log.Println(err)
		utils.SendJSON(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    bootcamp,
	})

}

// @desc    Update bootcamp
// @route   PUT /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) UpdateBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	err := Bootcamp.FindId(bson.ObjectIdHex(id)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no bootcamp with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	if bootcamp.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this bootcamp was deleted"))
		return
	}

	var d map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data"))
	}

	// The Update method is incompleted so the error is not handled
	// see https://github.com/zebresel-com/mongodm/issues/20
	bootcamp.Update(d)

	err = bootcamp.Save()
	if _, ok := err.(*mongodm.ValidationError); ok {
		// the updated data not comply with the model requirement
		// grab all error
		_, issues := bootcamp.Validate()
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	} else if err != nil {
		log.Println(err)
		utils.SendJSON(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bootcamp,
	})

}

// @desc    Delete bootcamp
// @route   DELETE /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) DeleteBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	err := Bootcamp.FindId(bson.ObjectIdHex(id)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no bootcamp with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	if bootcamp.Deleted {
		// should not found the deleted bootcamp
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no bootcamp with id of %s", id))
		return
	}

	bootcamp.SetDeleted(true)
	bootcamp.Save()
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    nil,
	})
}

func extractData(data map[string]interface{}) map[string]interface{} {
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
			val = extractData(nested)

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

func cleanData(m map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	var key string
	var val interface{}
	for k, v := range m {
		key = k

		switch refVal := reflect.ValueOf(v); refVal.Kind() {
		case reflect.Slice:

			tmps := make([]interface{}, refVal.Len())
			for i := 0; i < refVal.Len(); i++ {
				tmps[i] = stringToType(refVal.Index(i).String())
			}
			val = tmps

			// remove [] for only one element
			if refVal := reflect.ValueOf(val); refVal.Len() == 1 {
				val = refVal.Index(0).Interface()
			}

		case reflect.Map:
			val = cleanData(v.(map[string]interface{}))
		}
		result[key] = val
	}

	return result
}

func stringToType(s string) interface{} {
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
