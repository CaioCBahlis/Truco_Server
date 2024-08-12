package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	
)


var ServerIP  string = "52.67.111.171:8080"
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

}


func (S *Server ) Connect(){
	S.Server_Conn, _ = net.Dial("tcp", MyServer.Server_IP)
	go MyServer.ListenToServer()
	go MyServer.WriteToServer()
	select{}

   	
}

func (S *Server) ListenToServer(){
	fmt.Println("Now Listening")
	for{
		tmp := make([]byte, 4096)
		n, _ := S.Server_Conn.Read(tmp)
		if n > 0 {fmt.Println("Server:", string(tmp[:]))}
		
	}
}

func (S *Server) WriteToServer(){
	reader := bufio.NewReader(os.Stdin)
	for{
		fmt.Print("Enter Message: ")
		message, _ := reader.ReadString('\n')
		S.Server_Conn.Write([]byte(message))
	}
}

//func main(){
	//ServerINIT()
//}