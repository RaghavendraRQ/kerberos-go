package main

import (
	"kerberos/pkg/tgs"
	"log"
)

func main() {

	if err := tgs.Run(); err != nil {
		log.Fatalf("Error running TGS: %v", err)
	}
}
