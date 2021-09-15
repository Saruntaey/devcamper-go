package models

import (
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
	Name                 string   `json:"name" bson:"name" required:"true"`
	Slug                 string   `json:"slug" bson:"slug"`
	Description          string   `json:"description" bson:"description" required:"true"`
	Website              string   `json:"website" bson:"website"`
	Phone                string   `json:"phone" bson:"phone" maxLen:"20"`
	Email                string   `json:"email" bson:"email" validation:"email"`
	Address              string   `json:"address,omitempty" bson:"address,omitempty" required:"true"`
	Location             *GeoJson `json:"location" bson:"location"`
	Careers              []string `json:"careers" bson:"careers" required:"true"`
	AverageRating        float64  `json:"averageRating" bson:"averageRating"`
	AverageCost          int      `json:"averageCost" bson:"averageCost"`
	Photo                string   `json:"photo" bson:"photo"`
	Housing              bool     `json:"housing" bson:"housing"`
	JobAssistance        bool     `json:"jobAssistance" bson:"jobAssistance"`
	JobGuarantee         bool     `json:"jobGuarantee" bson:"jobGuarantee"`
	AcceptGi             bool     `json:"acceptGi" bson:"acceptGi"`

	// User          bson.ObjectId `json:"user,omitempty" bson:"user,omitempty"`
}
