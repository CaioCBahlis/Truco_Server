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
		
		fmt.Println(connection)
		message := []byte("Hello, World")
		connection.Write(message)
		/*
		if len(MyServer.Clients) {
			Waiting_Message, _:= "Waiting For Players... %d/4", len(MyServer.Clients)
			connection.Write([]byte(Waiting_Message))
		}else{
			connection.Write([]byte("Starting Match...."))
			
		}
		*/
		connection.Write([]byte("Starting Match..."))
		MyServer.Start_Game()
	}

}


func ShuffleHands() []cardpack.Card{
	Hands := []cardpack.Card{}
	
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cardpack.Cards), func(i, j int)  {cardpack.Cards[i], cardpack.Cards[j] = cardpack.Cards[j], cardpack.Cards[i]})

	for i := 0; i < 8; i++{
		Hands = append(Hands, cardpack.Card{Name: cardpack.Cards[i], Value: cardpack.Values[cardpack.Cards[i]]})
	}

	return Hands
}

func ListenToMe(connection net.Conn){
	mybuff := make([]byte, 1024)
	for {
		connection.Read(mybuff)
		fmt.Println(string(mybuff[:]))
	}
}

func (S *ServerStruct) Start_Game(){
	Card := ShuffleHands()
	S.Clients[0].IpAddress.Write([]byte(Card[0].Name))
}