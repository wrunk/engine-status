package main

// import (
// 	"context"
// 	"fmt"

// 	term "github.com/nsf/termbox-go"
// )

// func syncInput(ctx context.Context, cf context.CancelFunc) {
// 	defer term.Close()
// 	err := term.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	for true {
// 		if ev := term.PollEvent(); ev.Type == term.EventKey {
// 			fmt.Printf("maxHTTP %d sleep %d\n", maxHTTPWorkers, serverSleepTime)
// 			if ev.Key == term.KeyEsc {
// 				return
// 			} else if ev.Ch == 'c' {
// 				maxHTTPWorkers += 1
// 			} else if ev.Ch == 'x' && maxHTTPWorkers > 1 {
// 				maxHTTPWorkers -= 1
// 			} else if ev.Ch == 's' {
// 				serverSleepTime += 100
// 			} else if ev.Ch == 'a' && serverSleepTime > 0 {
// 				serverSleepTime -= 100
// 			} else {

// 			}
// 		}
// 	}
// }
