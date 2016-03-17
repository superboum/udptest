package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	File    *log.Logger
	Console *log.Logger
)

func udpTest(Conn *net.UDPConn) bool {
	buf := make([]byte, 1024)

	Conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	Conn.Write([]byte("TEST"))
	n, _, _ := Conn.ReadFromUDP(buf)
	Conn.Close()
	return string(buf[:n]) == "OK\n"
}

func initLog() {
	f, err := os.OpenFile("port-checker.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error opening file: %v", err)
	}

	File = log.New(f,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Console = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime)

}

func req(method string, url string, pt PortTest) error {
	ptjson, _ := json.Marshal(pt)
	buf := []byte(ptjson)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(buf))
	req.Header.Set("X-Application", "udptester")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	return err

}

func main() {

	initLog()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("SERVER IP: ")
	url, _ := reader.ReadString('\n')
	url = url[:len(url)-1]

	fmt.Print("SERVER PORT: ")
	port, _ := reader.ReadString('\n')
	port = port[:len(port)-1]

	Console.Println("Starting UDP scan")
	File.Println("Starting UDP scan")

	for i := 1024; i < 65535; i++ {
		pt := PortTest{strconv.Itoa(i)}

		err := req("POST", "http://"+url+":"+port+"/port", pt)
		if err != nil {
			File.Println(err)
			Console.Println(err)
		}

		srvAddr, err := net.ResolveUDPAddr("udp", url+":"+pt.Port)
		Conn, err := net.DialUDP("udp", nil, srvAddr)
		CheckError(err)
		defer Conn.Close()
		if udpTest(Conn) {
			File.Println("UDP Port", pt.Port, "OPEN")
			Console.Println("UDP Port", pt.Port, "OPEN")
		} else {
			File.Println("UDP Port", pt.Port, "CLOSE")
		}

		err = req("DELETE", "http://"+url+":"+port+"/port", pt)
		if err != nil {
			File.Println(err)
			Console.Println(err)
		}

	}
}
