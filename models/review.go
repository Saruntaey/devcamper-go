package models

import (
	"github.com/zebresel-com/mongodm"
)

type Review struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Title                string      `json:"title" bson:"title" required:"true" maxLen:"50"`
	Text                 string      `json:"text" bson:"text" required:"true" maxLen:"100"`
	Rating               int         `json:"rating" bson:"rating"`
	Bootcamp             interface{} `json:"bootcamp" bson:"bootcamp" model:"Bootcamp" relation:"11" autosave:"true" required:"true"`
	User                 interface{} `json:"user" bson:"user" model:"User" relation:"11" autosave:"true" required:"true"`
}

// override validate function to aviod check before save (will check explicitly)
func (r *Review) Validate(values ...interface{}) (bool, []error) {
	return true, nil
}

// check data before create bootcamp
func (r *Review) ValidateCreate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = r.DefaultValidate()

	validationErrors = append(validationErrors, r.validateBothCreateAndUpdate()...)

	// append check here

	return len(validationErrors) == 0, validationErrors
}

// check data before update bootcamp
func (r *Review) ValidateUpdate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = r.DefaultValidate()

	validationErrors = append(validationErrors, r.validateBothCreateAndUpdate()...)

	// append check here

	return len(validationErrors) == 0, validationErrors
}

// common data to validate
func (r *Review) validateBothCreateAndUpdate() []error {
	var validationErrors []error

	// check rating range
	if r.Rating < 0 {
		r.AppendError(&validationErrors, "Rating must be at least 1")
	} else if r.Rating > 10 {
		r.AppendError(&validationErrors, "Rating cannot be more than 10")
	}

	return validationErrors
}
