package main

import (
	"kerberos/as"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		as.Main()

	}()

	wg.Wait()

}
