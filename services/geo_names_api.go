package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "http://api.geonames.org/searchJSON"

type geoNamesResponse struct {
	Items []geoName `json:"geonames"`
}

type geoName struct {
	Name string `json:"name"`
}

type GeoNamesApi struct {
	client *http.Client
}

func NewGeoNamesApi(client *http.Client) *GeoNamesApi {
	return &GeoNamesApi{
		client: client,
	}
}

func (gna *GeoNamesApi) CountCitiesByLetter(letter string) (count int, err error) {
	response, err := gna.get()
	if err != nil {
		return count, fmt.Errorf("failed to get GeoNames response: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return count, fmt.Errorf("GeoNames returned unexpected status code: %d", response.StatusCode)
	}

	var geonameResponse geoNamesResponse
	if err := json.NewDecoder(response.Body).Decode(&geonameResponse); err != nil {
		return count, fmt.Errorf("failed to decode GeoNames response: %w", err)
	}

	letter = strings.ToLower(letter)
	for _, item := range geonameResponse.Items {
		if strings.HasPrefix(strings.ToLower(item.Name), letter) {
			count++
		}
	}
	return count, nil
}

func (gna *GeoNamesApi) get() (response *http.Response, err error) {

	u, err := url.Parse(baseURL)
	if err != nil {
		return response, err
	}

	q := u.Query()
	q.Set("featureClass", "P")
	q.Set("maxRows", "50")
	q.Set("username", "hsample")
	q.Set("orderby", "population")
	u.RawQuery = q.Encode()

	response, err = gna.client.Get(u.String())
	return response, err
}
