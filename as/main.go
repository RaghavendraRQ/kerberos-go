package main

import (
	"fmt"
	kerberos "kerberos/common"
	"log"
	"net"
	"strconv"
)

func main() {
	listener, err := net.Listen("tcp", kerberos.AS_PORT)

	if err != nil {
		log.Fatalln("AS can't start: ", err)
	}
	fmt.Println("Authentication Server (AS) is running on port 8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("AS can't connect: ", err)
		}

		fmt.Println("Connection From: ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	username := string(kerberos.ReadData(conn))
	password := string(kerberos.ReadData(conn))

	// TODO: Authenicate the user
	fmt.Println("Username: ", string(username))

	if !authenticateUser(username, password) {
		log.Printf("Authentication failed for user: %s", username)
		kerberos.WriteData(conn, []byte{kerberos.AS_AUTH_ERR})
		return
	}
	kerberos.WriteData(conn, []byte{kerberos.AS_AUTH_OK})

	log.Printf("User %s authenticated", username)
	clientMut.Lock()
	clientId := "client" + strconv.Itoa(clientIDs)
	clientIDs++
	clientMut.Unlock()

	ticketGrantingTicker := kerberos.NewTicketGrantingTicket(clientId)

	encTicketGrantingTicket := kerberos.Encrypt(kerberos.TGS_KEY[:], kerberos.Encode(ticketGrantingTicker))

	kerberos.WriteData(conn, encTicketGrantingTicket)
	fmt.Println("TGT issued to: ", clientId)
}

func authenticateUser(username, password string) bool {

	return users[username] == password

}
