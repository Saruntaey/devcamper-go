package models

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/zebresel-com/mongodm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Name                 string    `json:"name" bson:"name" required:"true"`
	Email                string    `json:"email" bson:"email" validation:"email" required:"true"`
	Role                 string    `json:"role" bson:"role"`
	PasswordRaw          string    `json:"password,omitempty" bson:"-"`
	PasswordHash         string    `json:"-" bson:"password"`
	ResetPasswordToken   string    `json:"-" bson:"resetPasswordToken,omitempty"`
	ResetPasswordExpired time.Time `json:"-" bson:"resetPasswordExpired,omitempty"`
}

// override validate function to aviod check before save (will check explicitly)
func (u *User) Validate(values ...interface{}) (bool, []error) {
	return true, nil
}

// check data before create user
func (u *User) ValidateCreate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = u.DefaultValidate()

	validationErrors = append(validationErrors, u.validateBothCreateAndUpdate()...)

	// append check here

	return len(validationErrors) == 0, validationErrors
}

// check data before update user
func (u *User) ValidateUpdate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = u.DefaultValidate()

	validationErrors = append(validationErrors, u.validateBothCreateAndUpdate()...)

	// append check here

	return len(validationErrors) == 0, validationErrors
}

// common data to validate
func (u *User) validateBothCreateAndUpdate() []error {
	var validationErrors []error

	// check if role in category
	roles := []string{
		"user",
		"publisher",
	}
	valid := false
	for _, v := range roles {
		if v == u.Role {
			valid = true
			break
		}
	}
	if !valid {
		u.AppendError(&validationErrors, fmt.Sprintf("Please select minimum skill in [ %s ]", strings.Join(roles, ", ")))
	}

	return validationErrors
}

func (u *User) HashPassword() error {
	if len(u.PasswordRaw) < 6 {
		return errors.New("password shoud be at least 6 characters")
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(u.PasswordRaw), bcrypt.DefaultCost)
	if err != nil {
		log.Println("cannot hash password: ", err)
		return errors.New("bad data")
	}
	u.PasswordHash = string(bs)
	// prevent exporting password to json
	u.PasswordRaw = ""
	return nil
}

func (u *User) MatchPassword(pwdRaw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pwdRaw))
	return err == nil
}

func (u *User) GenResetPwdToken() string {
	bs := make([]byte, 20)
	io.ReadFull(rand.Reader, bs)
	h := hmac.New(sha256.New, []byte(os.Getenv("JWT_SECRET")))
	h.Write(bs)
	u.ResetPasswordToken = fmt.Sprintf("%x", h.Sum(nil))
	u.ResetPasswordExpired = time.Now().Add(time.Minute * time.Duration(10))

	return fmt.Sprintf("%x", bs)
}

func (u *User) IsUserInRoles(roles ...string) bool {
	for _, v := range roles {
		if u.Role == v {
			return true
		}
	}
	return false
}
