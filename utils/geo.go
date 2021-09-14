package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Location struct {
	Results []struct {
		Locations []struct {
			Street  string `json:"street"`
			City    string `json:"adminArea5"`
			State   string `json:"adminArea3"`
			Country string `json:"adminArea1"`
			Zipcode string `json:"postalCode"`
			LatLng  struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"locations"`
	} `json:"results"`
}

func GetLocation(add string) *Location {
	b := os.Getenv("GEOCODER_URL")
	k := os.Getenv("GEOCODER_API_KEY")
	ss := strings.Split(add, " ")
	var addTrim string
	for _, v := range ss {
		if v != " " {
			addTrim += v + "%20"
		}
	}
	uri := fmt.Sprintf("%s?key=%s&location=%s", b, k, addTrim)
	resp, err := http.Get(uri)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bs, _ := io.ReadAll(resp.Body)

	var loc Location
	err = json.Unmarshal(bs, &loc)
	if err != nil {
		log.Fatal(err)
	}
	return &loc

}
