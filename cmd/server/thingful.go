package main

import (
	"encoding/json"
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

	return c.doGetRequest(url)
}

func (c *client) SearchByLocation(lat, long, radius float64, limit int) (*Results, error) {

	template := "https://api.thingful.net/search?geo-lat=%f&geo-long=%f&geo-radius=%f&limit=%d&sort=distance"
	url := fmt.Sprintf(template, lat, long, radius, limit)

	return c.doGetRequest(url)
}

func (c *client) DistinctByLocationAndCategory(lat, long, radius float64) (*Results, error) {

	results, err := c.SearchByLocation(lat, long, radius, 500)

	if err != nil {
		return nil, err
	}

	distinct := map[string]Thing{}

	for i, _ := range results.Data {

		category := CategoriseThing(results.Data[i].Relationships.Provider.Data.ID)

		_, found := distinct[category]

		if !found {
			distinct[category] = results.Data[i]
		}
	}

	toReturn := &Results{
		Data: []Thing{},
	}

	for k, _ := range distinct {
		toReturn.Data = append(toReturn.Data, distinct[k])

	}

	return toReturn, nil

}

func (c *client) doGetRequest(url string) (*Results, error) {

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
		return nil, fmt.Errorf("error loading %s, statusCode :%d ", url, res.StatusCode)
	}

	defer res.Body.Close()

	jsonBytes, err := ioutil.ReadAll(res.Body)

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

type Thing struct {
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
}

type Results struct {
	Data  []Thing `json:"data"`
	Links `json:"links"`
}

type Links struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
	Self string `json:"self"`
}
