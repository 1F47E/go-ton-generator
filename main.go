package main

import (
	"log"
	"math/rand"
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
	count := 10
	wg := sync.WaitGroup{}
	wg.Add(count)
	goodCnt := 0
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			// random sleep from 0.1 to 1 sec
			time.Sleep(time.Duration(rand.Intn(1000)+100) * time.Millisecond)
			words := wallet.NewSeed()

			w, _ := wallet.FromSeed(api, words, wallet.V4R2)
			addr := w.Address().String()

			// filter out addresses with - and _
			if !strings.Contains(addr, "-") && !strings.Contains(addr, "_") {
				log.Println("good wallet: ", addr)
				goodCnt++
			}
		}()
	}
	wg.Wait()
	log.Printf("%d/%d good wallets generated in %s\n", goodCnt, count, time.Since(now))
}
