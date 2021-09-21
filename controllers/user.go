package controllers

import (
	"devcamper/utils"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// @desc    Get all users
// @route   GET /api/v1/users
// @access  Private/Admin
func (u *User) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "get all users",
	})
}

// @desc    Get single user
// @route   GET /api/v1/users/:id
// @access  Private/Admin
func (u *User) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("get user id: %s", id),
	})
}

// @desc    Create user
// @route   POST /api/v1/users
// @access  Private/Admin
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "create user",
	})
}

// @desc    Update user
// @route   PUT /api/v1/users/:id
// @access  Private/Admin
func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("update users id: %s", id),
	})
}

// @desc    Delete user
// @route   DELETE /api/v1/users/:id
// @access  Private/Admin
func (u *User) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    fmt.Sprintf("delete users id: %s", id),
	})
}
