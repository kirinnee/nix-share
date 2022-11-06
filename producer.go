package main

import (
	"fmt"
	"net"
)

func send(ip, sendPort, port, content string) {

	pc, err := net.ListenPacket("udp4", ":"+sendPort)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	fmt.Println("Sending: " + ip + ":" + port)

	addr, err := net.ResolveUDPAddr("udp4", ip+":"+port)
	if err != nil {
		panic(err)
	}

	_, err = pc.WriteTo([]byte(content), addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("done!")
}
