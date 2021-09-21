package controllers

import (
	"devcamper/utils"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
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
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "get all reviews",
	})
}

// @desc    Get reviews
// @route   GET /api/v1/bootcamps/:id/reviews
// @access  Public
func (rw *Review) GetReviewsInBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("get all reviews in bootcamp id %s", id),
	})
}

// @desc    Get single review
// @route   GET /api/v1/reviews/:id
// @access  Public
func (rw *Review) GetReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("get review id %s", id),
	})
}

// @desc    Add review
// @route   GET /api/v1/bootcamps/:id/reviews
// @access  Private
func (rw *Review) AddReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("add review for bootcamp id %s", id),
	})
}

// @desc    Update review
// @route   PUT /api/v1/reviews/:id
// @access  Private
func (rw *Review) UpdateReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("update review id %s", id),
	})
}

// @desc    Delete review
// @route   DELETE /api/v1/reviews/:id
// @access  Private
func (rw *Review) DeleteReview(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("delete review id %s", id),
	})
}
