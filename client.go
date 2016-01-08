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
	Nonce int64
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

/////////// Helper Functions:

func openConnection(localAddr, remoteAddr string) *net.UDPConn {

	_, port, _ := net.SplitHostPort(localAddr)
	port = ":" + port

	laddr, err := net.ResolveUDPAddr("udp", port)

	if err != nil {
		fmt.Println("Something is Wrong with the given local address\n")
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	raddr, err := net.ResolveUDPAddr("udp", remoteAddr)

	if err != nil {
		fmt.Println("Something is Wrong with the given remote address\n")
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", laddr, raddr)

	if err != nil {
		fmt.Println("Something has gone wrong in the initial connection\n")
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	return conn

}

func sendString(conn *net.UDPConn, message string) {

	buffer := []byte(message)
	conn.Write(buffer)
}

func sendBytes(conn *net.UDPConn, message []byte) {

	conn.Write(message)

}

func readMessage(conn *net.UDPConn) []byte {

	buffer := make([]byte, 1024)

	bytesRead, _, _ := conn.ReadFromUDP(buffer)

	buffer = buffer[:bytesRead]

	return buffer
}

func printConsole(message string) {
	fmt.Printf("%s\n", message)
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

func parseNonceMessage(segment []byte) NonceMessage {

	nonce := NonceMessage{}
	err := json.Unmarshal(segment, &nonce)

	if err != nil {
		fmt.Printf("Error Json Nonce: %s\n", err)
	}

	return nonce
}

func createHashMessage(secret string, nonce NonceMessage) []byte {
	var hashString string = computeMd5(secret, nonce.Nonce)

	hash := &HashMessage{
		Hash: hashString,
	}

	packet, _ := json.Marshal(hash)

	return packet
}

func parseFortuneInfoMessage(message []byte) FortuneInfoMessage {

	fortuneInfo := FortuneInfoMessage{}
	err := json.Unmarshal(message, &fortuneInfo)

	if err != nil {
		fmt.Printf("Error Json Fortune Info: %s\n", err)
	}

	return fortuneInfo

}

func createFortuneReqMessage(nonce int64) []byte {

	fortuneReq := &FortuneReqMessage{
		FortuneNonce: nonce,
	}

	packet, _ := json.Marshal(fortuneReq)

	return packet
}

func parseFortuneMessage(message []byte) FortuneMessage {

	fortune := FortuneMessage{}
	err := json.Unmarshal(message, &fortune)

	if err != nil {
		fmt.Printf("Error Json Fortune: %s\n", err)
	}

	return fortune

}

// Main workhorse method.
func main() {

	// Arguments
	var localAddr string = os.Args[1]
	var remoteAddr string = os.Args[2]
	var secret string = os.Args[3]

	// Hardcoded Arguments for Easier Debugging
	//var localAddr string = "127.0.0.1:2020"
	//var remoteAddr string = "198.162.52.206:1999"
	//var secret string = "2016"

	var conn *net.UDPConn = openConnection(localAddr, remoteAddr)

	sendString(conn, "Arbitrary Payload")

	var message []byte = readMessage(conn)

	var nonce NonceMessage = parseNonceMessage(message)

	var packet []byte = createHashMessage(secret, nonce)

	sendBytes(conn, packet)

	message = readMessage(conn)

	var fortuneinfo FortuneInfoMessage = parseFortuneInfoMessage(message)

	conn.Close()

	conn = openConnection(localAddr, fortuneinfo.FortuneServer)

	packet = createFortuneReqMessage(fortuneinfo.FortuneNonce)

	sendBytes(conn, packet)

	message = readMessage(conn)

	var fortune FortuneMessage = parseFortuneMessage(message)

	fmt.Println(fortune.Fortune)

	conn.Close()

}
