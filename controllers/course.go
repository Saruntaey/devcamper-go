package controllers

import (
	"devcamper/models"
	"devcamper/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type Course struct {
	connection *mongodm.Connection
}

func NewCourse(conn *mongodm.Connection) *Course {
	return &Course{
		connection: conn,
	}
}

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps
// @access  Public
func (c *Course) GetCourses(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Course := c.connection.Model("Course")
	courses := []*models.Course{}

	err := Course.Find(bson.M{"deleted": false}).Exec(&courses)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(courses),
		"data":    courses,
	})
}

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps/:id/courses
// @access  Public
func (c *Course) GetCoursesInBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bootcampId := ps.ByName("id")
	if !bson.IsObjectIdHex(bootcampId) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid course id format"))
		return
	}

	Bootcamp := c.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}
	err := Bootcamp.FindId(bson.ObjectIdHex(bootcampId)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found bootcamp with id of %s", bootcampId))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if bootcamp.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("the bootcamp with the id of %s was deleted", bootcampId))
		return
	}

	Course := c.connection.Model("Course")
	courses := []*models.Course{}

	query := bson.M{
		"bootcamp": bson.ObjectIdHex(bootcampId),
		"deleted":  false,
	}
	err = Course.Find(query).Exec(&courses)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(courses),
		"data":    courses,
	})
}

// @desc    Get course
// @route   GET /api/v1/courses/:id
// @access  Public
func (c *Course) GetCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid course id format"))
		return
	}

	Course := c.connection.Model("Course")
	course := &models.Course{}

	err := Course.FindId(bson.ObjectIdHex(id)).Exec(course)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found course with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if course.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this course was deleted"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    course,
	})
}

// @desc    Add course
// @route   POST /api/v1/bootcamps/:id/courses
// @access  Private
func (c *Course) AddCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bootcampId := ps.ByName("id")
	if !bson.IsObjectIdHex(bootcampId) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := c.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}
	err := Bootcamp.FindId(bson.ObjectIdHex(bootcampId)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found bootcamp with id of %s", bootcampId))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if bootcamp.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("the bootcamp with the id of %s was deleted", bootcampId))
		return
	}

	Course := c.connection.Model("Course")
	course := &models.Course{}
	Course.New(course)

	json.NewDecoder(r.Body).Decode(course)
	course.Bootcamp = bson.ObjectIdHex(bootcampId)
	if valid, issue := course.ValidateCreate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issue...)
		return
	}
	course.Save()

	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    course,
	})
}

// @desc    Update course
// @route   PUT /api/v1/courses/:id
// @access  Private
func (c *Course) UpdateCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid course id format"))
		return
	}

	Course := c.connection.Model("Course")
	course := &models.Course{}

	err := Course.FindId(bson.ObjectIdHex(id)).Exec(course)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found course with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if course.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this course was deleted"))
		return
	}
	var data map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data"))
		return
	}
	// delete unexpected field
	delete(data, "bootcamp")

	// The Update method is incompleted so the error is not handled
	// see https://github.com/zebresel-com/mongodm/issues/20
	course.Update(data)

	if valid, issues := course.ValidateUpdate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}

	course.Save()

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    course,
	})
}

// @desc    Delete course
// @route   DELETE /api/v1/courses/:id
// @access  Private
func (c *Course) DeleteCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid course id format"))
		return
	}

	Course := c.connection.Model("Course")
	course := &models.Course{}

	err := Course.FindId(bson.ObjectIdHex(id)).Exec(course)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no course with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if course.Deleted {
		// should not found the deleted course
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no course with id of %s", id))
		return
	}
	course.SetDeleted(true)
	course.Save()
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    nil,
	})
}
