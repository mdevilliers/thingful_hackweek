package main

import "strings"

type classifier struct {
	d map[string]bool
}

// NewWeatherClassifier returns the simplest text/weather classifier I can think of.
func NewWeatherClassifier() *classifier {

	keywords := []string{
		"hot", "cold", "warm", "cool",
		"frosty", "frost",
		"beautiful",
		"day", "night",
		"outside",
		"sunny", "sun", "sunshine", "sunset", "sunrise", "sun-rise",
		"bright", "dark",
		"humid",
		"icy", "ice",
		"morning",
		"nice",
		"cloud", "cloudy",
		"freezing", "freeze",
		"summer", "autumn", "winter", "spring",
		"shining",
		"degrees",
		"temp", "temperature",
		"gorgeous",
		"snow", "snowy",
		"hail",
		"wind",
		"rain", "rainy",
		"frost",
		"leaves",
		"bluesky", "blue-sky",
		"weather", "lovelyweather"}

	d := map[string]bool{}

	for _, w := range keywords {
		d[w] = true
	}

	return &classifier{
		d: d,
	}
}

func (c *classifier) IsWeather(text string) (bool, int) {

	clean := strings.Split(cleanString(text), " ")

	score := 0
	for _, w := range clean {
		if c.d[strings.ToLower(w)] == true {
			score++
		}
	}

	return score > 2, score
}

// string cleaning functions
const delim = "?!.;,*#@"

func isDelim(c string) bool {
	if strings.Contains(delim, c) {
		return true
	}
	return false
}

func cleanString(input string) string {

	size := len(input)
	temp := ""
	var prevChar string

	for i := 0; i < size; i++ {
		str := string(input[i])
		if (str == " " && prevChar != " ") || !isDelim(str) {
			temp += str
			prevChar = str
		} else if prevChar != " " && isDelim(str) {
			temp += " "
		}
	}
	return temp
}
