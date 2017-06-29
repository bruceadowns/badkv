package main

import (
	"log"
	"sync"

	"github.com/bruceadowns/badkv/lib"
)

func main() {
	in, err := lib.ParseArgs()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(in)

	// Optionally catch Ctrl-C
	//c := make(chan os.Signal)
	//signal.Notify(c, os.Interrupt)
	//<-c

	log.Println("Start http handlers")
	var wg sync.WaitGroup
	lib.StartRootHandler(&wg, in)

	log.Println("Wait for http handlers to finish")
	wg.Wait()

	log.Println("Done")
}
