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
func pingServer() (*GAEInstance, error) {
	req, err := http.NewRequest(http.MethodGet, url+"?k="+key, nil)
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
	if resp.StatusCode != 200 {
		log("Error: Didn't get back 200, got (%d)", resp.StatusCode)
		return nil, fmt.Errorf("Error: Didn't get back 200, got (%d)",
			resp.StatusCode)
	}

	gi := &GAEInstance{}
	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("Error: Failed to read response body %v", err)
		return nil, err
	}
	err = json.Unmarshal(bys, gi)
	if err != nil {
		log("Error: Failed to unmarshal json %v", err)
		return nil, err
	}
	return gi, nil
}

func serverIsCompatible() bool {
	gi, err := pingServer()
	if err != nil || gi.GAEInstance == "" || gi.GAEVersion == "" {
		fmt.Printf("Invalid resp from server: err: %s gi: %#v", err, gi)
		return false
	}
	return true
}
