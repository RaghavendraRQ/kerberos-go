package main

import "sync"

var (
	clientIDs int = 1
	clientMut sync.Mutex
	users     = map[string]string{
		"RenuMadam":  "Password",
		"RenuMadam1": "Password1",
	}
)
