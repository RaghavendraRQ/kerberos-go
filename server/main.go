package main

import (
	"fmt"
	"log"
	"net"
	kerberos "renumadam/common"
)

func main() {
	listner, err := net.Listen("tcp", kerberos.SERVER_PORT)

	if err != nil {
		log.Fatalln("Can't able to start: ", err)
	}
	fmt.Println("Service is running on", kerberos.SERVER_PORT)

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

	encServiceTicket := kerberos.ReadData(conn)

	if encServiceTicket == nil {
		log.Println("No data received")
		return
	}

	serviceTicket_ := kerberos.Decrypt(kerberos.SERVICE_KEY[:], encServiceTicket)

	if serviceTicket_ == nil {
		log.Println("Can't decrypt the ticket")
		return
	}

	serviceTicket := kerberos.Decode[kerberos.ServiceTicket](serviceTicket_)

	if serviceTicket.IsExpired() {
		log.Println("Service Ticket is expired")
		kerberos.WriteData(conn, []byte{kerberos.TGS_EXPIRED_ERR})
		return
	}

	kerberos.WriteData(conn, []byte{kerberos.TGS_OK})

	serviceTicket.PrintPretty()
	fmt.Println("Service Ticket is valid")
	fmt.Println("Access Granted to: ", serviceTicket.ClientId)

}
