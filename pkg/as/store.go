package as

import (
	"kerberos/internal"
	"sync"
	"time"
)

var (
	clientIDs int = 1
	clientMut sync.Mutex
	users     = map[string]string{
		"kerberos":  "Password",
		"kerberos1": "Password1",
	}

	sharedKey []byte
)

type Response struct {
	SessionKey []byte
	TGDId      string
	TimeStamp  time.Time
	LifeTime   time.Duration
	TGT        internal.TicketGrantingTicket
}
