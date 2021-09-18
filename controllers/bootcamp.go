package controllers

import (
	"devcamper/models"
	"devcamper/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type Bootcamp struct {
	connection *mongodm.Connection
	// Model      *Bootcamp
	// Models     []*Bootcamp
}

func NewBootcamp(conn *mongodm.Connection) *Bootcamp {
	return &Bootcamp{
		connection: conn,
	}
}

// func (bc *Bootcamp) GetModel() *mongodm.Model {
// 	return bc.connection.Model("Bootcamp")
// }

// func (bc *Bootcamp) MallocOne() {
// 	bc.Model = &Bootcamp{}
// }

// func (bc *Bootcamp) MallocMany() {
// 	bc.Models = []*Bootcamp{}
// }

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps
// @access  Public
func (bc *Bootcamp) GetBootcamps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad request data"))
		return
	}
	// bootcamps := []*models.Bootcamp{}
	// respData, httpStatus, err := utils.AdvanceQuery(r.Form, bc.connection.Model("Bootcamp"), &bootcamps)

	// respData, httpStatus, err := models.AdvanceQuery(r.Form, &Bootcamp{})
	respData, httpStatus, err := models.AdvanceQuery(r.Form, bc.connection.Model("Bootcamp"))
	if err != nil {
		utils.ErrorResponse(w, httpStatus, err)
		return
	}

	utils.SendJSON(w, httpStatus, respData)
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
