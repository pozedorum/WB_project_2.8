package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	PORT := ":9090"
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Error listening:", err.Error())

	}

	defer listener.Close()
	fmt.Println("Server is listening on " + PORT)

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error accepting:", err.Error())

		}
		fmt.Println("Connected with", conn.RemoteAddr().String())

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		clientMessage := scanner.Text()
		fmt.Printf("Received from client: %s\n", clientMessage)
		if _, err := conn.Write([]byte(clientMessage + "\n")); err != nil {
			log.Fatal("request error: ", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading:", err.Error())
	}
}
