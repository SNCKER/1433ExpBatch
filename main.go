package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
)

var (
	filename     string
	maxWorkerNum int
)

func init() {
	flag.StringVar(&filename, "f", "", "The name of the file in which the target information is stored.")
	flag.IntVar(&maxWorkerNum, "t", 32, "max goroutines(threads).")
}

func main() {
	flag.Parse()
	targets, err := TargetsLoad(filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("targets loaded success.")

	wg := sync.WaitGroup{}
	ch := make(chan bool, maxWorkerNum)
	defer close(ch)
	for {
		fmt.Print("B@tch-Cons0le>")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(os.Stderr, err)
			return
		}
		log.Print("[Exec]", command)

		for _, target := range targets {
			wg.Add(1)
			ch <- true
			go Xp_cmdshell(target, command, &wg, ch)
		}
		wg.Wait()
	}
}
