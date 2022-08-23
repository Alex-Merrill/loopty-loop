package main

import (
	"fmt"
	"os"
	"strconv"

	looper "github.com/Alex-Merrill/loopty-loop/Looper"
)

func main() {
	var path string
	var minDuration, maxDuration int
	var err error

	if len(os.Args) > 4 || len(os.Args) < 2 {
		fmt.Println("ARgs wrong")
		os.Exit(0)
	} else if len(os.Args) == 2 {
		path = os.Args[1]
		minDuration = 1
		maxDuration = 5
	} else if len(os.Args) == 3 {
		path = os.Args[1]
		minDuration, err = strconv.Atoi(os.Args[2])
		checkErr(err)
		maxDuration = 5
	} else {
		path = os.Args[1]
		minDuration, err = strconv.Atoi(os.Args[2])
		checkErr(err)
		maxDuration, err = strconv.Atoi(os.Args[3])
		checkErr(err)
	}

	loop := looper.NewLoop(path, minDuration, maxDuration)

	if res, err := loop.Start(); err != nil {
		panic(err)
	} else {
		fmt.Println(res)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
