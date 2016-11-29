package main

func CategoriseThing(providerType string) string {

	switch providerType {

	case "api.citybik.es", "bikes.oobrien", "marlin.casa", "globalbikeshare", "citybikes":
		return "Bike Dock"

	case "openweathermap", "metoffice", "wunderground":
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
