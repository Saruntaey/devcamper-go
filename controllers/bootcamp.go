package controllers

import (
	"devcamper/models"
	"devcamper/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type Bootcamp struct {
	connection *mongodm.Connection
}

func NewBootcamp(conn *mongodm.Connection) *Bootcamp {
	return &Bootcamp{
		connection: conn,
	}
}

// @desc    Get all bootcamps
// @route   GET /api/v1/bootcamps
// @access  Public
func (bc *Bootcamp) GetBootcamps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad request data"))
		return
	}

	// create advance query
	query, pagination, err := models.AdvanceQuery(r.Form, bc.connection.Model("Bootcamp"))
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad request data"))
		return
	}

	// execute query
	bootcamps := []*models.Bootcamp{}
	err = query.Exec(&bootcamps)
	if err != nil {
		log.Println(err)
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	// grab courses for each bootcamp (virtual field)
	Course := bc.connection.Model("Course")
	for _, bootcamp := range bootcamps {
		courses := []*models.Course{}
		query := bson.M{
			"bootcamp": bootcamp.Id,
			"deleted":  false,
		}
		err := Course.Find(query).Exec(&courses)
		if err != nil {
			continue
		}
		tmp := []interface{}{}
		for _, v := range courses {
			tmp = append(tmp, v)
		}
		bootcamp.Courses = tmp
	}

	// prepare response data
	respData := map[string]interface{}{
		"success":    true,
		"count":      len(bootcamps),
		"pagination": pagination,
	}

	// hide data that user not request
	selectField := r.Form["select"]
	if len(selectField) != 0 {
		selects := strings.Split(selectField[0], ",")
		respData["data"] = models.ExtractSelectField(bootcamps, selects)
	} else {
		respData["data"] = bootcamps
	}

	utils.SendJSON(w, http.StatusOK, respData)
}

// @desc    Get single bootcamp
// @route   GET /api/v1/bootcamps/:id
// @access  Public
func (bc *Bootcamp) GetBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}
	err := Bootcamp.FindId(bson.ObjectIdHex(id)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("not found bootcamp with id of %s", id))
		return
	} else if bootcamp.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this bootcamp was deleted"))
		return
	} else if err != nil {
		utils.ErrorHandler(w, err)
		return
	}
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bootcamp,
	})
}

// @desc    Create bootcamp
// @route   POST /api/v1/bootcamps
// @access  Private
func (bc *Bootcamp) CreateBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(bc.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("publisher", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	Bootcamp.New(bootcamp)
	err := json.NewDecoder(r.Body).Decode(bootcamp)
	if err != nil {
		log.Println("bad data")
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data"))
		return
	}

	bootcamp.User = cUser.Id
	if valid, issues := bootcamp.ValidateCreate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}

	loc := utils.GetLocation(bootcamp.Address)
	tmp := loc.Results[0].Locations[0]
	geo := &models.GeoJson{
		Type:             "Point",
		Coordinates:      []float64{tmp.LatLng.Lng, tmp.LatLng.Lat},
		FormattedAddress: fmt.Sprintf("%s, %s, %s %s, %s", tmp.Street, tmp.City, tmp.State, tmp.Zipcode, tmp.Country),
		Street:           tmp.Street,
		City:             tmp.City,
		State:            tmp.State,
		Zipcode:          tmp.Zipcode,
		Country:          tmp.Country,
	}
	bootcamp.Location = geo
	bootcamp.Address = ""
	bootcamp.Photo = "no-photo.jpg"
	bootcamp.Slug = strings.Join(strings.Split(strings.ToLower(bootcamp.Name), " "), "-")

	err = bootcamp.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    bootcamp,
	})

}

// @desc    Update bootcamp
// @route   PUT /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) UpdateBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(bc.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("publisher", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	err := Bootcamp.FindId(bson.ObjectIdHex(id)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no bootcamp with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	if bootcamp.Deleted {
		utils.ErrorResponse(w, http.StatusNotFound, errors.New("this bootcamp was deleted"))
		return
	}

	if bootcamp.User != cUser.Id && cUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("you do not have permission"))
		return
	}

	var d map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("bad data"))
	}

	// The Update method is incompleted so the error is not handled
	// see https://github.com/zebresel-com/mongodm/issues/20
	bootcamp.Update(d)

	if valid, issues := bootcamp.ValidateUpdate(); !valid {
		utils.ErrorResponse(w, http.StatusBadRequest, issues...)
		return
	}

	err = bootcamp.Save()
	if err != nil {
		utils.ErrorHandler(w, err)
		return
	}

	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    bootcamp,
	})

}

// @desc    Delete bootcamp
// @route   DELETE /api/v1/bootcamps/:id
// @access  Private
func (bc *Bootcamp) DeleteBootcamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cUser := getCurrentUser(bc.connection, r)
	if cUser == nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	if !cUser.IsUserInRoles("publisher", "admin") {
		utils.ErrorResponse(w, http.StatusForbidden, fmt.Errorf("user with %s role do not autorize for this route", cUser.Role))
		return
	}

	id := ps.ByName("id")
	if !bson.IsObjectIdHex(id) {
		utils.ErrorResponse(w, http.StatusBadRequest, errors.New("invalid bootcamp id format"))
		return
	}

	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamp := &models.Bootcamp{}

	err := Bootcamp.FindId(bson.ObjectIdHex(id)).Exec(bootcamp)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no bootcamp with id of %s", id))
		return
	} else if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, errors.New("server error"))
		return
	}

	if bootcamp.Deleted {
		// should not found the deleted bootcamp
		utils.ErrorResponse(w, http.StatusNotFound, fmt.Errorf("no bootcamp with id of %s", id))
		return
	}

	if bootcamp.User != cUser.Id && cUser.Role != "admin" {
		utils.ErrorResponse(w, http.StatusForbidden, errors.New("you do not have permission"))
		return
	}

	bootcamp.SetDeleted(true)
	bootcamp.Save()
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    nil,
	})
}

// @desc    Get bootcamps within radius
// @route   GET /api/v1/bootcamps/radius/:zipcode/:distance
// @access  Public
func (bc *Bootcamp) GetBootcampsInRadius(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	zipcode := ps.ByName("zipcode")
	distance, err := strconv.ParseFloat(ps.ByName("distance"), 64)
	if err != nil {
		utils.SendJSON(w, http.StatusBadRequest, errors.New("distance should be number"))
	}
	fmt.Printf("zipcode: %s, distance: %v\n", zipcode, distance)

	// get lat/lng from mapquest usint zipcode
	// TO DO
	lng := -71.1
	lat := 42.34

	// Calc radius using radians
	// Divide dist by radius of earth
	// earth radius = 3,963mi (6,378km)
	radius := distance / 3963.0

	Bootcamp := bc.connection.Model("Bootcamp")
	bootcamps := []*models.Bootcamp{}

	query := bson.M{
		"location": bson.M{
			"$geoWithin": bson.M{
				"$centerSphere": []interface{}{
					[]float64{
						lng,
						lat,
					},
					radius,
				},
			},
		},
	}

	Bootcamp.Find(query).Exec(&bootcamps)
	utils.SendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"count":   len(bootcamps),
		"data":    bootcamps,
	})

}
