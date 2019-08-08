// Package util contains utility functions for logging system failures
// and converting data between types []byte and float64
package util

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sync"
)

var mu = &sync.Mutex{}

// FailLog logs system failures
func FailLog(err error) {
	mu.Lock()
	f, oErr := os.OpenFile("log/fail_log.log", os.O_RDWR|os.O_APPEND, 644)
	if oErr != nil {
		log.Fatal(oErr)
	}
	defer f.Close()
	wr := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wr)
	log.Println(err)
	mu.Unlock()
	fmt.Printf("system failure logged: %v\n", err)
}

// BytesToUint64 decodes a byte slice representing sinlge int value to type uint64
func BytesToUint64(bs []byte) uint64 {
	return binary.BigEndian.Uint64(bs)
}

// Float64ToBytes encodes a uint64 value to a byte slice
func Float64ToBytes(fl float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(fl))
	return b
}

// BytesToFloat64 decodes a byte slice representing sinlge int value to type float64
func BytesToFloat64(bs []byte) float64 {
	return math.Float64frombits(BytesToUint64(bs))
}
