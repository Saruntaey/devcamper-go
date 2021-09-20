package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"devcamper/models"
	"devcamper/utils"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2"
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

type UpdatePassword struct {
	CPwd string `json:"currentPassword"`
	NPwd string `json:"newPassword"`
}

type ForgotPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Pwd string `json:"password"`
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
	user := getCurrentUser(u.connection, r)
	if user == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	updatePwd := UpdatePassword{}
	json.NewDecoder(r.Body).Decode(&updatePwd)
	// check if the data is provided
	if len(updatePwd.CPwd) == 0 || len(updatePwd.NPwd) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("please provied current password and new password"))
		return
	}

	// check current password match
	if !user.MatchPassword(updatePwd.CPwd) {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("current password not match"))
		return
	}

	user.PasswordRaw = updatePwd.NPwd
	err := user.HashPassword()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
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

// @desc    Forgot password
// @route   POST /api/v1/auth/forgotpassword
// @access  Public
func (u *User) ForgotPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	forgotPwd := ForgotPassword{}
	json.NewDecoder(r.Body).Decode(&forgotPwd)
	// check if the email is provided
	if len(forgotPwd.Email) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("please provide email"))
		return
	}

	User := u.connection.Model("User")
	user := &models.User{}

	query := bson.M{
		"email":   forgotPwd.Email,
		"deleted": false,
	}
	User.FindOne(query).Exec(user)
	token := user.GenResetPwdToken()
	err := user.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
	}

	resetPwdURL := fmt.Sprintf("/api/v1/auth/resetpassword/%s", token)
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"sueecee": true,
		"data":    resetPwdURL,
	})
}

// @desc    Reset password
// @route   PUT /api/v1/auth/resetpassword/:token
// @access  Public
func (u *User) ResetPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := ps.ByName("token")
	h := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET")))
	bs, err := hex.DecodeString(token)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	_, err = h.Write(bs)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}
	x := fmt.Sprintf("%x", h.Sum(nil))

	User := u.connection.Model("User")
	user := &models.User{}

	query := bson.M{
		"resetPasswordToken": x,
		"resetPasswordExpired": bson.M{
			"$gt": time.Now(),
		},
		"deleted": false,
	}
	err = User.FindOne(query).Exec(user)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("your token is expired"))
		return
	} else if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	resetPwd := ResetPassword{}
	json.NewDecoder(r.Body).Decode(&resetPwd)
	// check if the password is provided
	if len(resetPwd.Pwd) == 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("please provide a new password"))
		return
	}
	user.PasswordRaw = resetPwd.Pwd
	err = user.HashPassword()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = user.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
	}
	// remove field resetPasswordToken and resetPasswordExpired from DB
	change := mgo.Change{
		Update: bson.M{
			//TO DO
			"$unset": bson.M{
				"resetPasswordToken":   "",
				"resetPasswordExpired": "",
			},
		},
	}
	u.connection.Session.DB(os.Getenv("MONGO_DB")).C("users").FindId(user.Id).Apply(change, user)
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    user,
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
