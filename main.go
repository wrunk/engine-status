package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"

	tm "github.com/buger/goterm"
)

var (
	table       = tm.NewTable(0, 10, 5, ' ', 0)
	consoleMsgQ = []string{}
	conTex      = &sync.Mutex{}
	instances   = map[string]GAEInstance{}
	instanceIDs = []string{}
	insTex      = &sync.Mutex{}
)

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
	var err error
	maxHTTPReqs, err = strconv.Atoi(maxH)
	if err != nil || maxHTTPReqs < 1 {

		fmt.Printf("Invalid max requests %s\n", maxH)
		helpExit()

	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Http orchestration. Will keep fetching until we are quiting or it has
	// fetched more than maxHTTPReqs reqs
	go orcFetch(ctx)
	syncInput(ctx, cancelFunc)

	_ = cancelFunc
	// Render loop. Http processing will happen in its own loop
	// for true {
	// 	gi, _ := sendHTTP()
	// 	addGI(gi)
	// 	render()
	// 	time.Sleep(time.Second * 3)
	// 	if totalHTTPReqs > maxHTTPReqs {
	// 		fmt.Printf("Hit %d requests, exiting\n", maxHTTPReqs)
	// 		break
	// 	}
	// }
	// cancelFunc()
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
