package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var (
	url           = ""
	maxHTTPReqs   = 0
	totalHTTPReqs = 0
	// How many active go routines we have calling the server
	curHTTPWorkers = 0
	// Max go routines we can have at any one time
	maxHTTPWorkers = 5

	// Control the outer loop
	done = false

	serverSleepTime = 0
)

// Http orchestration loop
// If we have less than maxCon go routines out,
// spawn more aggressively
// Otherwise wait on a synchronous channel receive
func orcFetch(ctx context.Context) {

	wg := &sync.WaitGroup{}
	resChan := make(chan SvrResp)
	for !done {
		// If we're at max con, sleep for 50ms

		select {
		// If our ctx is canceled, bail
		case <-ctx.Done():
			fmt.Println("orcFetch found context done signal!, going to bail")
		case r := <-resChan:
			curHTTPWorkers-- // This is only func dealing with these vars synchronously
			totalHTTPReqs++
			fmt.Printf("Found response %#v\n", r)
			fmt.Printf("Now at totalHTTPReqs %d\n", totalHTTPReqs)
		default:
			if totalHTTPReqs >= maxHTTPReqs {
				fmt.Printf("Now at totalHTTPReqs %d, exiting\n", totalHTTPReqs)
				done = true
				break
			}
			if curHTTPWorkers < maxHTTPWorkers {
				fmt.Println("Going to add another httpworker")
				// spawn another this loop cycle
				wg.Add(1)
				go httpWorker(wg, resChan)
				curHTTPWorkers++ // This is only func dealing with these vars synchronously
			} else {
				fmt.Println("Select Stmt got bored, waiting for a bit")
				time.Sleep(800 * time.Millisecond)
			}

		}
	}
	wg.Wait()

}

// TODO refactor this once we know how this will all work. Probably
// dont need this and sendHTTP
func httpWorker(wg *sync.WaitGroup, resChan chan SvrResp) {
	// Remove a waitgroup item in case we are shutting down
	// Always do this in the beginning in case this function changes
	// and we forget
	defer wg.Done()
	gaeIns, code, err := sendHTTP()
	// Send the response for processing
	resChan <- SvrResp{*gaeIns, code, err}
}

// Wrap the server response since we always need to return something
// even if there was a failure
type SvrResp struct {
	GAEI   GAEInstance
	Status int
	Err    error
}

type GAEInstance struct {
	GAEApplication  string `json:"gae_application"`
	GAEDeploymentID string `json:"gae_deployment_id"`
	GAEEnv          string `json:"gae_env"`
	GAEInstance     string `json:"gae_instance"`
	GAEMBMem        string `json:"gae_memory_mb"`
	GAEService      string `json:"gae_service"`
	GAERuntime      string `json:"gae_runtime"`
	GAEVersion      string `json:"gae_version"`
	GCP             string `json:"google_cloud_project"`
	Port            string `json:"port"`
}

func sendHTTP() (*GAEInstance, int, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		// log("Error: Failed to create http request %v", err)
		return nil, 0, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// log("Error: Failed to make http request %v", err)
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log("Error: Didn't get back 200, got (%d)", resp.StatusCode)
		return nil, resp.StatusCode, fmt.Errorf("Error: Didn't get back 200, got (%d)",
			resp.StatusCode)
	}

	gi := &GAEInstance{}
	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("Error: Failed to read response body %v", err)
		return nil, resp.StatusCode, err
	}
	err = json.Unmarshal(bys, gi)
	if err != nil {
		log("Error: Failed to unmarshal json %v", err)
		return nil, resp.StatusCode, err
	}
	return gi, resp.StatusCode, nil
}

// Don't really need this...
// But could be used to dedupe instances in case there were any concerns
func (g *GAEInstance) Hash() string {
	t := reflect.TypeOf(*g)
	v := reflect.ValueOf(*g)
	s := ""

	for i := 0; i < t.NumField(); i++ {
		s = fmt.Sprintf("%s:%v", s, v.Field(i))
		// The below will print the type and value
		// fmt.Printf("(%v) %+v\n", v.Field(i), t.Field(i))
	}
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
