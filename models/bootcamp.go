package models

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zebresel-com/mongodm"
)

type GeoJson struct {
	Type             string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates      []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	FormattedAddress string    `json:"formattedAddress,omitempty" bson:"formattedAddress,omitempty"`
	Street           string    `json:"street,omitempty" bson:"street,omitempty"`
	City             string    `json:"city,omitempty" bson:"city,omitempty"`
	State            string    `json:"state,omitempty" bson:"state,omitempty"`
	Zipcode          string    `json:"zipcode,omitempty" bson:"zipcode,omitempty"`
	Country          string    `json:"country,omitempty" bson:"country,omitempty"`
}

type Bootcamp struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Name                 string        `json:"name" bson:"name" required:"true" maxLen:"50"`
	Slug                 string        `json:"slug" bson:"slug"`
	Description          string        `json:"description" bson:"description" required:"true" maxLen:"500"`
	Website              string        `json:"website" bson:"website"`
	Phone                string        `json:"phone" bson:"phone" maxLen:"20"`
	Email                string        `json:"email" bson:"email" validation:"email"`
	Address              string        `json:"address,omitempty" bson:"address,omitempty"`
	Location             *GeoJson      `json:"location" bson:"location"`
	Careers              []string      `json:"careers" bson:"careers" required:"true"`
	AverageRating        float64       `json:"averageRating" bson:"averageRating"`
	AverageCost          int           `json:"averageCost" bson:"averageCost"`
	Photo                string        `json:"photo" bson:"photo"`
	Housing              bool          `json:"housing" bson:"housing"`
	JobAssistance        bool          `json:"jobAssistance" bson:"jobAssistance"`
	JobGuarantee         bool          `json:"jobGuarantee" bson:"jobGuarantee"`
	AcceptGi             bool          `json:"acceptGi" bson:"acceptGi"`
	Courses              []interface{} `json:"courses,omitempty" bson:"-"`
	User                 interface{}   `json:"user" bson:"user" model:"User" relation:"11" autosave:"true" required:"true"`
}

// override validate function to aviod check before save (will check explicitly)
func (bc *Bootcamp) Validate(values ...interface{}) (bool, []error) {
	return true, nil
}

// check data before create bootcamp
func (bc *Bootcamp) ValidateCreate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = bc.DefaultValidate()

	validationErrors = append(validationErrors, bc.validateBothCreateAndUpdate()...)

	// check if address is proveded
	if bc.Address == "" {
		bc.AppendError(&validationErrors, "Please add an address")
	}

	return len(validationErrors) == 0, validationErrors
}

// check data before update bootcamp
func (bc *Bootcamp) ValidateUpdate() (bool, []error) {
	var validationErrors []error

	_, validationErrors = bc.DefaultValidate()

	validationErrors = append(validationErrors, bc.validateBothCreateAndUpdate()...)

	return len(validationErrors) == 0, validationErrors
}

// common data to validate
func (bc *Bootcamp) validateBothCreateAndUpdate() []error {
	var validationErrors []error

	// check website format
	if regex := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`); !regex.Match([]byte(bc.Website)) {
		bc.AppendError(&validationErrors, "Please use a valid URL with HTTP or HTTPS")
	}

	// check careers in list
	careers := []string{
		"Web Development",
		"Mobile Development",
		"UI/UX",
		"Data Science",
		"Business",
		"Other",
	}
	inCareers := true
	for _, v := range bc.Careers {
		for i, n := 0, len(careers); i < n; i++ {
			if v == careers[i] {
				break
			}
			if i == n-1 && v != careers[i] {
				inCareers = false
				break
			}
		}
		if !inCareers {
			bc.AppendError(&validationErrors, fmt.Sprintf("Please select careers in [ %s ]", strings.Join(careers, ", ")))
			break
		}
	}

	// check averageRating range
	if bc.AverageRating < 0 {
		bc.AppendError(&validationErrors, "Rating must be at least 1")
	} else if bc.AverageRating > 10 {
		bc.AppendError(&validationErrors, "Rating cannot be more than 10")
	}

	return validationErrors
}
