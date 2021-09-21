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

// @desc    Get all users
// @route   GET /api/v1/users
// @access  Private/Admin
func (u *User) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(u.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}
	User := u.connection.Model("User")
	users := []*models.User{}

	query := bson.M{
		"deleted": false,
	}
	err := User.Find(query).Exec(&users)
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(users),
		"data":    users,
	})
}

// @desc    Get single user
// @route   GET /api/v1/users/:id
// @access  Private/Admin
func (u *User) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(u.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid user id format"))
		return
	}

	User := u.connection.Model("User")
	user := &models.User{}

	query := bson.M{
		"_id":     bson.ObjectIdHex(id),
		"deleted": false,
	}
	err := User.FindOne(query).Exec(user)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no user with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// @desc    Create user
// @route   POST /api/v1/users
// @access  Private/Admin
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(u.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	User := u.connection.Model("User")
	user := &models.User{}
	User.New(user)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data request"))
		return
	}
	if valid, issues := user.ValidateCreate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}
	// check if the email is unique (the email of deleted user should not be reuse for resore account feature)
	if n, _ := User.Find(bson.M{"email": user.Email}).Count(); n > 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("email %s was taken, please use the new one", user.Email))
		return
	}
	err = user.HashPassword()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	err = user.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// @desc    Update user
// @route   PUT /api/v1/users/:id
// @access  Private/Admin
func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(u.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}
	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid user id format"))
		return
	}

	User := u.connection.Model("User")
	user := &models.User{}

	query := bson.M{
		"_id":     bson.ObjectIdHex(id),
		"deleted": false,
	}

	err := User.FindOne(query).Exec(user)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no user with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	var d map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data"))
	}

	// The Update method is incompleted so the error is not handled
	// see https://github.com/zebresel-com/mongodm/issues/20
	user.Update(d)

	if valid, issues := user.ValidateCreate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}
	// check if the email is unique (the email of deleted user should not be reuse for resore account feature)
	query = bson.M{
		"email": user.Email,
		"_id": bson.M{
			"$ne": user.Id,
		},
	}
	if n, _ := User.Find(query).Count(); n > 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("email %s was taken, please use the new one", user.Email))
		return
	}

	// hash password if the new password is provided
	if _, ok := d["password"]; ok {
		err = user.HashPassword()
		if err != nil {
			utils.ErrorResponse(w, http.StatusBadRequest, err)
			return
		}
	}
	err = user.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
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
