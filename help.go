package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getSha1Hash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func log(f string, args ...interface{}) {
	conTex.Lock()
	consoleMsgQ = append(consoleMsgQ, fmt.Sprintf(f, args...))
	conTex.Unlock()
}

func sendHTTP() (*GAEInstance, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
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
