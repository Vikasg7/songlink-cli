package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SonglinkResponse struct {
	PageURL         string          `json:"pageUrl"`
	LinksByPlatform LinksByPlatform `json:"linksByPlatform"`
}

type LinksByPlatform struct {
	Spotify PlatformMusic `json:"spotify"`
}

type PlatformMusic struct {
	URL string `json:"url"`
}

func GetLinks(searchURL string) (string, error) {
	response, err := makeRequest(searchURL); if err != nil {
		return "", err
	}

	platform := PlatformMusic{
		URL: "",
	}
	links := LinksByPlatform{
		Spotify: platform,
	}
	linksResponse := SonglinkResponse{
		PageURL:         "",
		LinksByPlatform: links,
	}

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&linksResponse); if err != nil {
		return "", fmt.Errorf("error decoding JSON response: %w", err)
	}

	nonLocalURL := strings.ReplaceAll(linksResponse.PageURL, "/fi", "")
	spotifyURL := linksResponse.LinksByPlatform.Spotify.URL
	outputString := nonLocalURL

	if nonLocalURL != "" && spotifyURL != "" {
		outputString = fmt.Sprintf("<%s>\n%s", nonLocalURL, spotifyURL)
	}
	return outputString, nil
}

func makeRequest(searchURL string) (*http.Response, error) {
	url := buildURL(searchURL)
	response, err := http.Get(url); if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP response status: %s", response.Status)
	}

	return response, nil
}

func buildURL(searchURL string) string {
	url := url.URL{
		Scheme: "https",
		Host:   "api.song.link",
		Path:   "/v1-alpha.1/links",
	}
	values := url.Query()
	values.Add("url", searchURL)
	url.RawQuery = values.Encode()
	return url.String()
}
