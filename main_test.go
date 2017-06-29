package main

import (
	"log"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	t.Skip()

	t.Log("Start Test Main")
	defer t.Log("End Test Main")
}

func TestTimeConversion(t *testing.T) {
	t.Skip()

	now := time.Now()
	log.Printf("now string: %s", now)
	log.Printf("now nanos: %d", now.UnixNano())
	log.Printf("now millis: %d", now.UnixNano()/1000000)

	d1 := now.UnixNano()
	n1 := time.Unix(0, d1)
	log.Printf("now unix nanos: %d", n1.UnixNano())
}
