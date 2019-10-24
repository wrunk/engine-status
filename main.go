package main

import (
	"fmt"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func init() {
	// Turn on some spew options to avoid certain death
	spew.Config.DisableMethods = true
}

func main() {
	if len(os.Args) != 3 {
		helpExit()
	}
	url := os.Args[1]
	key := os.Args[2]
	maxHTTPReqs := 100
	totalHTTPReqs := 0
	sleepTime := 0
	es := NewES()

	if !serverIsCompatible(url, key, sleepTime) {
		fatal("Server URL isn't compatible")
	}

	// var cancelFunc func()
	// ctx, cancelFunc := context.WithCancel(context.Background())

	// _ = cancelFunc
	// Render loop. Http processing will happen in its own loop
	for true {
		er, _ := pingServer(url, key, sleepTime)
		totalHTTPReqs++
		es.AddResp(er)
		fmt.Println(er)
		fmt.Println("Sleeping for 500 milli")
		es.RenderTable()
		time.Sleep(time.Millisecond * 500)
		if totalHTTPReqs > maxHTTPReqs {
			fmt.Printf("Hit %d requests, exiting\n", maxHTTPReqs)
			break
		}
	}
	// cancelFunc()
}

func helpExit() {
	fmt.Printf("Please run the program like:\n")
	fmt.Printf("go run *.go <url> <key>\n")
	fmt.Println("Where key is tied to your server endpoint")
	os.Exit(-1)
}
