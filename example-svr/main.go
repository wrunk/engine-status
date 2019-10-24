package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
This folder contains an example App Engine app.

It is a working GAE2.0 Go1.13 example current as of October 2019

For this example I just route all http paths to the engineStatus
function/endpoint. In real go webapps you'll want to use a real
router like the mux from https://www.gorillatoolkit.org/

You can easily deploy this to app engine by:

gcloud app deploy --project=<your-gcp-proj-id> ./app.yaml

You will need to setup a gcp project, enable app engine, and
setup a billing account although this initial deploy and minimal test
should be free
*/

const (
	secretKey = "HoglBoglBooBooJay"
)

var (
	errJSON   = `{"error":"%s"}`
	goProcSig = RandomString(40)
)

func newER() *EngineResponse {
	return &EngineResponse{
		GAEMemoryMB:     os.Getenv("GAE_APPLICATION"),
		CGOEnabled:      os.Getenv("CGO_ENABLED"),
		GAEInstance:     os.Getenv("GAE_INSTANCE"),
		Home:            os.Getenv("HOME"),
		Port:            os.Getenv("PORT"),
		Goroot:          os.Getenv("GOROOT"),
		GAEService:      os.Getenv("GAE_SERVICE"),
		Path:            os.Getenv("PATH"),
		GAEDeploymentID: os.Getenv("GAE_DEPLOYMENT_ID"),
		DebianFrontend:  os.Getenv("DEBIAN_FRONTEND"),
		GCP:             os.Getenv("GOOGLE_CLOUD_PROJECT"),
		GAEEnv:          os.Getenv("GAE_ENV"),
		PWD:             os.Getenv("PWD"),
		GAEApplication:  os.Getenv("GAE_APPLICATION"),
		GAERuntime:      os.Getenv("GAE_RUNTIME"),
		GAEVersion:      os.Getenv("GAE_VERSION"),

		GoProcSig: goProcSig,
	}
}

func main() {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	http.HandleFunc("/", engineStatus)
	fmt.Println("Going to start test server on port: " + port)
	for _, line := range os.Environ() {
		s := strings.Split(line, "=")
		fmt.Printf("[%s] ====> [%s]\n", s[0], s[1])
	}
	http.ListenAndServe(":"+port, nil)
}

func engineStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("**************HEADERS:::::::::::::::::::::;")
	for k, v := range r.Header {
		fmt.Printf("[%s][%#v]\n", k, v)
	}
	sleepTimeMS := r.URL.Query().Get("s")
	key := r.URL.Query().Get("k")
	if key != secretKey {
		writeJSONError(w, "Invalid key k")
		return
	}

	st, err := strconv.Atoi(sleepTimeMS)
	if err == nil {
		fmt.Printf("Sleeping for (%d) MS\n", st)
		time.Sleep(time.Duration(st))
	} else {
		writeJSONError(w, "Sleep time s was invalid")
		return
	}
	writeJSON(w, newER())
}

func writeJSON(w http.ResponseWriter, obj interface{}) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		writeJSONError(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func writeJSONError(w http.ResponseWriter, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	// Assuming this will write the full header back to the client so all
	// headers have to be written and can't be afterward
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, fmt.Sprintf(errJSON, errMsg))
}
