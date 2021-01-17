package inttest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"tagallery.com/api/config"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// apiURL takes a route and returns the full API url.
func apiURL(route string) string {
	return fmt.Sprintf("http://localhost:%v%v", config.Get().Port, route)
}

// GetRequest sends a HTTP GET request and parses the returned data into the type of {response}.
func GetRequest(url string, response interface{}) error {
	return Request("GET", url, nil, response)
}

// PostRequest sends a HTTP POST request and parses the returned data into the type of {response}.
func PostRequest(url string, body interface{}, response interface{}) error {
	return Request("POST", url, body, response)
}

// PostRequest sends a HTTP DELETE request and parses the returned data into the type of {response}.
func DeleteRequest(url string, response interface{}) error {
	return Request("DELETE", url, nil, response)
}

// Request sends a HTTP request to {url} and parses the returned data into the type of {response}.
func Request(method string, url string, body interface{}, response interface{}) error {
	var resp *http.Response

	content, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(content))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(respBody, response)
}
