package controllers

import (
	"devcamper/models"
	"devcamper/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	// parse form
	err := r.ParseForm()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad request data"))
		return
	}

	// create advance query
	query, pagination, err := models.AdvanceQuery(r.Form, c.connection.Model("Course"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad request data"))
		return
	}

	courses := []*models.Course{}

	// execute query
	err = query.Exec(&courses)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	// prepare response data
	respData := map[string]interface{}{
		"success":    true,
		"count":      len(courses),
		"pagination": pagination,
	}

	// hide data that user not request
	selectField := r.Form["select"]
	if len(selectField) != 0 {
		selects := strings.Split(selectField[0], ",")
		respData["data"] = models.ExtractSelectField(courses, selects)
	} else {
		respData["data"] = courses
	}

	utils.SendJSON(w, http.StatusOK, respData)
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
	cUser := getCurrentUser(c.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("publisher", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

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

	if bootcamp.User != cUser.Id && cUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("you do not have permission"))
		return
	}

	Course := c.connection.Model("Course")
	course := &models.Course{}
	Course.New(course)

	json.NewDecoder(r.Body).Decode(course)
	course.Bootcamp = bson.ObjectIdHex(bootcampId)
	course.User = cUser.Id
	if valid, issue := course.ValidateCreate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issue...)
		return
	}
	course.Save()

	// update averageCost for bootcamp
	bootcamp.AverageCost = getAvgCost(c.connection, bootcamp.Id)
	bootcamp.Save()

	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    course,
	})
}

// @desc    Update course
// @route   PUT /api/v1/courses/:id
// @access  Private
func (c *Course) UpdateCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(c.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("publisher", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

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

	if course.User != cUser.Id && cUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("you do not have permission"))
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
	delete(data, "user")

	// The Update method is incompleted so the error is not handled
	// see https://github.com/zebresel-com/mongodm/issues/20
	course.Update(data)

	if valid, issues := course.ValidateUpdate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}

	course.Save()

	if _, ok := data["tuition"]; ok {
		// update averageCost for bootcamp
		Bootcamp := c.connection.Model("Bootcamp")
		bootcamp := &models.Bootcamp{}

		Bootcamp.FindId(course.Bootcamp.(bson.ObjectId)).Exec(bootcamp)
		bootcamp.AverageCost = getAvgCost(c.connection, bootcamp.Id)
		bootcamp.Save()
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    course,
	})
}

// @desc    Delete course
// @route   DELETE /api/v1/courses/:id
// @access  Private
func (c *Course) DeleteCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(c.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("publisher", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

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

	if course.User != cUser.Id && cUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("you do not have permission"))
		return
	}

	course.SetDeleted(true)
	course.Save()

	// update averageCost for bootcamp
	Bootcamp := c.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	Bootcamp.FindId(course.Bootcamp.(bson.ObjectId)).Exec(bootcamp)
	bootcamp.AverageCost = getAvgCost(c.connection, bootcamp.Id)
	bootcamp.Save()

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    nil,
	})
}

func getAvgCost(conn *mongodm.Connection, bootcampId bson.ObjectId) int {
	var avg int
	Course := conn.Model("Course")
	courses := []*models.Course{}

	query := bson.M{
		"bootcamp": bootcampId,
		"deleted":  false,
	}
	err := Course.Find(query).Exec(&courses)
	if err != nil {
		return avg
	}

	sum := float64(0)
	for _, v := range courses {
		sum += v.Tuition
	}
	// force last digit to be zero
	avg = int((sum/float64(len(courses)))/10) * 10
	return avg
}
