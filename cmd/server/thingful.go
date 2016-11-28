package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient(apiKey string) *client {

	return &client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}

}

func (c *client) Access(thinguid string) (*Results, error) {

	template := "https://api.thingful.net/access?uid=%s"
	url := fmt.Sprintf(template, thinguid)

	jsonBytes, err := c.doGetRequest(url)

	if err != nil {
		return nil, err
	}

	results := Results{}
	fmt.Println(string(jsonBytes))
	err = json.Unmarshal(jsonBytes, &results)

	if err != nil {
		return nil, err
	}

	return &results, nil

}

func (c *client) SearchByLocation(lat, long, radius float64) (*Results, error) {

	template := "https://api.thingful.net/search?geo-lat=%f&geo-long=%f&geo-radius=%f&limit=10&sort=distance"
	url := fmt.Sprintf(template, lat, long, radius)

	jsonBytes, err := c.doGetRequest(url)

	if err != nil {
		return nil, err
	}

	results := Results{}

	err = json.Unmarshal(jsonBytes, &results)

	if err != nil {
		return nil, err
	}

	return &results, nil

}

func (c *client) doGetRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("StatusCode :%d ", res.StatusCode))
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)

}

type Results struct {
	Links struct {
		Next string `json:"next"`
		Prev string `json:"prev"`
		Self string `json:"self"`
	} `json:"links"`
	Data []struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Links struct {
			Access string `json:"access"`
			Self   string `json:"self"`
		} `json:"links"`
		Attributes struct {
			Channels []struct {
				ID    string `json:"id"`
				Value string `json:"value"`
				Units string `json:"units"`
			} `json:"channels"`
			CreatedAt   time.Time `json:"created_at"`
			Datasource  string    `json:"datasource"`
			Description string    `json:"description"`
			Distance    float64   `json:"distance"`
			IndexedAt   time.Time `json:"indexed_at"`
			Location    struct {
				Longitude float64 `json:"longitude"`
				Latitude  float64 `json:"latitude"`
			} `json:"location"`
			Score      float64   `json:"score"`
			Title      string    `json:"title"`
			UpdatedAt  time.Time `json:"updated_at"`
			Visibility string    `json:"visibility"`
			Webpage    string    `json:"webpage"`
		} `json:"attributes"`
		Relationships struct {
			Category struct {
				Data struct {
					Type string `json:"type"`
					ID   string `json:"id"`
				} `json:"data"`
			} `json:"category"`
			Provider struct {
				Data struct {
					Type  string `json:"type"`
					ID    string `json:"id"`
					Links struct {
						Related string `json:"related"`
					} `json:"links"`
				} `json:"data"`
			} `json:"provider"`
		} `json:"relationships"`
	} `json:"data"`
}
