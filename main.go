package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {
	// init ton client
	api := ton.NewAPIClient(nil)

	// generate N wallets
	now := time.Now()
	count := 100
	wg := sync.WaitGroup{}
	goodCnt := 0
	resCh := make(chan string)

	// get num of cores
	numCpu := runtime.NumCPU()
	workPerG := int(count / numCpu)
	wg.Add(numCpu)
	for i := 0; i < numCpu; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < workPerG; j++ {
				words := wallet.NewSeed()

				w, _ := wallet.FromSeed(api, words, wallet.V4R2)
				addr := w.Address().String()

				// filter out addresses with - and _
				if strings.Contains(addr, "-") || strings.Contains(addr, "_") {
					continue
				}
				resCh <- addr
			}
		}()
	}

	// read results
	go func() {
		for addr := range resCh {
			fmt.Println(addr)
		}
	}()

	wg.Wait()
	close(resCh)
	log.Printf("%d/%d good wallets generated in %s\n", goodCnt, count, time.Since(now))
}
