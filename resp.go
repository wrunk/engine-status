package main

import "github.com/davecgh/go-spew/spew"

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

func (er *EngineResponse) String() string {
	return spew.Sdump(er) // Cheating
}
