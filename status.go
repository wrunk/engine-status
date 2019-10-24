package main

import (
	"fmt"
	"time"

	tm "github.com/buger/goterm"
)

type EngineStatus struct {
	TimerStart time.Time // For use with finding elapsed time
	Instances  map[string]*InstanceMeta
}

type InstanceMeta struct {
	Added       time.Time
	NumRequests int
	Key         string // AKA go proc sig
	FirstResp   EngineResponse
}

func (es *EngineStatus) NumInstances() int {
	return len(es.Instances)
}

// Returns true if we've added a new instance
func (es *EngineStatus) AddResp(er *EngineResponse) bool {
	if _, ok := es.Instances[er.GoProcSig]; ok {
		im := es.Instances[er.GoProcSig]
		im.NumRequests += 1
		return false
	} else {
		es.Instances[er.GoProcSig] = &InstanceMeta{
			Added:       time.Now(),
			NumRequests: 1,
			Key:         er.GoProcSig,
			FirstResp:   *er,
		}
		return true
	}
}

func NewES() *EngineStatus {
	return &EngineStatus{
		TimerStart: time.Now(),
		Instances:  map[string]*InstanceMeta{},
	}
}

func (es *EngineStatus) RenderTable() {
	// tm.Clear() // Clear current screen
	totals := tm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprintf(totals, "Timer\tTotal Instances\tTotal Reqs\n")
	tr := 0
	for _, i := range es.Instances {
		tr += i.NumRequests
	}
	elapsed := time.Now().Sub(es.TimerStart)
	fmt.Fprintf(totals, "%s\t%d\t%d\n", elapsed, len(es.Instances), tr)
	tm.Println(totals)
	tm.Flush()
}
