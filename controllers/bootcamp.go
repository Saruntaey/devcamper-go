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
	q := extractData(r.Form)
	q["deleted"] = false
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

func extractData(data map[string][]string) map[string]interface{} {
	dataFormated := map[string]interface{}{}
	operator := []string{"gt", "gte", "lt", "lte", "in"}
	for k, v := range data {
		// extract data
		var key string
		var val interface{}
		ks := strings.Split(k, "[")
		if len(ks) == 1 {
			key = ks[0]
			// insert operator
			for _, o := range operator {
				if key == o {
					key = fmt.Sprintf("$%s", key)
				}
			}
			// convert data type
			vType := []interface{}{}
			for _, s := range v {
				if d, err := strconv.ParseInt(s, 0, 64); err == nil {
					vType = append(vType, d)
				} else if d, err := strconv.ParseFloat(s, 64); err == nil {
					vType = append(vType, d)
				} else if d, err := strconv.ParseBool(s); err == nil {
					vType = append(vType, d)
				} else {
					vType = append(vType, s)
				}
			}

			if len(vType) == 1 {
				val = vType[0]
			} else {
				val = vType
			}
		} else {
			key = ks[0]
			nested := map[string][]string{
				strings.TrimRight(ks[1], "]"): v,
			}
			val = extractData(nested)
		}

		// check if data exist in dataFormated
		if v, ok := dataFormated[key]; ok {
			refVal := reflect.ValueOf(v)
			if refVal.Kind() == reflect.Map {
				comb := map[string]interface{}{}

				for _, mKey := range refVal.MapKeys() {
					comb[fmt.Sprintf("%s", mKey)] = stringToType(fmt.Sprintf("%v", refVal.MapIndex(mKey)))
					fmt.Println("old", refVal.MapIndex(mKey))
				}

				newRefVal := reflect.ValueOf(val)
				for _, nKey := range newRefVal.MapKeys() {
					comb[fmt.Sprintf("%s", nKey)] = stringToType(fmt.Sprintf("%v", newRefVal.MapIndex(nKey)))
					fmt.Println("new", newRefVal.MapIndex(nKey))
				}
				val = comb
				// fmt.Println(key, val)
			} else {
				vs := []interface{}{v}
				val = append(vs, val)
			}

		}
		dataFormated[key] = val
	}
	return dataFormated
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
