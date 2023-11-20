package utils

import (
	"encoding/json"
	"net/http"
)

type IPAPIResponse struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

func GetLocationFromIP(ip string) *IPAPIResponse {
	url := "https://ipapi.co/" + ip + "/json/"

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode == http.StatusTooManyRequests {
		return nil
	}

	defer resp.Body.Close()

	var location IPAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return nil
	}

	return &location
}
