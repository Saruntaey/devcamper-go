package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/zebresel-com/mongodm"
)

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Name                 string    `json:"name" bson:"name" required:"true"`
	Email                string    `json:"email" bson:"eamil" validation:"email" required:"true"`
	Role                 string    `json:"role" bson:"role"`
	Password             string    `json:"-" bson:"password" minLen:"6"`
	ResetPasswordToken   string    `json:"-" bson:"resetPasswordToken"`
	ResetPasswordExpired time.Time `json:"-" bson:"resetPasswordExpired"`
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
