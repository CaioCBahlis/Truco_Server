package main

import (
	"Truco_Server/cardpack"
	"fmt"
	"math/rand"
	"net"
	"time"
)

type Client struct{
	Name string
	IpAddress string
}


func main(){
	Server, err := net.Listen("tcp", ":8080")
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
		go ListenToMe(connection)
		
		fmt.Println(connection)
		message := []byte("Hello, World")
		connection.Write(message)
	}

}


func ShuffleHands(){
	Hands := []cardpack.Card{}
	
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cardpack.Cards), func(i, j int)  {cardpack.Cards[i], cardpack.Cards[j] = cardpack.Cards[j], cardpack.Cards[i]})

	for i := 0; i < 8; i++{
		Hands = append(Hands, cardpack.Card{Name: cardpack.Cards[i], Value: cardpack.Values[cardpack.Cards[i]]})
	}

	fmt.Println(Hands)
}

func ListenToMe(connection net.Conn){
	mybuff := make([]byte, 1024)
	for {
		connection.Read(mybuff)
		fmt.Println(string(mybuff[:]))
	}
}