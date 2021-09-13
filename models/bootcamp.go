package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type GeoJson struct {
	Type             string    `json:"type" bson:"type"`
	Coordinate       []float32 `json:"coordinate" bson:"coordinate"`
	FormattedAddress string    `json:"formattedAddress" bson:"formattedAddress"`
	Street           string    `json:"street" bson:"street"`
	City             string    `json:"city" bson:"city"`
	State            string    `json:"state" bson:"state"`
	Zipcode          string    `json:"zipcode" bson:"zipcode"`
	Country          string    `json:"country" bson:"country"`
}

type Bootcamp struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	Name          string        `json:"name" bson:"name"`
	Slug          string        `json:"slug" bson:"slug"`
	Website       string        `json:"website" bson:"website"`
	Phone         string        `json:"phone" bson:"phone"`
	Email         string        `json:"email" bson:"email"`
	Address       string        `json:"address" bson:"address"`
	Location      GeoJson       `json:"location" bson:"location"`
	AverageRating float32       `json:"averageRating" bson:"averageRating"`
	AverageCost   int           `json:"averageCost" bson:"averageCost"`
	Photo         string        `json:"photo" bson:"photo"`
	Housing       bool          `json:"housing" bson:"housing"`
	JobAssistance bool          `json:"jobAssistance" bson:"jobAssistance"`
	JobGuarantee  bool          `json:"jobGuarantee" bson:"jobGuarantee"`
	AcceptGi      bool          `json:"acceptGi" bson:"acceptGi"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	User          bson.ObjectId `json:"user" bson:"user"`
}

type Bootcamps []Bootcamp
