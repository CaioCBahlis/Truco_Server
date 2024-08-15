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
	PlayingOrder []Client
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
						Index := make([]byte, 1024)
						S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" +"Enter the index of your card (1-3)"))

						sz, _ := S.Clients[PlayerIndex].IpAddress.Read(Index)
						Num, _ := strconv.Atoi(strings.TrimSpace(string(Index[:sz])))
						CardIndex := Num -1

						for CardIndex > len(S.Clients[PlayerIndex].CurHand)-1 || CardIndex < 0{
							S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" +"Invalid Index"))
							S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" +"Enter the index of your card (1-3)"))
							S.Clients[PlayerIndex].IpAddress.Read(Index)
							Num, _ = strconv.Atoi(string(Index[:]))
							CardIndex = Num -1
						}
						
						PlayedCard := S.Clients[PlayerIndex].CurHand[CardIndex]
						S.Clients[PlayerIndex].CurHand = append(S.Clients[PlayerIndex].CurHand[:CardIndex], S.Clients[PlayerIndex].CurHand[CardIndex+1:]...)
						S.CardsOnTable = append(S.CardsOnTable, PlayedCard)
						S.Clients[PlayerIndex].Played = true
						
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
		}else{
			S.Clients[PlayerIndex].IpAddress.Write([]byte("\n" + "Not your turn"))
		}
	}
}

func (S *ServerStruct) Start_Game(){
	S.PlayingOrder = S.Clients


	var Gui []string
	for S.Clients[0].Points < 12 && S.Clients[1].Points < 12{

		Card := ShuffleHands()
		CardNum := 0
		for idx := range(len(S.Clients)){
			S.PlayingOrder[idx].CurHand = append(S.PlayingOrder[idx].CurHand, Card[CardNum], Card[CardNum+1], Card[CardNum+2])
			CardNum += 3
		}

		for S.Round <= 3 && S.Clients[0].RoundsWon != 2 && S.Clients[1].RoundsWon != 2 {

			S.CardsOnTable  = []cardpack.Card{}
			for idx := range(len(S.PlayingOrder)){
				Gui = cardpack.UpdateGui(S.Round,  S.PlayingOrder[idx].CurHand)
				S.PlayingOrder[idx].IpAddress.Write([]byte("\n"))

				for i := range(18){
					S.PlayingOrder[idx].IpAddress.Write([]byte(Gui[i] + "\n"))
				}
			}
			
			for idx := range(len(S.PlayingOrder)){
				S.PlayingOrder[idx].IsTurn = true
				S.PlayingOrder[idx].IpAddress.Write([]byte("\n" + "It's Your Turn!"))
				for !S.PlayingOrder[idx].Played{
					fmt.Println("\n" + "Waiting for" +  S.PlayingOrder[idx].Name + "...")
					time.Sleep(5 * time.Second)
				}

				S.BroadCast(S.PlayingOrder[idx].Name + "Played: ")
				Card := cardpack.CreateTerminalRepr(S.CardsOnTable[len(S.CardsOnTable)-1].Name)
				for i := range(7){
					S.BroadCast(Card[i])
				}

				S.PlayingOrder[idx].Played = false
				S.PlayingOrder[idx].IsTurn = false
			}
			

			if cardpack.Values[S.CardsOnTable[0].Name] > cardpack.Values[S.CardsOnTable[1].Name] {
				S.PlayingOrder[0].RoundsWon += 1
				S.BroadCast(S.PlayingOrder[0].Name  + "Won the Round")

			}else if cardpack.Values[S.CardsOnTable[0].Name] < cardpack.Values[S.CardsOnTable[1].Name]{
				S.PlayingOrder[1].RoundsWon += 1
				S.BroadCast(S.PlayingOrder[1].Name + "Won the Round")
				S.PlayingOrder = []Client{S.Clients[1], S.Clients[[0]]}
			
				
			}else{
				S.BroadCast("Draw")
			}
			S.Round += 1
			S.CardsOnTable = make([]cardpack.Card, 4)
		}

	if S.PlayingOrder[0].RoundsWon > S.PlayingOrder[1].RoundsWon{
		fmt.Println("Player 1 Won")
		S.PlayingOrder[0].Points += 1
	}else if S.PlayingOrder[0].RoundsWon < S.PlayingOrder[1].RoundsWon{
		fmt.Println("Player 2 Won")
		S.PlayingOrder[1].Points += 1
	}else{
		fmt.Println("Draw")
	}
	S.BroadCast(fmt.Sprintf(S.PlayingOrder[0].Name, S.PlayingOrder[0].Points))
	S.BroadCast(fmt.Sprintf(S.PlayingOrder[1].Name, S.PlayingOrder[1].Points))
	}
}