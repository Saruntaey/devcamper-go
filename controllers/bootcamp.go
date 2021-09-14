package controllers

import (
	"devcamper/models"
	"devcamper/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Bootcamp struct {
	session *mgo.Session
	db      string
	c       string
}

func NewBootcamp(s *mgo.Session) *Bootcamp {
	return &Bootcamp{s, os.Getenv("MONGO_DB"), "bootcamps"}
}

func (bc *Bootcamp) collection() *mgo.Collection {
	return bc.session.DB(bc.db).C(bc.c)
}

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps
// @access  Public
func (bc *Bootcamp) GetBootcamps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	var bootcamps models.Bootcamps
	err := bc.collection().Find(nil).All(&bootcamps)
	if err != nil {
		fmt.Println("GetBootcamps err: ", err)

		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(bootcamps),
		"data":    bootcamps,
	})
}

// @desc    Get single bootcamp
// @route   GET /api/v1/bootcamps/:id
// @access  Public
func (bc *Bootcamp) GetBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")

	// validate id
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bootcamp id not in correct format"))
		return
	}

	var bootcamp models.Bootcamp

	err := bc.collection().FindId(bson.ObjectIdHex(id)).One(&bootcamp)
	if err != nil {
		fmt.Println("GetBootcamp err: ", err)
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found bootcamp with the id of %s", id))
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
	var bootcamp models.Bootcamp
	err := json.NewDecoder(r.Body).Decode(&bootcamp)
	if err != nil {
		fmt.Println("CreateBootcamp err: ", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	err = bootcamp.ValidateData(bc.collection(), true)

	if err != nil {
		fmt.Println("CreateBootcamp err: ", err)
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = bc.collection().Insert(bootcamp)
	if err != nil {
		fmt.Println("CreateBootcamp err: ", err)
		if mgo.IsDup(err) {
			utils.ErrorResponse(w, http.StatusBadRequest, errors.New("duplicate field value enter"))
			return
		}
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
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
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bootcamp id not in correct format"))
		return
	}
	var newBootcamp models.Bootcamp
	json.NewDecoder(r.Body).Decode(&newBootcamp)
	err := bc.collection().UpdateId(bson.ObjectIdHex(id), newBootcamp)
	if err != nil {
		fmt.Println("UpdateBootcamp err: ", err)
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    id,
	})
	// var bootcamp models.Bootcamp
	// err := bc.collection().FindId(bson.ObjectIdHex(id)).One(&bootcamp)
	// if err != nil {
	// 	fmt.Println("UpdateBootcamp err: ", err)
	// 	utils.SendJSON(w, http.StatusInternalServerError, errors.New("server error"))
	// 	return
	// }
	// var newBootcamp models.Bootcamp
	// err = json.NewDecoder(r.Body).Decode(&newBootcamp)
	// if err != nil {
	// 	fmt.Print("UpdateBootcamp err: ", err)
	// 	utils.SendJSON(w, http.StatusInternalServerError, errors.New("server error"))
	// 	return
	// }
}

// @desc    Delete bootcamp
// @route   DELETE /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) DeleteBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete bootcamp"))
}
