package main

import "sync"

var (
	clientIDs int = 1
	clientMut sync.Mutex
	users     = map[string]string{
		"kerberos":  "Password",
		"kerberos1": "Password1",
	}
)
