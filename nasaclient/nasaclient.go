package nasaclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	nasaApodAPIGet    = "https://api.nasa.gov/planetary/apod?api_key="
	nasaAPIDefaultKey = "DEMO_KEY"
	nasaTimeFormat    = "2006-01-02"
)

// NasaApodClient represents the web Client.
type NasaApodClient struct {
	apiKey string
}

// MakeNasaApodClient creates a web client to make http request
// to the Neo Nasa API: https://api.nasa.gov/api.html#NeoWS
func MakeNasaApodClient() *NasaApodClient {
	log.Println("[nasa-apod] making nasa apod client")
	apiKey := os.Getenv("NASA_API_KEY")
	if len(apiKey) == 0 {
		apiKey = nasaAPIDefaultKey
	}
	return &NasaApodClient{
		apiKey: apiKey,
	}
}

// Apod represents the metadata information of an image fetched via the
// Astronomy Picture Of the Day (APOD) API: https://api.nasa.gov/api.html#apod
type Apod struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// FetchAPOD fetches metadata information of the image of the given 'date'
// in 'hd' format potentially.
// The 'date' must be in the following format: 'YYYY-MM-DD'
func (n *NasaApodClient) FetchAPOD(date string, hd bool) (*Apod, error) {
	url := nasaApodAPIGet +
		n.apiKey
	if len(date) != 0 {
		url += "&date=" + date
	}
	if hd {
		url += "&hd=true"
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if strings.Contains(string(bytes), "OVER_RATE_LIMIT") {
		return nil, fmt.Errorf("http get rate limit reached, wait or use a proper key instead of the default one")
	}
	apod := &Apod{}
	json.Unmarshal(bytes, apod)
	return apod, nil
}

// FetchTodayAPOD fetches metadata information of today's Apod image
// in 'hd' format potentially.
func (n *NasaApodClient) FetchTodayAPOD(hd bool) (*Apod, error) {
	return n.FetchAPOD("", hd)
}
