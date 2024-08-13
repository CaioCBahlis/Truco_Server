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
	fmt.Println("Server is Running")
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
		
		
		if len(MyServer.Clients) != 2{
			Waiting_Message, _:= "Waiting For Players... %d/2", len(MyServer.Clients)
			connection.Write([]byte(Waiting_Message))
			time.Sleep(1 * time.Second)
		}else{
			break
		}
	}

		
		MyServer.Clients[0].IpAddress.Write([]byte("Starting Match...."))
		MyServer.Clients[1].IpAddress.Write([]byte("Starting Match...."))
		time.Sleep(1 * time.Second)

		
	
		MyServer.Clients[0].IpAddress.Write([]byte("3"))
		MyServer.Clients[1].IpAddress.Write([]byte("3"))
		time.Sleep(1 * time.Second)

		MyServer.Clients[0].IpAddress.Write([]byte("2"))
		MyServer.Clients[1].IpAddress.Write([]byte("2"))
		time.Sleep(1 * time.Second)

		MyServer.Clients[0].IpAddress.Write([]byte("1"))
		MyServer.Clients[1].IpAddress.Write([]byte("1"))
		time.Sleep(1 * time.Second)

		MyServer.Start_Game()
	
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
	S.Clients[1].IpAddress.Write([]byte("--------------------------------------------"))

	for i := range(7){
		ImageLine:=  Card[0].Repr[i] + Card[1].Repr[i] + Card[2].Repr[i]
		S.Clients[0].IpAddress.Write([]byte(ImageLine + "\n"))
		Image2Line := Card[3].Repr[i] + Card[4].Repr[i] + Card[5].Repr[i]
		S.Clients[0].IpAddress.Write([]byte(Image2Line + "\n"))
	}
	S.Clients[0].IpAddress.Write([]byte("--------------------------------------------"))
	S.Clients[1].IpAddress.Write([]byte("--------------------------------------------"))

}