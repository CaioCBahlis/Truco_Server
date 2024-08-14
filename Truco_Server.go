package main

import (
	"Truco_Server/cardpack"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

type ServerStruct struct{
	Port string
	Clients []Client
	OnGame bool
}

type Client struct{
	Name string
	IpAddress net.Conn
	CurHand []cardpack.Card
	IsTurn bool
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

		connection.Write([]byte("What would you like to be called"))		
		NameBuff := []byte{}
		n, _ := connection.Read(NameBuff)
		for n == 0{
			n, _ = connection.Read(NameBuff)
		}

		ConnClient :=  Client{Name: string(NameBuff[:]) ,IpAddress:  connection}
		MyServer.Clients = append(MyServer.Clients, ConnClient)
		go MyServer.ListenToMe(ConnClient)
		
		
		if len(MyServer.Clients) != 2{
			Waiting_Message := "Waiting For Players... " + string(len(MyServer.Clients)) + "/2" 
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

func (S *ServerStruct) ListenToMe(MyClient Client){
	if !S.OnGame || MyClient.IsTurn{
		mybuff := make([]byte, 1024)
		for {
			n, _ := MyClient.IpAddress.Read(mybuff)
			if n > 0{fmt.Println(string(mybuff[:]))}
		}
	}
}

func (S *ServerStruct) Start_Game(){
	Card := ShuffleHands()


	FPUi := make([]string, len(cardpack.TripleUI))
    SPUi := make([]string, len(cardpack.TripleUI))
    copy(FPUi, cardpack.TripleUI)
    copy(SPUi, cardpack.TripleUI)

	CardNum := 0
	for _, MyClient := range S.Clients{
		MyClient.CurHand = append(MyClient.CurHand, Card[CardNum], Card[CardNum+1], Card[CardNum+2])
		CardNum += 2
	}

	FPUi[3] = strings.Replace(FPUi[3], "X", string(Card[0].Name[0]), 1)
	FPUi[3] = strings.Replace(FPUi[3], "Y", string(Card[1].Name[0]), 1)
	FPUi[3] = strings.Replace(FPUi[3], "Z", string(Card[2].Name[0]), 1)

	FPUi[5] = strings.Replace(FPUi[5], "X", string([]rune(Card[0].Name)[1]), 1)
	FPUi[5] = strings.Replace(FPUi[5], "Y", string([]rune(Card[1].Name)[1]), 1)
	FPUi[5] = strings.Replace(FPUi[5], "Z", string([]rune(Card[2].Name)[1]), 1)

	FPUi[7] = strings.Replace(FPUi[7], "X", string(Card[0].Name[0]), 1)
	FPUi[7] = strings.Replace(FPUi[7], "Y", string(Card[1].Name[0]), 1)
	FPUi[7] = strings.Replace(FPUi[7], "Z", string(Card[2].Name[0]), 1)

	SPUi[3] = strings.Replace(SPUi[3], "X", string(Card[3].Name[0]), 1)
	SPUi[3] = strings.Replace(SPUi[3], "Y", string(Card[4].Name[0]), 1)
	SPUi[3] = strings.Replace(SPUi[3], "Z", string(Card[5].Name[0]), 1)

	SPUi[5] = strings.Replace(SPUi[5], "X", string([]rune(Card[3].Name)[1]), 1)
	SPUi[5] = strings.Replace(SPUi[5], "Y", string([]rune(Card[4].Name)[1]), 1)
	SPUi[5] = strings.Replace(SPUi[5], "Z", string([]rune(Card[5].Name)[1]), 1)

	SPUi[7] = strings.Replace(SPUi[7], "X", string(Card[3].Name[0]), 1)
	SPUi[7] = strings.Replace(SPUi[7], "Y", string(Card[4].Name[0]), 1)
	SPUi[7] = strings.Replace(SPUi[7], "Z", string(Card[5].Name[0]), 1)


	for i := range(18){
	
		S.Clients[0].IpAddress.Write([]byte(FPUi[i] + "\n"))
		S.Clients[1].IpAddress.Write([]byte(SPUi[i] + "\n"))
	}
	S.Clients[0].IsTurn = true



}