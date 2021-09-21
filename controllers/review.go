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

type Review struct {
	connection *mongodm.Connection
}

func NewReview(conn *mongodm.Connection) *Review {
	return &Review{
		connection: conn,
	}
}

// @desc    Get reviews
// @route   GET /api/v1/reviews
// @access  Public
func (rw *Review) GetReviews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Review := rw.connection.Model("Review")
	reviews := []*models.Review{}

	err := Review.Find(bson.M{"deleted": false}).Exec(&reviews)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(reviews),
		"data":    reviews,
	})
}

// @desc    Get reviews
// @route   GET /api/v1/bootcamps/:id/reviews
// @access  Public
func (rw *Review) GetReviewsInBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bootcampId := ps.ByName("id")
	if !bson.IsObjectIdHex(bootcampId) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := rw.connection.Model("Bootcamp")
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

	Review := rw.connection.Model("Review")
	reviews := []*models.Course{}

	query := bson.M{
		"bootcamp": bson.ObjectIdHex(bootcampId),
		"deleted":  false,
	}
	err = Review.Find(query).Exec(&reviews)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(reviews),
		"data":    reviews,
	})
}

// @desc    Get single review
// @route   GET /api/v1/reviews/:id
// @access  Public
func (rw *Review) GetReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid review id format"))
		return
	}

	Review := rw.connection.Model("Review")
	review := &models.Review{}

	err := Review.FindId(bson.ObjectIdHex(id)).Exec(review)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found review with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if review.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this review was deleted"))
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    review,
	})
}

// @desc    Add review
// @route   GET /api/v1/bootcamps/:id/reviews
// @access  Private
func (rw *Review) AddReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(rw.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("user", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	bootcampId := ps.ByName("id")
	if !bson.IsObjectIdHex(bootcampId) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := rw.connection.Model("Bootcamp")
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

	Review := rw.connection.Model("Review")
	review := &models.Review{}
	Review.New(review)

	json.NewDecoder(r.Body).Decode(review)
	review.Bootcamp = bson.ObjectIdHex(bootcampId)
	review.User = cUser.Id
	if valid, issue := review.ValidateCreate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issue...)
		return
	}
	review.Save()

	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    review,
	})
}

// @desc    Update review
// @route   PUT /api/v1/reviews/:id
// @access  Private
func (rw *Review) UpdateReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(rw.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("user", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid review id format"))
		return
	}

	Review := rw.connection.Model("Review")
	review := &models.Review{}

	err := Review.FindId(bson.ObjectIdHex(id)).Exec(review)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found review with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if review.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this review was deleted"))
		return
	}

	if review.User != cUser.Id && cUser.Role != "admin" {
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
	review.Update(data)

	if valid, issues := review.ValidateUpdate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}

	review.Save()

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    review,
	})
}

// @desc    Delete review
// @route   DELETE /api/v1/reviews/:id
// @access  Private
func (rw *Review) DeleteReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(rw.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("user", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid course id format"))
		return
	}

	Review := rw.connection.Model("Review")
	review := &models.Review{}

	err := Review.FindId(bson.ObjectIdHex(id)).Exec(review)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no review with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	if review.Deleted {
		// should not found the deleted course
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no review with id of %s", id))
		return
	}

	if review.User != cUser.Id && cUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("you do not have permission"))
		return
	}

	review.SetDeleted(true)
	review.Save()
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    nil,
	})
}
