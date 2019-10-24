package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"sort"
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
	key            string
	wg             = &sync.WaitGroup{}
	ctx            context.Context
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
	GAEInstance   string `json:"gae_instance"`
	GAEService    string `json:"gae_service"`
	GAEServer     string `json:"gae_server"`
	GAEVersion    string `json:"gae_version"`
	GCP           string `json:"google_cloud_project"`
	GAEDatacenter string `json:"gae_datacenter"`

	// Unused stuff, possibly supported for second gen runtimes
	GAEEnv          string `json:"gae_env"`
	GAEApplication  string `json:"gae_application"`
	GAEDeploymentID string `json:"gae_deployment_id"`
	GAEMBMem        string `json:"gae_memory_mb"`
	GAERuntime      string `json:"gae_runtime"`
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
	fmt.Printf("go run *.go <url> <key>\n")
	fmt.Println("Where key is tied to your server endpoint")
	os.Exit(-1)
}
func main() {

	if len(os.Args) != 3 {
		helpExit()
	}
	url = os.Args[1]
	key = os.Args[2]

	if !serverIsCompatible() {
		fatal("Server URL isn't compatible")
	}

	var cancelFunc func()
	ctx, cancelFunc = context.WithCancel(context.Background())

	maxHTTPReqs = 50

	_ = cancelFunc
	// Render loop. Http processing will happen in its own loop
	for true {
		gi, _ := pingServer()
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

func render() {
	tm.Clear()
	tm.MoveCursor(0, 0)
	table = tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(table, "InstanceID\tVersion\tProject\tService\tDatacenter\n")
	for _, id := range instanceIDs {
		i := instances[id]
		fmt.Fprintf(table, "...%s\t%s\t%s\t%s\t%s\t\n",
			i.GAEInstance[len(i.GAEInstance)-10:],
			i.GAEVersion,
			i.GCP,
			i.GAEService,
			i.GAEDatacenter,
		)

	}

	tm.Println("==================================================================================================")
	tm.Println("Active App Engine Instances")
	tm.Println("==================================================================================================")
	tm.Println(table)
	tm.Println("==================================================================================================")
	tm.Println("[Concurrency: 5]  [Server Sleep (MS) 0]  [Total Requests 13]")
	tm.Println("Press up/down to change concurrency || Press s to increase sleep 100ms")
	tm.Println("d to decrease 100ms")
	tm.Println("==================================================================================================")
	if len(consoleMsgQ) > 0 {
		conTex.Lock()
		ss := consoleMsgQ[0]
		consoleMsgQ = consoleMsgQ[1:]
		tm.Println(ss)
		conTex.Unlock()
	} else {
		tm.Println("No messages...")

	}
	tm.Println("==================================================================================================")

	tm.Flush()
}
