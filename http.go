package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Send a single, synchronous request to the server.
// With normal use, we'll send thousands of these.
// TODO use a client, request pool
func pingServer(url, key string, sleepTime int) (*EngineResponse, error) {
	fullURL := fmt.Sprintf("%s?k=%s&s=%d", url, key, sleepTime)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		log("Error: Failed to create http request %v", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log("Error: Failed to make http request %v", err)
		return nil, err
	}

	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("Error: Failed to read response body %v and http status code was %d",
			err, resp.StatusCode)
		return nil, err
	}

	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("Error: Didn't get back 200, got (%d), body (%s)",
			resp.StatusCode, string(bys))
		log(msg)
		return nil, fmt.Errorf(msg)
	}

	gi := &EngineResponse{}
	err = json.Unmarshal(bys, gi)
	if err != nil {
		log("Error: Failed to unmarshal json %v", err)
		return nil, err
	}
	return gi, nil
}

func serverIsCompatible(url, key string, sleepTime int) bool {
	gi, err := pingServer(url, key, sleepTime)
	if err != nil || gi.GAEInstance == "" || gi.GAEVersion == "" {
		fmt.Printf("Invalid resp from server: err: %s gi: %#v", err, gi)
		return false
	}
	return true
}
