package main

import (
	"Truco_Server/cardpack"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type ServerStruct struct{
	Port string
	Clients []Client
}

type Client struct{
	IpAddress net.Conn
}


func main(){
	MyServer := ServerStruct{Port: ":8080", Clients: []Client{}}
	Server, err := net.Listen("tcp", MyServer.Port)
	if err != nil{
		fmt.Println("Error Openning the server")
		return
	}

	for {
		connection, err := Server.Accept()
		if err != nil{
			fmt.Println("Error accepting Client")
			return
		}
		MyServer.Clients = append(MyServer.Clients, Client{connection})
		go ListenToMe(connection)
		
		/*
		fmt.Println(connection)
		message := []byte("Hello, World")
		connection.Write(message)
		
		if len(MyServer.Clients) {
			Waiting_Message, _:= "Waiting For Players... %d/4", len(MyServer.Clients)
			connection.Write([]byte(Waiting_Message))
		}else{
			connection.Write([]byte("Starting Match...."))
		}
		*/

		connection.Write([]byte("Starting Match..."))
		time.Sleep(1 * time.Second)

		connection.Write([]byte("3"))
		time.Sleep(1 * time.Second)

		connection.Write([]byte("2"))
		time.Sleep(1 * time.Second)

		connection.Write([]byte("1"))
		time.Sleep(1 * time.Second)

		MyServer.Start_Game()
	}

}


func ShuffleHands() []cardpack.Card{
	Hands := []cardpack.Card{}
	
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cardpack.Cards), func(i, j int)  {cardpack.Cards[i], cardpack.Cards[j] = cardpack.Cards[j], cardpack.Cards[i]})

	for i := 0; i < 8; i++{
		Hands = append(Hands, cardpack.Card{Name: cardpack.Cards[i], Value: cardpack.Values[cardpack.Cards[i]], Repr: cardpack.CreateTerminalRepr(cardpack.Cards[i])})
	}

	return Hands
}

func ListenToMe(connection net.Conn){
	mybuff := make([]byte, 1024)
	for {
		n, _ := connection.Read(mybuff)
		if n > 0{fmt.Println(string(mybuff[:]))}
	}
}

func (S *ServerStruct) Start_Game(){
	Card := ShuffleHands()
	S.Clients[0].IpAddress.Write([]byte("--------------------------------------------"))
	for i := range(7){
		ImageLine:=  Card[0].Repr[i] + Card[1].Repr[i] + Card[2].Repr[i]
		S.Clients[0].IpAddress.Write([]byte(ImageLine))
	}
	S.Clients[0].IpAddress.Write([]byte("--------------------------------------------"))
}