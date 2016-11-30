package main

import (
	"fmt"
	"log"
)

func DistinctByLocationAndCategory(client *client, lat, long, radius float64) ([]item, error) {

	results, err := client.SearchByLocation(lat, long, radius, 500)

	if err != nil {
		return nil, err
	}

	distinct := map[string]Thing{}

	for i, _ := range results.Data {

		category := categoriseThing(results.Data[i].Relationships.Provider.Data.ID)

		_, found := distinct[category]

		if !found {
			distinct[category] = results.Data[i]
		}
	}

	items := []item{}

	for k, _ := range distinct {
		r, err := client.Access(distinct[k].ID)

		if err != nil {
			log.Print(err.Error())
			continue
		}

		insight := getInsight(k, distinct[k])

		thing := item{
			Type: "thing",
			Location: location{
				Latitude:  distinct[k].Attributes.Location.Latitude,
				Longitude: distinct[k].Attributes.Location.Longitude,
			},
			Data:       r,
			Categories: []string{k},
			Insight:    insight,
			Distance:   distinct[k].Attributes.Distance,
		}

		items = append(items, thing)

	}

	return items, nil

}

func categoriseThing(providerType string) string {

	switch providerType {

	case "api.citybik.es", "bikes.oobrien", "marlin.casa", "globalbikeshare", "citybikes":
		return "Bike Dock"

	case "openweathermap", "metoffice", "wunderground", "wowmet":
		return "Weather Station"

	case "aqicn", "netatmo":
		return "Air Quality"

	case "webcams":
		return "Webcam"

	case "chargepoints", "chargepointsuk":
		return "Electric Car Charging Dock"

	case "environment.data.gov", "gaugemap":
		return "Water Level"

	case "transportapi":
		return "Transport"

	case "thingspeak", "xively", "thethingsnetwork", "smartcitizen":
		return "Random Stuff"

	case "wikibeacon":
		return "iBeacon"

	}
	return "Unknown"
}

func getInsight(category string, thing Thing) string {

	switch category {

	case "Bike Dock":
		return fmt.Sprintf("Weather good enough for a cycle? Nearest cycle dock is %f meters away?", thing.Attributes.Distance)

	case "Transport":
		return fmt.Sprintf("Nice enough to leave the car behind? How about getting public transport from %s?", thing.Attributes.Title)
	}

	return "No insight gained."

}
