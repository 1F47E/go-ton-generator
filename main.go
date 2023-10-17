package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/xssnick/tonutils-go/ton/wallet"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go-ton-gen <address regexp> like ^EQCKOR to start with or for suffix cat$ to end with")
		os.Exit(1)
	}

	// Format regexp
	addrRegexp := os.Args[1]
	patterns := strings.Split(addrRegexp, ",")
	regexList := make([]*regexp.Regexp, 0)
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			continue
		}
		regexList = append(regexList, re)
	}
	if len(regexList) == 0 {
		log.Fatal("No valid regexp")
	}
	fmt.Printf("Got %d regexp\n", len(regexList))

	// create res dir
	dir := "wallets"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			panic(err)
		}
	}

	// generate N wallets
	// now := time.Now()
	// count := 100
	wg := sync.WaitGroup{}

	// get num of cores
	numCpu := runtime.NumCPU()
	// workPerG := int(count / numCpu)
	wg.Add(numCpu)
	for i := 0; i < numCpu; i++ {
		go func() {
			defer wg.Done()
			// for j := 0; j < workPerG; j++ {
			for {
				words := wallet.NewSeed()

				w, _ := wallet.FromSeed(nil, words, wallet.V4R2)
				addr := w.Address().String()
				// println(addr)
				addrLower := strings.ToLower(addr)
				for _, re := range regexList {
					if re.MatchString(addrLower) {
						// MATCH!
						fmt.Printf("Found address %s\n", addr)

						res := fmt.Sprintf("%s:%s", words, addr)

						keyFile := fmt.Sprintf("%s/%s", dir, addr)
						err := os.WriteFile(keyFile, []byte(res), 0644)
						if err != nil {
							panic(err)
						}

						// found++
					}
				}

				// filter out addresses with - and _
				// if strings.Contains(addr, "-") || strings.Contains(addr, "_") {
				// 	continue
				// }
				// resCh <- addr
			}
		}()
	}

	// read results
	// go func() {
	// 	for addr := range resCh {
	// 		fmt.Println(addr)
	// 	}
	// }()

	wg.Wait()
	// close(resCh)
	// log.Printf("%d/%d good wallets generated in %s\n", goodCnt, count, time.Since(now))
}
