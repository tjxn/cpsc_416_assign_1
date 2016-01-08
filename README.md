# CPSC 416 - Distributed Systems: Assignment 1
## University of British Columbia

High-level protocol description

There are three kinds of nodes in the system: a client (that you will implement), an authentication server (aserver) 
to verify a client's credentials, and a fortune server (fserver) that returns a fortune string to the client. 
You will test and debug your solution against running aserver and fserver instances. However, you will not have access 
to their code. You are given initial starter code (below) which contains the expected message type declarations. Both 
aserver and fserver serve multiple clients simultaneously. Each client implements a sequential control flow, interacting 
with aserver first, and later with the fserver. The client communicates with both servers over UDP, using 
binary-encoded JSON messages.

The client is run knowing the UDP IP:port of the aserver and an int64 secret value. It follows the following steps:

1. the client sends a UDP message with arbitrary payload to the aserver
2. the client receives a NonceMessage reply containing an int64 nonce from the aserver
3. the client computes an MD5 hash of the (nonce + secret) value and sends this value as a hex 
		string to the aserver as part of a HashMessage
4. the aserver verifies the received hash and replies with a FortuneInfoMessage that contains information 
		for contacting fserver (its UDP IP:port and an int64 fortune nonce to use when connecting to it)
5. the client sends a FortuneReqMessage to fserver
6. the client receives a FortuneMessage from the fserver
7. the client prints out the received fortune string as the last thing before exiting on a new newline-terminated line
8. The communication steps in this protocol are illustrated in the following space-time diagram:

![Space-Time Diagram](/tjxn/cpsc_416_assign_1/assign1-proto.jpg)
 
Protocol corner-cases

- The aserver and fserver expect the client to use the same UDP IP:port throughout the protocol. 
- An ErrMessage will be generated if the client uses different addresses.
- The aserver will reply with an ErrMessage in case the MD5 hex string value does not match the expected value.
- The fserver will reply to all messages that it perceives as not being HashMessage with a new NonceMessage.
- The fserver will reply with an ErrMessage in case it cannot unmarshal the message from the client.
- The fserver will reply with an ErrMessage in case an incorrect fortune nonce is supplied in the FortuneReqMessage.
- All messages fit into 1024 bytes.


Implementation requirements

- The client code must be runnable on CS ugrad machines and be compatible with Go version 1.4.3.
- Your code does not need to check or handle ErrMessage replies from the aserver or fserver. However, 
- you may find it useful to check for these replies during debugging.
- Your code may assume that UDP is reliable and not implement any retransmission.
- You must use UDP and the message types given out in the initial code.
- Your solution can only use standard library Go packages.
- Your solution code must be Gofmt'd using gofmt.


Solution spec

Write a single go program called client.go that acts as a client in the protocol described above. 
Your program must implement the following command line usage:
go run client.go [local UDP ip:port] [aserver ip:port] [secret]

- [local UDP ip:port] : local address that the client uses to connect to both the aserver and the fserver 
		(i.e., the external IP of the machine the client is running on)
- [aserver UDP ip:port] : the UDP address on which the aserver receives new client connections
- [secret] : an int64 secret

