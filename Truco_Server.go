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
	PointsOnWin int
}

type Client struct{
	Name string
	IpAddress net.Conn
	CurHand []cardpack.Card
	PlayerIndex int
	IsTurn bool
	Played bool
	RoundsWon int
	Points int
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

		connection.Write([]byte("\n" +"What would you like to be called"))		
		NameBuff := make([]byte, 1024)
		NmSize, _ := connection.Read(NameBuff)
		

		ConnClient :=  Client{Name: string(NameBuff[:NmSize]) ,IpAddress:  connection, PlayerIndex: PlayerIndex}
		MyServer.Clients = append(MyServer.Clients, ConnClient)
		go MyServer.ListenToMe(PlayerIndex)
		PlayerIndex += 1		
		
		if len(MyServer.Clients) != 2{
			Waiting_Message := "Waiting For Players... 1/2 Players Ready" 
			connection.Write([]byte("\n" + Waiting_Message))
		}else{
			break
		}
	}

		
		MyServer.Clients[0].IpAddress.Write([]byte("\n" +"Starting Match...."))
		MyServer.Clients[1].IpAddress.Write([]byte("\n" +"Starting Match...."))
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

func (S *ServerStruct) BroadCast(message string){
	for _, Client := range(S.Clients){
		Client.IpAddress.Write([]byte("\n" + message))
	}
}


func (S *ServerStruct) ListenToMe(PlayerIndex int){
	mybuff := make([]byte, 1024)
	for {
		n, _ := S.Clients[PlayerIndex].IpAddress.Read(mybuff)
		Message := strings.TrimSpace(string(mybuff[:n]))

		if !S.OnGame{
			if n > 0{fmt.Println(Message)}
		}else if S.Clients[PlayerIndex].IsTurn{
			fmt.Println(S.Clients[PlayerIndex].Name + ": "+ Message)
			if n > 0{
				switch Message{
					case "Jogar":
						S.Jogar(PlayerIndex)
						
					case "Truco":
						Response := make([]byte, 16)
						S.BroadCast(fmt.Sprintf("%s PEDIU TRUCO NEWBA", S.Clients[PlayerIndex].Name))
						var MyClient Client
						if PlayerIndex == 1{
							S.Clients[PlayerIndex-1].IpAddress.Write([]byte("VAI ACEITAR (y/n)"))
							S.Clients[PlayerIndex-1].IsTurn = true
							MyClient = S.Clients[PlayerIndex-1]
						}else{
							S.Clients[PlayerIndex+1].IpAddress.Write([]byte("VAI ACEITAR (y/n)"))
							S.Clients[PlayerIndex+1].IsTurn = true
							MyClient = S.Clients[PlayerIndex-1]
						}
							sz, _ :=  MyClient.IpAddress.Read(Response)
							ClientResponse := string(Response[:sz])
							if ClientResponse == "y"{
								S.PointsOnWin = 3
								S.BroadCast("Truco Aceito")
								S.Jogar(PlayerIndex)
								
							}else if ClientResponse == "n"{

								if S.Clients[PlayerIndex].PlayerIndex == 0{
									S.Clients[PlayerIndex + 1].RoundsWon = 2
								}else{
									S.Clients[PlayerIndex -1].RoundsWon = 2
								}
		
								PlayedCard := cardpack.Card{Name: "Queimar", Value: 0, Repr: cardpack.QueimadoCard}
								S.CardsOnTable = append(S.CardsOnTable, PlayedCard)
								S.Clients[PlayerIndex].Played = true
							}
						
						
					case "Envido":
						fmt.Println("Received")

					case "Queimar":
						S.Jogar(PlayerIndex)
						S.CardsOnTable[len(S.CardsOnTable)-1] = cardpack.Card{Name: "Resign", Value: 0, Repr: cardpack.ResignationCard}
						
					case "Correr":
						if S.Clients[PlayerIndex].PlayerIndex == 0{
							S.Clients[PlayerIndex + 1].RoundsWon = 2
						}else{
							S.Clients[PlayerIndex -1].RoundsWon = 2
						}

						PlayedCard := cardpack.Card{Name: "Resign", Value: 0, Repr: cardpack.ResignationCard}
						S.CardsOnTable = append(S.CardsOnTable, PlayedCard)
						S.Clients[PlayerIndex].Played = true

					case "Flor":
						fmt.Println("Received")
				}
			}
		}else{
			S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" + "Not your turn"))
		}
	}
}

func (S *ServerStruct) Jogar(PlayerIndex int){
	Index := make([]byte, 1024)
	S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" +"Enter the index of your card (1-3)"))

	sz, _ := S.Clients[PlayerIndex].IpAddress.Read(Index)
	Num, _ := strconv.Atoi(strings.TrimSpace(string(Index[:sz])))
	CardIndex := Num -1

	for CardIndex > len(S.Clients[PlayerIndex].CurHand)-1 || CardIndex < 0{
		S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" +"Invalid Index"))
		S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" +"Enter the index of your card (1-3)"))
		sz, _ = S.Clients[PlayerIndex].IpAddress.Read(Index)
		Num, _ = strconv.Atoi(strings.TrimSpace(string(Index[:sz])))
		CardIndex = Num -1
	}
						
	PlayedCard := S.Clients[PlayerIndex].CurHand[CardIndex]
	S.Clients[PlayerIndex].CurHand = append(S.Clients[PlayerIndex].CurHand[:CardIndex], S.Clients[PlayerIndex].CurHand[CardIndex+1:]...)
	S.CardsOnTable = append(S.CardsOnTable, PlayedCard)
	S.Clients[PlayerIndex].Played = true
}

func (S *ServerStruct) Start_Game(){
	PlayingOrder := S.Clients


	var Gui []string
	for S.Clients[0].Points < 12 && S.Clients[1].Points < 12{
		S.Round = 1
		S.Clients[0].RoundsWon = 0
		S.Clients[1].RoundsWon = 0
		S.PointsOnWin = 1
		Card := ShuffleHands()
		CardNum := 0

		for idx := range(len(S.Clients)){
			S.Clients[PlayingOrder[idx].PlayerIndex].CurHand = append(S.Clients[PlayingOrder[idx].PlayerIndex].CurHand, Card[CardNum], Card[CardNum+1], Card[CardNum+2])
	
			CardNum += 3
		}

		for S.Round <= 3 && S.Clients[0].RoundsWon < 2 && S.Clients[1].RoundsWon < 2 {
			S.Clients[0].Played = false
			S.Clients[1].Played = false
			S.CardsOnTable  = []cardpack.Card{}

			for idx := range(len(S.Clients)){
				Gui = cardpack.UpdateGui(S.Round,  S.Clients[PlayingOrder[idx].PlayerIndex].CurHand)
				S.Clients[PlayingOrder[idx].PlayerIndex].IpAddress.Write([]byte("\n"))

				for i := range(18){
					S.Clients[PlayingOrder[idx].PlayerIndex].IpAddress.Write([]byte(Gui[i] + "\n"))
				}
			}
			
			for idx := range(len(S.Clients)){
				S.Clients[PlayingOrder[idx].PlayerIndex].IsTurn = true
				S.Clients[PlayingOrder[idx].PlayerIndex].IpAddress.Write([]byte("\n" + "It's Your Turn!"))

				for !S.Clients[PlayingOrder[idx].PlayerIndex].Played {
					fmt.Println("\n" + "Waiting for" +  S.Clients[PlayingOrder[idx].PlayerIndex].Name + "...")
					time.Sleep(5 * time.Second)
				}

				
				S.BroadCast(S.Clients[PlayingOrder[idx].PlayerIndex].Name + "Played: ")
				Card := cardpack.CreateTerminalRepr(S.CardsOnTable[len(S.CardsOnTable)-1].Name)
				for i := range(7){
					S.BroadCast(Card[i])
				}

				S.Clients[PlayingOrder[idx].PlayerIndex].Played = false
				S.Clients[PlayingOrder[idx].PlayerIndex].IsTurn = false
			}
			

			if cardpack.Values[S.CardsOnTable[0].Name] > cardpack.Values[S.CardsOnTable[1].Name] {
				S.Clients[PlayingOrder[0].PlayerIndex].RoundsWon += 1
				S.BroadCast(S.Clients[PlayingOrder[0].PlayerIndex].Name  + "Won the Round")

			}else if cardpack.Values[S.CardsOnTable[0].Name] < cardpack.Values[S.CardsOnTable[1].Name]{
				S.Clients[PlayingOrder[1].PlayerIndex].RoundsWon += 1
				S.BroadCast(S.Clients[PlayingOrder[1].PlayerIndex].Name + "Won the Round")
				PlayingOrder = []Client{PlayingOrder[1], PlayingOrder[0]}
			
			}else{
				S.Clients[PlayingOrder[0].PlayerIndex].RoundsWon += 1
				S.Clients[PlayingOrder[1].PlayerIndex].RoundsWon += 1
				S.BroadCast("Draw")
			}
			S.Round += 1
			S.CardsOnTable = make([]cardpack.Card, 4)
		}
	
	
		if S.Clients[PlayingOrder[0].PlayerIndex].RoundsWon > S.Clients[PlayingOrder[1].PlayerIndex].RoundsWon{
			fmt.Println("Player 1 Won")
			S.Clients[PlayingOrder[0].PlayerIndex].Points += S.PointsOnWin
		}else if S.Clients[PlayingOrder[0].PlayerIndex].RoundsWon < S.Clients[PlayingOrder[1].PlayerIndex].RoundsWon{
			fmt.Println("Player 2 Won")
			S.Clients[PlayingOrder[1].PlayerIndex].Points += S.PointsOnWin
		}else{
			fmt.Println("Draw")
		}
		S.BroadCast(fmt.Sprintf("%s: %d",S.Clients[PlayingOrder[0].PlayerIndex].Name, S.Clients[PlayingOrder[0].PlayerIndex].Points))
		S.BroadCast(fmt.Sprintf("%s: %d", S.Clients[PlayingOrder[1].PlayerIndex].Name, S.Clients[PlayingOrder[1].PlayerIndex].Points))
	}
}