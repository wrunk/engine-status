package main

import (
	"math/rand"
	"time"
)

// Just copy needed stuff in here to avoid making an importable sub package

type EngineResponse struct {
	// Fields from environment variables
	GAEMemoryMB     string `json:"gae_memory_mb"`        // GAE_MEMORY_MB
	CGOEnabled      string `json:"cgo_enabled"`          // CGO_ENABLED
	GAEInstance     string `json:"gae_instance"`         // GAE_INSTANCE
	Home            string `json:"home"`                 // HOME
	Port            string `json:"port"`                 // PORT
	Goroot          string `json:"go_root"`              // GOROOT
	GAEService      string `json:"gae_service"`          // GAE_SERVICE
	Path            string `json:"path"`                 // PATH
	GAEDeploymentID string `json:"gae_deployment_id"`    // GAE_DEPLOYMENT_ID
	DebianFrontend  string `json:"debian_frontend"`      // DEBIAN_FRONTEND
	GCP             string `json:"google_cloud_project"` // GOOGLE_CLOUD_PROJECT
	GAEEnv          string `json:"gae_env"`              // GAE_ENV
	PWD             string `json:"pwd"`                  // PWD
	GAEApplication  string `json:"gae_application"`      // GAE_APPLICATION
	GAERuntime      string `json:"gae_runtime"`          // GAE_RUNTIME
	GAEVersion      string `json:"gae_version"`          // GAE_VERSION

	// A unique stamp given to the go process running on the server so we
	// know 100% this is a unique instance
	GoProcSig string `json:"go_proc_sig"`
}

var ranStrSetAlphaNum = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func init() {
	// Seed the random num gen, once for app
	rand.Seed(time.Now().UTC().UnixNano())
}

func RandomInt(min, max int) int {
	/* We will give back min - max inclusive.
	   0,0 always gives back 0
	   1,1 always gives back 1
	   1,2 gives back either 1 or 2
	   etc
	   1,0 panics (Intn panics on negative #)
	*/
	return min + rand.Intn(max+1-min)
}

func RandomString(length uint) string {
	if length < 1 {
		msg := "Cant ask for random string of length less than 1"
		panic(msg)
	}
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = ranStrSetAlphaNum[rand.Intn(len(ranStrSetAlphaNum))]
	}
	return string(bytes)
}
