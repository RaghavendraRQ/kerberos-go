package tgs

import (
	"fmt"
	kerberos "kerberos/internal"
	"log"
	"net"
)

var (
	services = map[int]string{
		1: "service1",
		2: "service2",
		3: "service3",
	}
)

func Run() error {
	listner, err := net.Listen("tcp", kerberos.TGS_PORT)

	if err != nil {
		return fmt.Errorf("can't start the tgs: %v", err)
	}
	log.Println("Ticket Granting Server (TGS) is running on", kerberos.TGS_PORT)
	defer listner.Close()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		fmt.Println("Connection from: ", conn.RemoteAddr())

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	encTicketGrantingTicket := kerberos.ReadData(conn)
	// TODO: Check for Authenticator as well

	if encTicketGrantingTicket == nil {
		log.Println("No data received")
		return
	}

	ticketGrantingTicket_ := kerberos.Decrypt(kerberos.TGS_KEY[:], encTicketGrantingTicket)

	if ticketGrantingTicket_ == nil {
		log.Println("Can't decrypt the ticket")
		return
	}

	ticketGrantingTicket := kerberos.Decode[kerberos.TicketGrantingTicket](ticketGrantingTicket_)
	ticketGrantingTicket.PrintPretty()

	if ticketGrantingTicket.IsExpired() {
		log.Println("TGT is expired")
		kerberos.WriteData(conn, []byte{kerberos.TGT_EXPIRED_ERR})
		return
	}
	kerberos.WriteData(conn, []byte{kerberos.TGT_OK})

	service := kerberos.ReadData(conn)

	if service == nil {
		log.Println("No service requested")
		return
	}

	reqService, ok := services[int(service[0])]

	if !ok {
		log.Println("Service not found")
		return
	}

	serviceTicket := kerberos.NewServiceTicket(ticketGrantingTicket.ClientId, reqService)

	encServiceTicket := kerberos.Encrypt(kerberos.SERVICE_KEY[:], kerberos.Encode(serviceTicket))

	kerberos.WriteData(conn, encServiceTicket)
	log.Println("SeriveTicket issued to: ", serviceTicket.ClientId)
}
