package main

import (
	"fmt"
	"net"
	"bufio"
)


var ServerIP  string = "15.228.205.193:8080"
var MyServer Server

type Server struct{
	Server_IP string
	Server_Conn net.Conn
}



func ServerINIT(){
	MyServer = Server{Server_IP: ServerIP}  
	fmt.Println("Connecting to Server")
	MyServer.Connect()
	fmt.Println("You're Now Connected to the Server")
	go MyServer.ListenToServer()
	
}


func (S *Server ) Connect(){
	S.Server_Conn, _ = net.Dial("tcp", MyServer.Server_IP)
	

	fmt.Fprintf(S.Server_Conn,"GET / HTTP/1.0\r\n\r\n")
	status, _ := bufio.NewReader(S.Server_Conn,).ReadString('\n')
	fmt.Println(status)
}

func (S *Server) ListenToServer(){
	var MessageBytes []byte
	for{
		S.Server_Conn.Read(MessageBytes)
		fmt.Println(MessageBytes)
	}
}
