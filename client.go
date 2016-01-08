/*
Implements the solution to assignment 1 for UBC CS 416 2015 W2.

Usage:
$ go run client.go [local UDP ip:port] [aserver UDP ip:port] [secret]

Example:
$ go run client.go 127.0.0.1:2020 198.162.52.206:1999
*/

package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

/////////// Msgs used by both auth and fortune servers:

// An error message from the server.
type ErrMessage struct {
	Error string
}

/////////// Auth server msgs:

// Message containing a nonce from auth-server.
type NonceMessage struct {
	Nonce int64 `json:"Nonce"`
}

// Message containing an MD5 hash from client to auth-server.
type HashMessage struct {
	Hash string
}

// Message with details for contacting the fortune-server.
type FortuneInfoMessage struct {
	FortuneServer string
	FortuneNonce  int64
}

/////////// Fortune server msgs:

// Message requesting a fortune from the fortune-server.
type FortuneReqMessage struct {
	FortuneNonce int64
}

// Response from the fortune-server containing the fortune.
type FortuneMessage struct {
	Fortune string
}

// Open a UDP connection to a given ip address and port
// Will Exit the program with an error code of 1 if connection fails
func openConnection(localAddr, remoteAddr string) *net.UDPConn {

	laddr, _ := net.ResolveUDPAddr("udp", localAddr)
	raddr, _ := net.ResolveUDPAddr("udp", remoteAddr)

	// Open connection to aServer
	conn, err := net.DialUDP("udp", laddr, raddr)

	// Check if connection was successful
	if err != nil {
		fmt.Println("Something has gone wrong in the initial connection\n")
		os.Exit(1)
	}

	return conn

}

// Send a String using an existing connection to a Server
func sendString(conn *net.UDPConn, message string) {

	// Send a message to a Server

	buffer := []byte(message)
	conn.Write(buffer)
}

// Send bytes using an existing connection to a Server
func sendBytes(conn *net.UDPConn, message []byte) {

	conn.Write(message)

}

// read a String using an existing connection to a Server
// If read fails then Exits wit error code 2
// returns the read message
func readMessage(conn *net.UDPConn) []byte {

	buffer := make([]byte, 1024)

	conn.ReadFromUDP(buffer)

	return buffer
}

func parseNonce(segment []byte) NonceMessage {
	var jsonSection []byte = segment[:29]

	nonce := NonceMessage{}
	err := json.Unmarshal([]byte(jsonSection), &nonce)

	if err != nil {
		fmt.Printf("Error Json: %s\n", err)
	}

	return nonce
}

// Prints the given message on the console
func printConsole(message string) {
	fmt.Println(message)
}

func computeMd5(secret string, nounce int64) string {
	secretNum, _ := strconv.ParseInt(secret, 10, 64)
	var toHash int64 = secretNum + nounce

	buf := make([]byte, 64)

	num := binary.PutVarint(buf, toHash)
	hashByte := buf[:num]

	hasher := md5.New()
	hasher.Write(hashByte)
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}

func createHashMessage(secret string, nounce NonceMessage) []byte {
	var hashString string = computeMd5(secret, nonce.Nonce)

	hash := &HashMessage{
		Hash: hashString,
	}

	packet, _ := json.Marshal(hash)

	return packet
}

// Main workhorse method.
func main() {

	// Arguments
	var localAddr string = os.Args[1]
	var remoteAddr string = os.Args[2]
	var secret string = os.Args[3]

	var conn *net.UDPConn = openConnection(localAddr, remoteAddr)

	sendString(conn, "Arbitrary Payload")

	var message []byte = readMessage(conn)

	var nonce NonceMessage = parseNonce(message)

	var packet []byte = createHashMessage()

	sendBytes(conn, packet)

	response := readMessage(conn)

	fmt.Printf("%s\n", response)
	conn.Close()

}
