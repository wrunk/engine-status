package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	l "log"
	"os"
)

func getSha1Hash(text string) string {
	h := sha1.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func log(f string, args ...interface{}) {
	l.Printf(f, args...)
	// conTex.Lock()
	// consoleMsgQ = append(consoleMsgQ, fmt.Sprintf(f, args...))
	// conTex.Unlock()
}

func fatal(f string, args ...interface{}) {
	fmt.Printf(f+"\n", args...)
	os.Exit(-1)
}
