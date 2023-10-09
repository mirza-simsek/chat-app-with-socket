package main

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn, clients map[net.Conn]bool) {
	defer conn.Close()

	for client := range clients {
		if client != conn {
			conn.Write([]byte("Sohbete yeni bir kullanıcı eklendi"))
		}
	}
	clients[conn] = true

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			conn.Close()
			break
		}
		message := string(buffer[:n])
		fmt.Println("Message is :", message)
		for client := range clients {
			if client != conn {
				_, err := client.Write([]byte(message))
				if err != nil {
					fmt.Println(err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP is listening on 12345")

	clients := make(map[net.Conn]bool)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go HandleConnection(conn, clients)
	}
}
