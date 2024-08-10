package main

import (
	"fmt"
	"net"
	"bufio"
)


var ServerIP  string = "172.31.12.212:8080"
var MyServer Server

type Server struct{
	Server_IP string
}



func ServerINIT(){
	MyServer = Server{Server_IP: ServerIP}  
	fmt.Println("Connecting to Server")
	MyServer.Connect()
	fmt.Println("You're Now Connected to the Server")
	
}


func (S *Server ) Connect(){
	conn, err := net.Dial("tcp", MyServer.Server_IP)
	if err != nil{
		fmt.Println("Error connecting to the server")
	}

	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(status)


}