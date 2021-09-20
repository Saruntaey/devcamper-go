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

type User struct {
	connection *mongodm.Connection
}

type LoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateDetails struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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
	sendToken(w, user)
}

// @desc    Login user
// @route   POST /api/v1/auth/login
// @access  Public
func (u *User) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginDetails := LoginDetails{}
	err := json.NewDecoder(r.Body).Decode(&loginDetails)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data"))
		return
	}

	// check if email and password is provided
	if len(loginDetails.Email) == 0 || len(loginDetails.Password) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("please provied email and password"))
		return
	}

	User := u.connection.Model("User")
	user := &models.User{}

	err = User.FindOne(bson.M{"email": loginDetails.Email}).Exec(user)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	if !user.MatchPassword(loginDetails.Password) {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("invalid email or password"))
		return
	}

	sendToken(w, user)
}

// @desc    Log user out / clear cookie
// @route   GET /api/v1/auth/logout
// @access  Private
func (u *User) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := "none"
	// delete cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		MaxAge:   -1,
	})
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   token,
	})
}

// @desc    Get current logged in user
// @route   GET /api/v1/auth/me
// @access  Private
func (u *User) GetMe(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := getCurrentUser(u.connection, r)
	if user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

// @desc    Update user details
// @route   PUT /api/v1/auth/updatedetails
// @access  Private
func (u *User) UpdateDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := getCurrentUser(u.connection, r)
	if user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	updateDetails := UpdateDetails{}
	json.NewDecoder(r.Body).Decode(&updateDetails)

	if len(updateDetails.Email) > 0 {
		user.Email = updateDetails.Email
	}
	if len(updateDetails.Name) > 0 {
		user.Name = updateDetails.Name
	}

	if valid, issues := user.ValidateUpdate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}

	User := u.connection.Model("User")
	// check if the email is unique (the email of deleted user should not be reuse for resore account feature)
	query := bson.M{
		"email": user.Email,
		"_id": bson.M{
			"$ne": user.Id,
		},
	}
	if n, _ := User.Find(query).Count(); n > 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("email %s was taken, please use the new one", user.Email))
		return
	}
	err := user.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
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

func sendToken(w http.ResponseWriter, user *models.User) {
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

func getCurrentUser(conn *mongodm.Connection, r *http.Request) *models.User {
	var token string
	// grab token from header
	if auth := strings.Split(r.Header.Get("Authorization"), " "); auth[0] == "Bearer" {
		token = auth[1]
		// grab token from cookie
	} else if c, err := r.Cookie("token"); err == nil {
		token = c.Value
	}

	// no token
	if len(token) == 0 {
		return nil
	}

	// validate token
	payload, err := utils.ParseJwt(token)
	if err != nil {
		return nil
	}

	// find user
	userId := payload.(*utils.Payload).Id
	User := conn.Model("User")
	user := &models.User{}

	query := bson.M{
		"_id":     bson.ObjectIdHex(userId),
		"deleted": false,
	}
	err = User.FindOne(query).Exec(user)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		return nil
	}
	return user
}
