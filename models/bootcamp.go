package models

import (
	"devcamper/utils"
	"errors"
	"regexp"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	Id            bson.ObjectId `json:"id" bson:"_id"`
	Name          string        `json:"name,omitempty" bson:"name,omitempty"`
	Slug          string        `json:"slug,omitempty" bson:"slug,omitempty"`
	Description   string        `json:"description,omitempty" bson:"description,omitempty"`
	Website       string        `json:"website,omitempty" bson:"website,omitempty"`
	Phone         string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Email         string        `json:"email,omitempty" bson:"email,omitempty"`
	Address       string        `json:"address,omitempty" bson:"address,omitempty"`
	Location      GeoJson       `json:"location,omitempty" bson:"location,omitempty"`
	AverageRating float64       `json:"averageRating" bson:"averageRating"`
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

func (bc *Bootcamp) ValidateData(c *mgo.Collection, init bool) error {
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
	if len(bc.Phone) > 20 {
		return errors.New("phone can not be more than 20 characters")
	}
	if re := regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`); !re.Match([]byte(bc.Email)) {
		return errors.New("please add a vaild email")
	}
	if bc.Photo == "" {
		bc.Photo = "no-photo.jpg"
	}

	// to be check for newly create bootcamp
	if init {
		// address is required only on create new bootcamp (not reqried on update)
		if bc.Address == "" {
			return errors.New("pleas add an address")
		}

		loc := utils.GetLocation(bc.Address)
		// convert data to GeoJson type
		l := loc.Results[0].Locations[0]

		bc.Location = GeoJson{
			"Point",
			[]float64{l.LatLng.Lng, l.LatLng.Lat},
			l.Street + ", " + l.City + ", " + l.State + " " + l.Zipcode + ", " + l.Country,
			l.Street,
			l.City,
			l.State,
			l.Zipcode,
			l.Country,
		}

		bc.Address = ""
	}

	bc.Id = bson.NewObjectId()
	bc.CreatedAt = time.Now()

	return nil
}
