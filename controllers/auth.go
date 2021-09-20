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

type User struct {
	connection *mongodm.Connection
}

func NewUser(conn *mongodm.Connection) *User {
	return &User{
		connection: conn,
	}
}

// @desc    Register user
// @route   POST /api/v1/auth/register
// @access  Public
func (u *User) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	User := u.connection.Model("User")
	user := &models.User{}
	User.New(user)

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, errors.New("bad data request"))
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

	// send jwt via cookie
	ss, err := utils.GetJwt(user.Id.Hex())
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    ss,
		HttpOnly: true,
	})
	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"token":   ss,
	})
}

// @desc    Login user
// @route   POST /api/v1/auth/login
// @access  Public
func (u *User) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "Login",
	})
}

// @desc    Log user out / clear cookie
// @route   GET /api/v1/auth/logout
// @access  Private
func (u *User) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "Logout",
	})
}

// @desc    Get current logged in user
// @route   GET /api/v1/auth/me
// @access  Private
func (u *User) GetMe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "Get me route",
	})
}

// @desc    Update user details
// @route   PUT /api/v1/auth/updatedetails
// @access  Private
func (u *User) UpdateDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "Update details",
	})
}

// @desc    Update password
// @route   PUT /api/v1/auth/updatepassword
// @access  Private
func (u *User) UpdatePassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "Update password",
	})
}

// @desc    Forgot password
// @route   POST /api/v1/auth/forgotpassword
// @access  Public
func (u *User) ForgotPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"sueecee": true,
		"data":    "Forgot password",
	})
}

// @desc    Reset password
// @route   PUT /api/v1/auth/resetpassword/:token
// @access  Public
func (u *User) ResetPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    "Reset password",
	})
}
