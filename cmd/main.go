package main

import (
	"kerberos/pkg/as"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		as.Run()

	}()

	wg.Wait()

}
