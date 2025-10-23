package main

import (
	"fmt"
	"io"
	"net/http"
)

const baseURL string = "https://aviationweather.gov/api/data"

func getAirportID() string {
	return "KCGI"
}

func fetchMETAR(fieldID string) string {
	metarURL := fmt.Sprintf("%s/metar?ids=%s&format=raw", baseURL, fieldID)

	response, e := http.Get(metarURL)
	if e != nil {
		fmt.Println("Error during fetch: ", e)
		return ""
	}
	defer response.Body.Close()

	byteMETAR, e := io.ReadAll(response.Body)
	if e != nil {
		fmt.Println("Error reading response: ", e)
	}

	return string(byteMETAR)
}

func main() {
	ID := getAirportID()
	METAR := fetchMETAR(ID)

	fmt.Println(METAR)
}
