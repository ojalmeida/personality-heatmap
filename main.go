package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"personality-heatmap/requesting"
	"time"
)

func main() {

	requesting.Start()

}

func test() {

	writer := uilive.New()
	// start listening for updates and render
	writer.Start()

	for i := 0; i <= 100; i++ {
		fmt.Fprintf(writer, "Downloading.. (%d/%d) GB\n", i, 100)
		time.Sleep(time.Millisecond * 50)
	}

	fmt.Fprintln(writer, "Finished: Downloaded 100GB")
	writer.Stop() // flush and stop rendering

}
