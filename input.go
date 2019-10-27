package main

import (
	"context"
	"fmt"
	"time"

	term "github.com/nsf/termbox-go"
)

func syncInput(ctx context.Context, cf context.CancelFunc, es *EngineStatus) {
	fmt.Println("Going to sync some inputs!")
	defer term.Close()
	err := term.Init()
	if err != nil {
		panic(err)
	}
	for true {
		if ev := term.PollEvent(); ev.Type == term.EventKey {
			fmt.Printf("Conn %d Sleep Delay %d\n", es.Concurrency, es.SleepDelayMS)
			// EDSF!!
			if ev.Key == term.KeyEsc {
				// TODO make this call cancel func!
				fmt.Println("Esc pressed, exiting input goroutine")
				return
			} else if ev.Ch == 'e' {
				es.Concurrency++
				fmt.Printf("E key pressed, rockin up concurrency to %d\n",
					es.Concurrency)

			} else if ev.Ch == 'd' {
				if es.Concurrency < 2 { // Can't go lower than one
					continue
				}
				es.Concurrency--
				fmt.Printf("D key pressed, Woa tiger, slowin down. Now at %d\n",
					es.Concurrency)
			} else if ev.Ch == 'f' {
				es.SleepDelayMS += 100
				fmt.Printf("F key pressed, Pokin the server, causing delays, now at %d MS\n",
					es.SleepDelayMS)
			} else if ev.Ch == 'a' {
				if es.SleepDelayMS == 0 { // Can't go lower than zero
					continue
				}
				es.SleepDelayMS -= 100
				fmt.Printf("A key pressed, Giving the server a break, now at %d MS\n",
					es.SleepDelayMS)
			} else if ev.Ch == 't' {
				es.TimerStart = time.Now()
				fmt.Println("T key pressed, restarting timer")
			}
		}
	}
}
