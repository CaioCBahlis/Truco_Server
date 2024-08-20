package main

import (
	"Truco_Server/cardpack"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

const(
	NROUNDS = 3
	MAXWONROUND = 2
	MAXPOINTS = 12
	NCARDS = 3
)

type ServerStruct struct{
	Port string
	Clients []Client
	OnGame bool
	PlayerSlots int
}

type Client struct{
	Name string
	IpAddress net.Conn
}

type Team struct{
	TeamName string
	TeamPlayers []*Player
	TeamPoints int
	RoundsWon int
	Resigned bool
	Challenged bool
	Accepted string
}

type Player struct{
	Client
	MyTeam *Team
	CurHand []cardpack.Card
	IsTurn bool
	Played bool
}

type Game struct{
	Teams []*Team
	Players []*Player
	CardsOnTable []cardpack.Card
	Round int
	PointsOnWin int
}


func main(){
	MyServer := ServerStruct{Port: ":8080", Clients: []Client{}}
	fmt.Println("Server is Listening on port 8080")
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

		connection.Write([]byte("\n" +"What would you like to be called"))		
		NameBuff := make([]byte, 1024)
		NmSize, _ := connection.Read(NameBuff)
		

		ConnectedClient :=  Client{Name: string(NameBuff[:NmSize]) , IpAddress:  connection}
		MyServer.Clients = append(MyServer.Clients, ConnectedClient)

		HostResponse := make([]byte, 1024)
		var Response string

		if len(MyServer.Clients) == 1{
			MyServer.Clients[0].IpAddress.Write([]byte("How many players? (2 or 4)"))
			sz, _ := MyServer.Clients[0].IpAddress.Read(HostResponse)
			Response = strings.TrimSpace(string(HostResponse[:sz]))

			for Response != "2" && Response != "4"{
				MyServer.Clients[0].IpAddress.Write([]byte("Invalid Number of Players"))
				MyServer.Clients[0].IpAddress.Write([]byte("How many players? (2 or 4)"))
				sz, _ := MyServer.Clients[0].IpAddress.Read(HostResponse)
				Response = string(HostResponse[:sz])
			}
	
			MyServer.PlayerSlots, _ = strconv.Atoi(Response)
		}		
		
		
		if len(MyServer.Clients) != MyServer.PlayerSlots{
			Waiting_Message := fmt.Sprintf("Waiting For Players... %d/%d Players Ready", len(MyServer.Clients), MyServer.PlayerSlots) 
			MyServer.BroadCast(Waiting_Message)
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
		NewGame := GameInit(MyServer.Clients)
		NewGame.Start_Game()
}

func GameInit(ServerClients []Client) *Game{
	NewGame := Game{Round:  0}
	for _, MyClient := range(ServerClients){
		NewPlayer := Player{Client: MyClient}
		NewGame.Players = append(NewGame.Players,  &NewPlayer)
	}

	rand.Shuffle(len(NewGame.Players), func(i, j int)  {NewGame.Players[i], NewGame.Players[j] = NewGame.Players[j], NewGame.Players[i]})

	
	if len(NewGame.Players) == 4{
		Team1 := Team{TeamPlayers: []*Player{NewGame.Players[0],  NewGame.Players[1]}, 
		TeamName: fmt.Sprintf("Team %s&%s", NewGame.Players[0].Name, NewGame.Players[1].Name), TeamPoints: 0,}
		
		Team2 := Team{TeamPlayers: []*Player{NewGame.Players[2], NewGame.Players[3]}, 
		TeamName: fmt.Sprintf("Team %s&%s", NewGame.Players[2].Name, NewGame.Players[3].Name), TeamPoints: 0}


		NewGame.Players[0].MyTeam, NewGame.Players[1].MyTeam = &Team1,&Team1
		NewGame.Players[2].MyTeam, NewGame.Players[3].MyTeam = &Team2,&Team2
		NewGame.Teams = []*Team{&Team1, &Team2}
	} else {

		Team1 := Team{TeamPlayers: []*Player{NewGame.Players[0]}, 
		TeamName: fmt.Sprintf("Team %s", NewGame.Players[0].Name), TeamPoints: 0}

		Team2 := Team{TeamPlayers: []*Player{NewGame.Players[1]}, 
		TeamName: fmt.Sprintf("Team %s", NewGame.Players[1].Name), TeamPoints: 0}


		NewGame.Players[0].MyTeam = &Team1
		NewGame.Players[1].MyTeam = &Team2
		NewGame.Teams = []*Team{&Team1, &Team2}
	}

	return &NewGame

}

func (G *Game) Start_Game(){
	RoundPlayingOrder, RoundNameOrder := G.PlayingOrder(len(G.Players))
	InternalOrder := make([]*Player, len(RoundPlayingOrder))
	copy(InternalOrder, RoundPlayingOrder)
	G.BroadCast("The Teams are")
	G.BroadCast(G.Teams[0].TeamName + " x " + G.Teams[1].TeamName)

	
	for G.Teams[0].TeamPoints < MAXPOINTS && G.Teams[1].TeamPoints < MAXPOINTS{
		G.MatchINIT(RoundNameOrder)
		Cards := G.ShuffleDeck()
		G.DealCards(Cards)

		for G.Round <= NROUNDS && G.Teams[0].RoundsWon < MAXWONROUND && G.Teams[1].RoundsWon < MAXWONROUND{
			G.NextGui()
			G.ClearRound()
			InternalOrder = G.PlayRound(InternalOrder)
			if G.Teams[0].Resigned || G.Teams[1].Resigned{
				break
			}
		}

		
		if G.Teams[0].Resigned{

			G.Teams[1].TeamPoints += G.PointsOnWin
			G.BroadCast(G.Teams[1].TeamName + " Won the Game")

		}else if G.Teams[1].Resigned{

			G.Teams[0].TeamPoints += G.PointsOnWin
			G.BroadCast(G.Teams[0].TeamName + " Won the Game")

		}else{
			if G.Teams[0].RoundsWon > G.Teams[1].RoundsWon{

				G.BroadCast(G.Teams[0].TeamName + " Won the Game")
				G.Teams[0].TeamPoints += G.PointsOnWin

			}else if G.Teams[0].RoundsWon < G.Teams[1].RoundsWon{

				G.BroadCast(G.Teams[1].TeamName + " Won the Game")
				G.Teams[1].TeamPoints += G.PointsOnWin

			}else{
				G.BroadCast("Draw")
			}
		}
		RoundPlayingOrder = append(RoundPlayingOrder[1:], RoundPlayingOrder[0])
	}

	if G.Teams[0].TeamPoints > G.Teams[1].TeamPoints{
		G.BroadCast(G.Teams[0].TeamName + " Won the Match")
	}else{
		G.BroadCast(G.Teams[1].TeamName + " Won The Match")
	}

	G.BroadCast("Thank You For Playing")
	
}


func (S *ServerStruct) BroadCast(message string){
	for _, Client := range(S.Clients){
		Client.IpAddress.Write([]byte("\n" + message))
	}
}

func (G *Game) BroadCast(message string){
	for _, Player := range(G.Players){
		Player.IpAddress.Write([]byte("\n" + message))
	}
}

func (G *Game) ListenToMe(MyPlayer *Player){
	mybuff := make([]byte, 1024)

	for{
		n, _ := MyPlayer.IpAddress.Read(mybuff)
		Message := string(mybuff[:n])
		FormattedMessage := strings.TrimSpace(strings.ToLower(Message))

		if MyPlayer.MyTeam.Challenged{
			var Response string
			ResponseBuff := make([]byte, 1024)

			sz, _ := MyPlayer.IpAddress.Read(ResponseBuff)
			Response =  strings.TrimSpace(strings.ToLower(string(ResponseBuff[:sz])))

			for Response != "y" && Response != "n"{
				MyPlayer.IpAddress.Write([]byte("Invalid Response"))
				sz, _ := MyPlayer.IpAddress.Read(ResponseBuff)
				Response =  strings.TrimSpace(strings.ToLower(string(ResponseBuff[:sz])))
			}

			MyPlayer.MyTeam.Accepted = Response
			

		}else if !MyPlayer.IsTurn{
			if n > 0 {
				G.BroadCast(MyPlayer.Name + ": " + Message)
			}

		}else if MyPlayer.IsTurn{
			switch FormattedMessage{
				case "jogar":
					CardIndex := G.Jogar(MyPlayer)

					PlayedCard := MyPlayer.CurHand[CardIndex]
					MyPlayer.CurHand = append(MyPlayer.CurHand[:CardIndex], MyPlayer.CurHand[CardIndex+1:]...)
					G.CardsOnTable = append(G.CardsOnTable, PlayedCard)
					MyPlayer.Played = true


				case "truco":
					OpponentTeam := G.Challenge(FormattedMessage, MyPlayer)

					if OpponentTeam.Accepted == "y"{
						G.PointsOnWin = 3
						G.BroadCast("Truco Aceito, A rodada Vale 3 Pontos")
						G.Jogar(MyPlayer)

					}else{
						G.BroadCast("Truco Negado")
						OpponentTeam.Resigned = true
					}

					OpponentTeam.Challenged = false
					OpponentTeam.Accepted = ""
					MyPlayer.Played = true

				case "envido":
					OpponentTeam := G.Challenge(FormattedMessage, MyPlayer)

					if OpponentTeam.Accepted == "y"{
						G.BroadCast("Envido Aceito")
						G.Envido()

					}else{
						G.BroadCast("Envido Negado")
						MyPlayer.MyTeam.TeamPoints += 1
					}

					OpponentTeam.Challenged = false
					OpponentTeam.Accepted = ""

				case "queimar":
					G.Jogar(MyPlayer)
					G.CardsOnTable[len(G.CardsOnTable)-1] = cardpack.Card{Name: "Queimar", Value: 0, Repr: cardpack.QueimadoCard}
					MyPlayer.Played = true

				case "correr":
					G.BroadCast(MyPlayer.MyTeam.TeamName + "Resigned")
					MyPlayer.MyTeam.Resigned = true
					MyPlayer.Played = true

				case "flor":
					Opponent := G.Flor(MyPlayer)
					Opponent.Accepted = ""
					Opponent.Challenged = false
			}
		}
	}
}

func (G *Game) Challenge(Challenge string, Challenger *Player) Team{

	G.BroadCast(fmt.Sprintf("%s PEDIU %s NEWBA", Challenger.Name, Challenge))
	var OpponentTeam Team
	if Challenger.MyTeam == G.Teams[0]{
		OpponentTeam = *G.Teams[1]
	}else{ 							  
		OpponentTeam = *G.Teams[0]
	}
				
	OpponentTeam.Challenged = true
	for _, Player := range(OpponentTeam.TeamPlayers){
		Player.IpAddress.Write([]byte("VAI ACEITAR (y/n)"))
	}

	G.BroadCast("Waiting for " + OpponentTeam.TeamName)
	for OpponentTeam.Accepted != "y" && OpponentTeam.Accepted != "n"{
		time.Sleep(1 * time.Second)
	}

	return OpponentTeam
}

func (G *Game) Jogar(MyClient *Player) int{
	CardBuff := make([]byte, 1024)
	MyClient.IpAddress.Write([]byte("\n" +"Enter the index of your card (1-3)"))

	sz, _ := MyClient.IpAddress.Read(CardBuff)
	Num, _ := strconv.Atoi(strings.TrimSpace(string(CardBuff[:sz])))
	CardIndex := Num-1

	for CardIndex > len(MyClient.CurHand) || CardIndex < 0{
		MyClient.IpAddress.Write([]byte("\n" +"Invalid Index"))
		MyClient.IpAddress.Write([]byte("\n" +"Enter the index of your card (1-3)"))
		sz, _ := MyClient.IpAddress.Read(CardBuff)
		Num, _ := strconv.Atoi(string(CardBuff[:sz]))
		CardIndex = Num-1
	}

	return CardIndex 
}

func (G *Game) Envido(){
	HighestInvido := [][]int{{0,0}}
	for PlayerIdx, Client := range(G.Players){
		MyCards := Client.CurHand
		Suits0 := []rune(MyCards[0].Name)[1]
		Suits1 := []rune(MyCards[1].Name)[1]
		Suits2 := []rune(MyCards[2].Name)[1]
		
		MyPoints := []int{}
		if Suits0 == Suits1 && Suits1 == Suits2{
			Values := []int{cardpack.EnvidoValues[MyCards[0].Value], cardpack.EnvidoValues[MyCards[1].Value], cardpack.EnvidoValues[MyCards[2].Value]}
			sort.Ints(Values)
			MyPoints = []int{Values[0]+Values[1], PlayerIdx}
			
		}else if Suits0 == Suits1{
			MyPoints = []int{cardpack.EnvidoValues[MyCards[0].Value] + cardpack.EnvidoValues[MyCards[1].Value], PlayerIdx }

		}else if Suits0 == Suits2{
			MyPoints = []int{cardpack.EnvidoValues[MyCards[0].Value] + cardpack.EnvidoValues[MyCards[2].Value], PlayerIdx} 

		}else if Suits1 == Suits2{
			MyPoints = []int{cardpack.EnvidoValues[MyCards[1].Value] + cardpack.EnvidoValues[MyCards[2].Value], PlayerIdx}
		}else{
			Values := []int{cardpack.EnvidoValues[MyCards[0].Value], cardpack.EnvidoValues[MyCards[1].Value], cardpack.EnvidoValues[MyCards[2].Value]}
			sort.Ints(Values)
			MyPoints = []int{Values[0], PlayerIdx}
		}

		if MyPoints[0] > HighestInvido[0][0]{
			HighestInvido[0] = MyPoints
		}
	}

	WinnerIDX := HighestInvido[0][1]
	Winner := G.Players[WinnerIDX]
	Winner.MyTeam.TeamPoints += 2
	G.BroadCast(Winner.MyTeam.TeamName + "Won")
}

func (G *Game) Flor(MyPlayer *Player) Team{
	MyHand := MyPlayer.CurHand
	suit1 := []rune(MyHand[0].Name)[1]
	suit2 := []rune(MyHand[1].Name)[1]
	suit3 := []rune(MyHand[2].Name)[1]

	var OpponentTeam Team
	if suit1 == suit2 && suit1 == suit3{
		OpponentTeam = G.Challenge("flor", MyPlayer)
	}else{
		MyPlayer.IpAddress.Write([]byte("Not a Flor"))
		return OpponentTeam
	}

	if OpponentTeam.Accepted == "y"{
		var P0Value int
		var P1Value int

		for _, Players := range(MyPlayer.MyTeam.TeamPlayers){
			for _, Card := range(Players.CurHand){
				if Card.Value >= 10{
					P0Value += 1
				}else{
					P0Value += Card.Value
				}
			}
		}

		for _, Players := range(OpponentTeam.TeamPlayers){
			for _, Card := range(Players.CurHand){
				if Card.Value >= 10{
					P1Value += 1
				}else{
					P1Value += Card.Value
				}
			}	
		}

		if P1Value > P0Value{
			OpponentTeam.TeamPoints += 6
			G.BroadCast(OpponentTeam.TeamName + "Won")
		}else if P0Value > P1Value{
			OpponentTeam.TeamPoints += 3
			G.BroadCast(MyPlayer.MyTeam.TeamName+ "Won")
		}
		OpponentTeam.Challenged = false
		OpponentTeam.Accepted = ""
	

	}else{
		G.BroadCast("Flor Negada")
		MyPlayer.MyTeam.TeamPoints += 3
	}

	return OpponentTeam
}


func (G  *Game) MatchINIT(RoundNameOrder []string){
	G.ClearGameVariables()
	G.BroadCast("Round: " + string(G.Round))
	G.BroadCast(G.Teams[0].TeamName + ": " + string(G.Teams[0].TeamPoints))
	G.BroadCast(G.Teams[1].TeamName + ": " +  string(G.Teams[1].TeamPoints))
	G.BroadCast("Round Order: " + strings.Join(RoundNameOrder, ","))

	for _, Players := range(G.Players){
		go G.ListenToMe(Players)
	}
}

func (G * Game) PlayRound(RoundPlayingOrder []*Player) []*Player{

	for _, CurPlayer := range(RoundPlayingOrder){
		CurPlayer.IsTurn = true
		G.BroadCast(CurPlayer.Name + "'s Turn")
		CurPlayer.IpAddress.Write([]byte("\n" + "It's Your Turn!"))
		

		for !CurPlayer.Played{
			time.Sleep(5 * time.Second)
			//G.BroadCast(fmt.Sprintf("\n Waiting For %s", CurPlayer.Name))
		}


		if CurPlayer.MyTeam.Resigned{
			G.BroadCast(CurPlayer.MyTeam.TeamName + "Has Resigned")
			return nil
		}else{
			G.BroadCast(CurPlayer.Name + "Has Played: ")
			PlayedCard := cardpack.CreateTerminalRepr(G.CardsOnTable[len(G.CardsOnTable)-1].Name)
			for line := range(7){
				G.BroadCast(PlayedCard[line])
			}
		}
	}

	HighestCard := [][]int{{0,0}}
	for PlayerIdx, Card := range(G.CardsOnTable){
		CardVal := cardpack.Values[Card.Name]
		if CardVal > HighestCard[0][0]{
				HighestCard = [][]int{{CardVal,  PlayerIdx}}
		}else if CardVal == HighestCard[0][0]{
				HighestCard = append(HighestCard, []int{CardVal,  PlayerIdx})
		}
	}

	WinnerIndex := HighestCard[0][1]
	WinnerPlayer := RoundPlayingOrder[WinnerIndex]
	WinnerTeam := WinnerPlayer.MyTeam

	if len(HighestCard) == 1{
		G.BroadCast(WinnerPlayer.Name + " Won The Round")
		RoundPlayingOrder = append(RoundPlayingOrder[WinnerIndex:], RoundPlayingOrder[:WinnerIndex]...)
		WinnerTeam.RoundsWon += 1
	}else{
		var WonRound int
		for _, CardValues := range(HighestCard){
			WonRound = 1
			Player := CardValues[1]
			RoundPlayingOrder = append(RoundPlayingOrder[WinnerIndex:], RoundPlayingOrder[:WinnerIndex]...)
			if RoundPlayingOrder[Player].MyTeam != WinnerTeam{
				G.BroadCast("Draw")
				WonRound = 0
				break
			}
		}	
			
		WinnerTeam.RoundsWon += WonRound
	}

	fmt.Println(RoundPlayingOrder)
	return RoundPlayingOrder
}

func (G * Game) ShuffleDeck() []cardpack.Card{

	Hands := []cardpack.Card{}
	
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cardpack.Cards), func(i, j int)  {cardpack.Cards[i], cardpack.Cards[j] = cardpack.Cards[j], cardpack.Cards[i]})

	for i := 0; i < 16; i++{
		Hands = append(Hands, cardpack.Card{Name: cardpack.Cards[i], Value: cardpack.Values[cardpack.Cards[i]], Repr: cardpack.CreateTerminalRepr(cardpack.Cards[i])})
	}

	return Hands
}

func (G *Game) DealCards(Cards []cardpack.Card){
	
	for _, Player := range(G.Players){
		Player.CurHand = append(Player.CurHand, Cards[0], Cards[1], Cards[2])
		Cards = Cards[3:]
	}

}

func (G *Game) NextGui(){

	for _, Player := range(G.Players){
		PlayerGUI := cardpack.UpdateGui(G.Round, Player.CurHand)

		Player.IpAddress.Write([]byte("\n"))
		for line := range(18){
			Player.IpAddress.Write([]byte(PlayerGUI[line] + "\n"))
		}
	}
}

func (G * Game) PlayingOrder(PlayerCount int) ([]*Player, []string){

	var RoundPlayingOrder []*Player
	var RoundNameOrder []string

	switch PlayerCount{
		case 2:
			RoundPlayingOrder = []*Player{G.Teams[0].TeamPlayers[0], G.Teams[1].TeamPlayers[0]}
			RoundNameOrder = []string{G.Teams[0].TeamPlayers[0].Name, G.Teams[1].TeamPlayers[0].Name}
		case 4:
			RoundPlayingOrder = []*Player{G.Teams[0].TeamPlayers[0], G.Teams[1].TeamPlayers[0], G.Teams[0].TeamPlayers[1], G.Teams[1].TeamPlayers[1]}
			RoundNameOrder = []string{G.Teams[0].TeamPlayers[0].Name, G.Teams[1].TeamPlayers[0].Name, G.Teams[0].TeamPlayers[1].Name, G.Teams[1].TeamPlayers[1].Name}
		}

	return RoundPlayingOrder, RoundNameOrder
}

func (G *Game) ClearGameVariables(){
	G.PointsOnWin = 1
	G.Round = 1
	G.CardsOnTable = []cardpack.Card{}
}

func (G *Game) ClearRound(){
	G.PointsOnWin = 1
	G.Round += 1
	
	for _, Client := range(G.Players){
		Client.IsTurn = false
		Client.Played = false
	}
}



