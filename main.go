package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"

	tm "github.com/buger/goterm"
)

var (
	table          = tm.NewTable(0, 10, 5, ' ', 0)
	consoleMsgQ    = []string{}
	conTex         = &sync.Mutex{}
	instances      = map[string]GAEInstance{}
	instanceIDs    = []string{}
	insTex         = &sync.Mutex{}
	url            = ""
	maxHTTPReqs    = 0
	totalHTTPReqs  = 0
	done           = false
	frameSleepTime = 2
	wg             = &sync.WaitGroup{}
)

/*
"gae_application":   os.Getenv("GAE_APPLICATION"),
"gae_deployment_id": os.Getenv("GAE_DEPLOYMENT_ID"),
"gae_env":           os.Getenv("GAE_ENV"),
"gae_instance":      os.Getenv("GAE_INSTANCE"),
"gae_memory_mb":     os.Getenv("GAE_MEMORY_MB"),
"gae_service":       os.Getenv("GAE_SERVICE"),
"gae_runtime":       os.Getenv("GAE_RUNTIME"),
"gae_version":       os.Getenv("GAE_VERSION"),

"google_cloud_project": os.Getenv("GOOGLE_CLOUD_PROJECT"),
// Only used with new, second generation runtimes
"port": os.Getenv("PORT"),
*/
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

// Don't really need this...
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

func helpExit() {
	fmt.Printf("Please run the program like:\n")
	fmt.Printf("go run *.go <url> <max-requests>\n")
	os.Exit(-1)
}
func main() {

	if len(os.Args) != 3 {
		helpExit()
	}
	url = os.Args[1]
	maxH := os.Args[2]
	maxHTTPReqs, err := strconv.Atoi(maxH)
	if err != nil || maxHTTPReqs < 1 {

		fmt.Printf("Invalid max requests %s\n", maxH)
		helpExit()

	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	go fetcho(ctx)

	// Render loop. Http processing will happen in its own loop
	for true {
		gi, _ := sendHTTP()
		addGI(gi)
		render()
		time.Sleep(time.Second * 3)
		if totalHTTPReqs > maxHTTPReqs {
			fmt.Printf("Hit %d requests, exiting\n", maxHTTPReqs)
			break
		}
	}
	cancelFunc()
}

func addGI(gi *GAEInstance) {
	/*
			sl := []string{"mumbai", "london", "tokyo", "seattle"}
		    sort.Sort(sort.StringSlice(sl))
		    fmt.Println(sl)
	*/
	if _, found := instances[gi.GAEInstance]; !found {
		fmt.Println("Adding gi, need lok")
		insTex.Lock()
		fmt.Println("Adding gi")
		instances[gi.GAEInstance] = *gi
		instanceIDs = append(instanceIDs, gi.GAEInstance)
		sort.Sort(sort.StringSlice(instanceIDs))
		insTex.Unlock()
	}
}

// Http orchestration loop
func fetcho(ctx context.Context) {

	wg := &sync.WaitGroup{}

	for true {
		// If we're at max con, sleep for 50ms

		// If our ctx is canceled, bail
	}

}

func httpWorker() {}

func render() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	table = tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(table, "InstanceID\tVersion\tProject\tService\tRuntime\n")
	for _, id := range instanceIDs {
		i := instances[id]
		fmt.Fprintf(table, "...%s\t%s\t%s\t%s\t%s\t\n",
			i.GAEInstance[len(i.GAEInstance)-10:],
			i.GAEVersion,
			i.GCP,
			i.GAEService,
			i.GAERuntime,
		)

	}

	tm.Println("=============================================================================")
	tm.Println("Active App Engine Instances")
	tm.Println("=============================================================================")
	tm.Println(table)
	tm.Println("=============================================================================")
	tm.Println("Concurrency: [5] || Server Sleep (MS) [0]")
	tm.Println("Press up/down to change concurrency || Press s to increase sleep 100ms")
	tm.Println("d to decrease 100ms")
	tm.Println("=============================================================================")
	if len(consoleMsgQ) > 0 {
		conTex.Lock()
		ss := consoleMsgQ[0]
		consoleMsgQ = consoleMsgQ[1:]
		tm.Println(ss)
		conTex.Unlock()
	} else {
		tm.Println("No messages...")

	}
	tm.Println("=============================================================================")

	tm.Flush()
}
