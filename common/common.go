package common

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

const (
	AS_AUTH_ERR = 0
	AS_AUTH_OK  = 1

	TGT_EXPIRED_ERR = 2
	TGT_OK          = 3

	TGS_EXPIRED_ERR = 4
	TGS_OK          = 5

	SERVICE_ERR = 6
	SERVICE_OK  = 7
)

var (
	TGS_KEY     [16]byte = [16]byte{0x4a, 0x3f, 0x8c, 0x1d, 0x7e, 0xb2, 0x5a, 0x9f, 0x0d, 0x6e, 0x2c, 0x4f, 0x8a, 0x3b, 0x7d, 0x1e}
	SERVICE_KEY [16]byte = [16]byte{0x4a, 0x3f, 0x8c, 0x1d, 0x7e, 0xb2, 0x5a, 0x9f, 0x0d, 0x6e, 0x2c, 0x4f, 0x8a, 0x3b, 0x7d, 0x1e}
	AS_PORT     string   = ":8000"
	TGS_PORT    string   = ":8001"
	SERVER_PORT string   = ":8002"
)

type TicketGrantingTicket struct {
	ClientId   string
	TimeStamp  time.Time
	Lifetime   time.Time
	SessionKey []byte
}

type ServiceTicket struct {
	ClientId   string
	ServiceId  string
	TimeStamp  time.Time
	Lifetime   time.Time
	SessionKey []byte
}

func (t TicketGrantingTicket) PrintPretty() {
	fmt.Println("\033[1;32mTicket Granting Ticket\033[0m")
	fmt.Println("\033[33mClient ID:\033[0m", t.ClientId)
	fmt.Println("\033[33mTime Stamp:\033[0m", t.TimeStamp)
	fmt.Println("\033[33mLifetime:\033[0m", t.Lifetime)
	fmt.Printf("\033[33mSession Key:\033[0m %x\n", t.SessionKey)
}

func (t TicketGrantingTicket) IsExpired() bool {
	return time.Now().After(t.Lifetime)
}

func (t ServiceTicket) PrintPretty() {
	fmt.Println("\033[1;32mService Ticket\033[0m")
	fmt.Println("\033[33mClient ID:\033[0m", t.ClientId)
	fmt.Println("\033[33mService ID:\033[0m", t.ServiceId)
	fmt.Println("\033[33mTime Stamp:\033[0m", t.TimeStamp)
	fmt.Println("\033[33mLifetime:\033[0m", t.Lifetime)
	fmt.Printf("\033[33mSession Key:\033[0m %x\n", t.SessionKey)
}

func (t ServiceTicket) IsExpired() bool {
	return time.Now().After(t.Lifetime)
}

func NewTicketGrantingTicket(clientId string) TicketGrantingTicket {
	return TicketGrantingTicket{
		ClientId:   clientId,
		TimeStamp:  time.Now(),
		Lifetime:   time.Now().Add(1 * time.Hour),
		SessionKey: generateSessionKey(),
	}
}

func NewServiceTicket(clientId, serviceId string) ServiceTicket {
	return ServiceTicket{
		ClientId:   clientId,
		ServiceId:  serviceId,
		TimeStamp:  time.Now(),
		Lifetime:   time.Now().Add(1 * time.Hour),
		SessionKey: generateSessionKey(),
	}
}

func generateSessionKey() []byte {
	var sessionkey = make([]byte, 16)
	if _, err := rand.Read(sessionkey); err != nil {
		log.Fatalln("Can't generate sessionkey", err)
	}
	return sessionkey

}

// func init() {
// 	if _, err := rand.Read(TGS_KEY[:]); err != nil {
// 		log.Fatalf("Error generating TGS_KEY: %v", err)
// 		fmt.Println("TGS_KEY: ", TGS_KEY)
// 	}

// 	if _, err := rand.Read(SERVICE_KEY[:]); err != nil {
// 		log.Fatalf("Error generating SERVICE KEY: %v", err)
// 	}
// }
