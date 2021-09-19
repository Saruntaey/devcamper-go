package models

import (
	"fmt"
	"strings"

	"github.com/zebresel-com/mongodm"
)

type Course struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Title                string      `json:"title" bson:"title" required:"true"`
	Description          string      `json:"description" bson:"description" required:"true"`
	Weeks                int         `json:"weeks" bson:"weeks" required:"true"`
	Tuition              float64     `json:"tuition" bson:"tuition" required:"true"`
	MinimumSkill         string      `json:"minimumSkill" bson:"minimumSkill" required:"true"`
	ScholarshipAvailable bool        `json:"scholarshipAvailable" bson:"scholarshipAvailable"`
	Bootcamp             interface{} `json:"bootcamp" bson:"bootcamp" model:"Bootcamp" relation:"11" autosave:"true" required:"true"`
	// User interface{}
}

// override validate function to aviod check before save (will check explicitly)
func (c *Course) Validate(values ...interface{}) (bool, []error) {
	return true, nil
}

// check data before create bootcamp
func (c *Course) ValidateCreate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = c.DefaultValidate()

	validationErrors = append(validationErrors, c.validateBothCreateAndUpdate()...)

	return len(validationErrors) == 0, validationErrors
}

// check data before update course
func (c *Course) ValidateUpdate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = c.DefaultValidate()

	validationErrors = append(validationErrors, c.validateBothCreateAndUpdate()...)

	return len(validationErrors) == 0, validationErrors
}

// common data to validate
func (c *Course) validateBothCreateAndUpdate() []error {
	var validationErrors []error

	// check if the minimum skill in category
	skills := []string{
		"beginner",
		"intermediate",
		"advanced",
	}
	valid := false
	for _, v := range skills {
		if v == c.MinimumSkill {
			valid = true
			break
		}
	}
	if !valid {
		c.AppendError(&validationErrors, fmt.Sprintf("Please select minimum skill in [ %s ]", strings.Join(skills, ", ")))
	}

	return validationErrors
}
