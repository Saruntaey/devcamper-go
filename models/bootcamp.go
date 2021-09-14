package models

import (
	"errors"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GeoJson struct {
	Type             string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates      []float32 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
	FormattedAddress string    `json:"formattedAddress,omitempty" bson:"formattedAddress,omitempty"`
	Street           string    `json:"street,omitempty" bson:"street,omitempty"`
	City             string    `json:"city,omitempty" bson:"city,omitempty"`
	State            string    `json:"state,omitempty" bson:"state,omitempty"`
	Zipcode          string    `json:"zipcode,omitempty" bson:"zipcode,omitempty"`
	Country          string    `json:"country,omitempty" bson:"country,omitempty"`
}

type Bootcamp struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	Name          string        `json:"name,omitempty" bson:"name,omitempty"`
	Slug          string        `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string        `json:"description,omitempty" bson:"description,omitempty"`
	Website       string        `json:"website,omitempty" bson:"website,omitempty"`
	Phone         string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Email         string        `json:"email,omitempty" bson:"email,omitempty"`
	Address       string        `json:"address,omitempty" bson:"address,omitempty"`
	Location      GeoJson       `json:"location,omitempty" bson:"location,omitempty"`
	AverageRating float32       `json:"averageRating" bson:"averageRating"`
	AverageCost   int           `json:"averageCost" bson:"averageCost"`
	Photo         string        `json:"photo,omitempty" bson:"photo,omitempty"`
	Housing       bool          `json:"housing" bson:"housing"`
	JobAssistance bool          `json:"jobAssistance" bson:"jobAssistance"`
	JobGuarantee  bool          `json:"jobGuarantee" bson:"jobGuarantee"`
	AcceptGi      bool          `json:"acceptGi" bson:"acceptGi"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	User          bson.ObjectId `json:"user,omitempty" bson:"user,omitempty"`
}

type Bootcamps []Bootcamp

func (bc *Bootcamp) ValidateData(c *mgo.Collection) error {
	// make field to contain only unique value
	c.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		Background: true,
	})

	// check required field
	if bc.Name == "" {
		return errors.New("please add a name")
	} else if len(bc.Name) > 50 {
		return errors.New("name can not be more than 50 characters")
	}
	if bc.Description == "" {
		return errors.New("please add a description")
	} else if len(bc.Description) > 500 {
		return errors.New("description can not be more than 500 characters")
	}
	if re := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`); !re.Match([]byte(bc.Website)) {
		return errors.New("please use a valid URL with HTTP or HTTPS")
	}
	bc.Id = bson.NewObjectId()
	bc.CreatedAt = time.Now()

	return nil
}
