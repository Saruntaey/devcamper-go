package controllers

import (
	"devcamper/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
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
	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    "Register",
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
