package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type GeoJson struct {
	Type             string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinate       []float32 `json:"coordinate,omitempty" bson:"coordinate,omitempty"`
	FormattedAddress string    `json:"formattedAddress,omitempty" bson:"formattedAddress,omitempty"`
	Street           string    `json:"street,omitempty" bson:"street,omitempty"`
	City             string    `json:"city,omitempty" bson:"city,omitempty"`
	State            string    `json:"state,omitempty" bson:"state,omitempty"`
	Zipcode          string    `json:"zipcode,omitempty" bson:"zipcode,omitempty"`
	Country          string    `json:"country,omitempty" bson:"country,omitempty"`
}

type Bootcamp struct {
	Id            bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string        `json:"name,omitempty" bson:"name,omitempty"`
	Slug          string        `json:"slug,omitempty" bson:"slug,omitempty"`
	Website       string        `json:"website,omitempty" bson:"website,omitempty"`
	Phone         string        `json:"phone,omitempty" bson:"phone,omitempty"`
	Email         string        `json:"email,omitempty" bson:"email,omitempty"`
	Address       string        `json:"address,omitempty" bson:"address,omitempty"`
	Location      GeoJson       `json:"location,omitempty" bson:"location,omitempty"`
	AverageRating float32       `json:"averageRating,omitempty" bson:"averageRating,omitempty"`
	AverageCost   int           `json:"averageCost,omitempty" bson:"averageCost,omitempty"`
	Photo         string        `json:"photo,omitempty" bson:"photo,omitempty"`
	Housing       bool          `json:"housing,omitempty" bson:"housing,omitempty"`
	JobAssistance bool          `json:"jobAssistance,omitempty" bson:"jobAssistance,omitempty"`
	JobGuarantee  bool          `json:"jobGuarantee,omitempty" bson:"jobGuarantee,omitempty"`
	AcceptGi      bool          `json:"acceptGi,omitempty" bson:"acceptGi,omitempty"`
	CreatedAt     time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	User          bson.ObjectId `json:"user,omitempty" bson:"user,omitempty"`
}

type Bootcamps []Bootcamp
