package as

import (
	"fmt"
	kerberos "kerberos/common"
	"log"
	"net"
	"strconv"
)

func Main() {
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

	exchageKey(conn) // Stroring the key directly on the sharedKey

	fmt.Printf("Shared Key: %x\n", sharedKey)
	username := string(kerberos.Decrypt(sharedKey, kerberos.ReadData(conn)))
	password := string(kerberos.Decrypt(sharedKey, kerberos.ReadData(conn)))

	// TODO: Authenicate the user
	fmt.Println("Username: ", string(username))

	if !AuthenticateUser(username, password) {
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

func AuthenticateUser(username, password string) bool {

	return users[username] == password

}

func exchageKey(conn net.Conn) error {
	publicKey_s, privateKey_s := kerberos.GenerateKeyPair()

	kerberos.WriteData(conn, publicKey_s.Bytes())

	publicKey_c_bytes := kerberos.ReadData(conn)

	publickey_c := kerberos.GetKeyFromBytes(publicKey_c_bytes)

	sharedKey = kerberos.GenerateSharedKey(publickey_c, privateKey_s)
	return nil

}
