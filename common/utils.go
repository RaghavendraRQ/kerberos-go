package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

func Encrypt(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Error creating AES cipher: %v", err)
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Fatalf("Error generating IV: %v", err)
	}

	paddedData := padData(data, aes.BlockSize)

	ciphertext := make([]byte, len(paddedData))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	return append(iv, ciphertext...)
}

func Decrypt(key, encryptedData []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalf("Error creating AES cipher: %v", err)
	}

	if len(encryptedData) < aes.BlockSize {
		log.Fatalf("Ciphertext too short")
	}
	iv := encryptedData[:aes.BlockSize]
	ciphertext := encryptedData[aes.BlockSize:]

	decrypted := make([]byte, len(ciphertext))

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, ciphertext)

	return unPadData(decrypted)
}

func GenerateSesssionKey() []byte {
	var sessionkey = make([]byte, 16)
	if _, err := rand.Read(sessionkey); err != nil {
		log.Fatalln("Can't generate sessionkey", err)
	}
	return sessionkey
}

func padData(data []byte, blockSize int) []byte {
	padLength := blockSize - (len(data) % blockSize)
	padding := bytes.Repeat([]byte{byte(padLength)}, padLength)
	return append(data, padding...)
}

func unPadData(data []byte) []byte {
	if len(data) == 0 {
		log.Fatal("Decryption error: empty data")
	}
	padding := int(data[len(data)-1])
	if padding <= 0 || padding > len(data) {
		log.Fatal("Decryption error: invalid padding")
	}
	return data[:len(data)-padding]
}

func ReadData(conn net.Conn) []byte {
	var size = make([]byte, 4)

	conn.Read(size)

	var data = make([]byte, int(binary.BigEndian.Uint32(size)))

	conn.Read(data)

	return data

}

func WriteData(conn net.Conn, data []byte) {

	buf := make([]byte, 4)
	size := len(data)
	binary.BigEndian.PutUint32(buf, uint32(size))

	if _, err := conn.Write(buf); err != nil {
		log.Fatalln("Can't write to client: ", err)
	}
	if _, err := conn.Write(data); err != nil {
		log.Fatalln("Can't write to client: ", err)
	}

}

func Encode[Ticket TicketGrantingTicket | ServiceTicket](tgt Ticket) []byte {
	data, err := json.Marshal(tgt)
	if err != nil {
		log.Fatalf("Error marshalling data: %v", err)
	}
	return data
}

func Decode[Ticket TicketGrantingTicket | ServiceTicket](data []byte) Ticket {
	var tgt Ticket
	if err := json.Unmarshal(data, &tgt); err != nil {
		log.Fatalf("Error unmarshalling data: %v", err)
	}
	return tgt
}
