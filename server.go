package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
)

var connections map[string]*net.UDPConn

func openPort(pt PortTest) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Port", pt.Port, "already in use")
		}
	}()

	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+pt.Port)
	CheckError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	connections[pt.Port] = ServerConn
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		log.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			return
		}
		_, err = ServerConn.WriteToUDP([]byte("OK\n"), addr)
		CheckError(err)
	}
}

func closePort(pt PortTest) {
	if connections[pt.Port] != nil {
		connections[pt.Port].Close()
	}
}

func portHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var pt PortTest
	err := decoder.Decode(&pt)

	if !CheckError(err) {
		log.Println("Invalid request")
		return
	} else if r.Method == "POST" {
		log.Println("Open Port", pt.Port)
		go openPort(pt)
	} else if r.Method == "DELETE" {
		log.Println("Close Port", pt.Port)
		closePort(pt)
	}

}

func main() {
	connections = make(map[string]*net.UDPConn, 1)
	http.HandleFunc("/port", portHandler)
	http.ListenAndServe(":8080", nil)
}
