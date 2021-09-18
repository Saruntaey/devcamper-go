package controllers

import (
	"devcamper/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
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
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "get all courses",
	})
}

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps/:id/courses
// @access  Public
func (c *Course) GetCoursesInBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bootcampId := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"data":     "get all courses in bootcamp",
		"bootcamp": bootcampId,
	})
}

// @desc    Get course
// @route   GET /api/v1/courses/:id
// @access  Public
func (c *Course) GetCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "get single course",
		"id":      id,
	})
}

// @desc    Add course
// @route   POST /api/v1/bootcamps/:id/courses
// @access  Private
func (c *Course) AddCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bootcampId := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"data":     "add course",
		"bootcamp": bootcampId,
	})
}

// @desc    Update course
// @route   PUT /api/v1/courses/:id
// @access  Private
func (c *Course) UpdateCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "update course",
		"id":      id,
	})
}

// @desc    Delete course
// @route   DELETE /api/v1/courses/:id
// @access  Private
func (c *Course) DeleteCourse(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "delete course",
		"id":      id,
	})
}
