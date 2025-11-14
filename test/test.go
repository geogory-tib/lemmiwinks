package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"lemiwinks/server"
	"lemiwinks/util"
)

func main() {
	ipandport := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])
	conn, err := net.Dial("tcp", ipandport)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	reader := bufio.NewScanner(os.Stdin)
	login(conn)
	msg := server.Message{}
	response_buffer := make([]byte, 4096)
	for {
		fmt.Println("who is the message for")
		reader.Scan()
		userin := reader.Text()
		msg.To = userin
		fmt.Println("type the contents")
		reader.Scan()
		userin = reader.Text()
		msg.Contents = userin
		msg.Time = time.Now().Format(time.ANSIC)
		json_data, err := json.Marshal(msg)
		if err != nil {
			log.Panic(err)
		}
		conn.Write(json_data)
		util.ReadFullJson(conn, response_buffer)
		fmt.Print(string(response_buffer))

	}
}

func login(svr net.Conn) {
	response_buffer := make([]byte, 1096)
	login_struct := server.User{}
	reader := bufio.NewScanner(os.Stdin)
	fmt.Print("Type user name and password")
	reader.Scan()
	userin := reader.Text()
	login_struct.Username = userin
	reader.Scan()
	userin = reader.Text()
	login_struct.Password = userin
	login_json, err := json.Marshal(login_struct)
	if err != nil {
		log.Panic(err)
	}
	svr.Write(login_json)
	n, err := util.ReadFullJson(svr, response_buffer)
	if err != nil {
		log.Panic(err)
	}
	response_json := response_buffer[:n-1]
	fmt.Println(string(response_json))
}
