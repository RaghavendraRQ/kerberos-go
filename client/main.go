package main

import (
	"fmt"
	kerberos "kerberos/common"
	"log"
	"net"
)

func main() {

	fmt.Println("Client is running.")

	var username, password string

	fmt.Print("Enter username: ")
	fmt.Scanln(&username)
	fmt.Println("Enter password: ")
	fmt.Scanln(&password)

	tgt := requestTGT(username, password)

	if tgt == nil {
		log.Fatalln("Can't get TGT")
	}

	var service int
	fmt.Print("Enter service: ")
	fmt.Scanln(&service)

	st := requestServiceTicket(tgt, service)

	if st == nil {
		log.Fatalln("Can't get Service Ticket")
	}

	checkService(st)

	fmt.Println("You are genius!!!!!")
}

func requestTGT(username, password string) []byte {
	conn, err := net.Dial("tcp", kerberos.AS_PORT)

	if err != nil {
		log.Fatalln("Can't connect to TGT: ", err)
	}
	defer conn.Close()

	kerberos.WriteData(conn, []byte(username))
	kerberos.WriteData(conn, []byte(password))

	code := kerberos.ReadData(conn)

	if code[0] == kerberos.AS_AUTH_ERR {
		log.Println("Authentication failed")
		return nil
	}

	tgt := kerberos.ReadData(conn)
	return tgt

}

func requestServiceTicket(tgt []byte, service int) []byte {
	conn, err := net.Dial("tcp", kerberos.TGS_PORT)
	if err != nil {
		log.Fatalln("Can't connect to TGT: ", err)
	}
	defer conn.Close()

	kerberos.WriteData(conn, tgt)

	code := kerberos.ReadData(conn)

	if code[0] == kerberos.TGT_EXPIRED_ERR {
		log.Println("TGT is expired")
		return nil
	}

	kerberos.WriteData(conn, []byte{byte(service)})

	serviceTicket := kerberos.ReadData(conn)
	return serviceTicket

}

func checkService(st []byte) {
	conn, err := net.Dial("tcp", kerberos.SERVER_PORT)
	if err != nil {
		log.Fatalln("Can't connect to TGT: ", err)
	}
	defer conn.Close()

	kerberos.WriteData(conn, st)

	code := kerberos.ReadData(conn)

	if code[0] == kerberos.TGS_EXPIRED_ERR {
		log.Println("Service Ticket is expired")
		return
	}

}
