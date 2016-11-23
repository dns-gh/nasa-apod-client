package nasaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	nasaApodAPIGet = "https://api.nasa.gov/planetary/apod?api_key="

	// The archive page of a picture of a particular day can be retrieve via the following info.
	nasaApodArchivePage          = "https://apod.nasa.gov/apod/ap"
	nasaApodArchivePageExtension = ".html"

	nasaApodServiceVersion = "v1"
	nasaAPIDefaultKey      = "DEMO_KEY"
	nasaTimeFormat         = "2006-01-02"
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

func makeArchiveDate(date string) string {
	return nasaApodArchivePage + strings.Replace(date[2:], "-", "", 2) + nasaApodArchivePageExtension
}

// FetchAPOD fetches metadata information of the image of the given 'date'
// in 'hd' format potentially.
// The 'date' must be in the following format: 'YYYY-MM-DD'
func (n *NasaApodClient) FetchAPOD(date string, hd bool) (*Apod, string, error) {
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
		return nil, "", err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if strings.Contains(string(bytes), "OVER_RATE_LIMIT") {
		return nil, "", fmt.Errorf("http get rate limit reached, wait or use a proper key instead of the default one")
	}
	apod := &Apod{}
	json.Unmarshal(bytes, apod)
	if apod.ServiceVersion != nasaApodServiceVersion {
		log.Printf("[WARNING] remote service version %s != %s\n", apod.ServiceVersion, nasaApodServiceVersion)
	}
	return apod, makeArchiveDate(apod.Date), nil
}

// FetchTodayAPOD fetches metadata information of today's Apod image
// in 'hd' format potentially.
func (n *NasaApodClient) FetchTodayAPOD(hd bool) (*Apod, string, error) {
	return n.FetchAPOD("", hd)
}

func loadImage(uri string) (string, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error request status: %s != 200", resp.Status)
	}
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FetchHD fetches the apod today's image and description.
func (n *NasaApodClient) FetchHD() (string, string, string, error) {
	apod, archiveURL, err := n.FetchTodayAPOD(true)
	if err != nil {
		return "", "", "", err
	}
	img, err := loadImage(apod.URL)
	if err != nil {
		return "", "", "", err
	}
	return apod.Explanation, img, archiveURL, nil
}
