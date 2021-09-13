package controllers

import (
	"devcamper/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
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
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		res := map[string]interface{}{
			"success": false,
			"data":    nil,
		}
		bs, _ := json.Marshal(res)
		w.Write(bs)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	res := map[string]interface{}{
		"success": true,
		"count":   len(bootcamps),
		"data":    bootcamps,
	}
	bs, _ := json.Marshal(res)
	w.Write(bs)
}

// @desc    Get single bootcamp
// @route   GET /api/v1/bootcamps/:id
// @access  Public
func (bc *Bootcamp) GetBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get single bootcamp"))
}

// @desc    Create bootcamp
// @route   POST /api/v1/bootcamps
// @access  Private
func (bc *Bootcamp) CreateBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("create bootcamp"))
}

// @desc    Update bootcamp
// @route   PUT /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) UpdateBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("update bootcamp"))
}

// @desc    Delete bootcamp
// @route   DELETE /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) DeleteBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("delete bootcamp"))
}
