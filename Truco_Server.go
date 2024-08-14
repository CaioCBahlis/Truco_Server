package main

import (
	"Truco_Server/cardpack"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
	"strconv"
)

type ServerStruct struct{
	Port string
	Clients []Client
	OnGame bool
	Round int
	CardsOnTable []cardpack.Card
}

type Client struct{
	Name string
	IpAddress net.Conn
	CurHand []cardpack.Card
	PlayerIndex int
	IsTurn bool
	Played bool
}

type Game struct{
	CardsOnTable []cardpack.Card
	Round int
}


func main(){
	MyServer := ServerStruct{Port: ":8080", Clients: []Client{}}
	fmt.Println("Server is Running")
	Server, err := net.Listen("tcp", MyServer.Port)
	if err != nil{
		fmt.Println("Error Openning the server")
		return
	}

	PlayerIndex := 0
	for {
		connection, err := Server.Accept()
		if err != nil{
			fmt.Println("Error accepting Client")
			return
		}

		connection.Write([]byte("What would you like to be called"))		
		NameBuff := make([]byte, 1024)
		connection.Read(NameBuff)
		

		ConnClient :=  Client{Name: string(NameBuff[:]) ,IpAddress:  connection, PlayerIndex: PlayerIndex}
		MyServer.Clients = append(MyServer.Clients, ConnClient)
		go MyServer.ListenToMe(PlayerIndex)
		PlayerIndex += 1		
		
		if len(MyServer.Clients) != 2{
			Waiting_Message := "Waiting For Players... 1/2 Players Ready" 
			connection.Write([]byte(Waiting_Message))
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

		MyServer.OnGame = true
		MyServer.Round = 1
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

func (S *ServerStruct) ListenToMe(PlayerIndex int){
	
	mybuff := make([]byte, 1024)
	for {
		n, _ := S.Clients[PlayerIndex].IpAddress.Read(mybuff)
		Message := strings.TrimSpace(string(mybuff[:n]))

		if !S.OnGame{
			if n > 0{fmt.Println(Message)}
		}else if S.Clients[PlayerIndex].IsTurn{
			fmt.Println(Message)
			if n > 0{
				switch Message{
					case "Jogar":
						Index := make([]byte, 1024)
						S.Clients[PlayerIndex].IpAddress.Write([]byte("Enter the index of your card (1-3)"))
						sz, _ := S.Clients[PlayerIndex].IpAddress.Read(Index)
						Num, _ := strconv.Atoi(string(Index[:sz]))
						CardIndex := Num -1

						for CardIndex > len(S.Clients[PlayerIndex].CurHand)-1 || CardIndex < 0{
							S.Clients[PlayerIndex].IpAddress.Write([]byte("Invalid Index"))
							S.Clients[PlayerIndex].IpAddress.Write([]byte("Enter the index of your card (1-3)"))
							S.Clients[PlayerIndex].IpAddress.Read(Index)
							Num, _ = strconv.Atoi(string(Index[:]))
						}
						
						PlayedCard := S.Clients[PlayerIndex].CurHand[CardIndex]
						S.Clients[PlayerIndex].CurHand = append(S.Clients[PlayerIndex].CurHand[:CardIndex], S.Clients[PlayerIndex].CurHand[CardIndex+1:]...)
						S.CardsOnTable = append(S.CardsOnTable, PlayedCard)

						fmt.Println(S.Clients[PlayerIndex].CurHand[Num-1].Name)
					case "Truco":
						fmt.Println("Received")
					case "Envido":
						fmt.Println("Received")
					case "Queimar":
						fmt.Println("Received")
					case "Correr":
						fmt.Println("Received")
					case "Flor":
						fmt.Println("Received")
				}
			}
		}
	}
}

func (S *ServerStruct) Start_Game(){
	Card := ShuffleHands()

	CardNum := 0
	for _, MyClient := range S.Clients{
		MyClient.CurHand = append(MyClient.CurHand, Card[CardNum], Card[CardNum+1], Card[CardNum+2])
		CardNum += 3
	}

	var Gui []string
	for _, Client := range S.Clients{
		Gui = cardpack.UpdateGui(S.Round, Client.CurHand)
		for i := range(18){
			Client.IpAddress.Write([]byte(Gui[i] + "\n"))
		}
	}

	for S.Round < 3 || S.OnGame{
		for idx, _ := range(S.Clients){
			Client := S.Clients[idx]
			Client.IsTurn = true
			for !Client.Played{
				fmt.Println("Waiting for Player...")
				time.Sleep(1 * time.Second)
			}
		}
		fmt.Println(S.CardsOnTable)
		S.Round += 1
		S.CardsOnTable = make([]cardpack.Card, 4)

	}


	select{}
}