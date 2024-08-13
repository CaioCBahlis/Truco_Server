package cardpack




type Card struct{
	Name string
	Value int
    Repr []string
}

var CardRepr = []string{
    ".-------.",
    "|   X   |",
    "|       |",
    "|   ♠   |",
    "|       |",
    "|   X   |",
    "'-------'",
}




var Values = map[string]int{
    "4♣": 1, "4♥": 1, "4♦": 1, "4♠": 1,
    "5♣": 2, "5♥": 2, "5♦": 2, "5♠": 2,
    "6♣": 3, "6♥": 3, "6♦": 3, "6♠": 3,
    "7♣": 4, "7♦": 4,
    "J♣": 5, "J♥": 5, "J♦": 5, "J♠": 5,
    "Q♣": 6, "Q♥": 6, "Q♦": 6, "Q♠": 6,
    "K♣": 7, "K♥": 7, "K♦": 7, "K♠": 7,
    "A♣": 8, "A♥": 8,
    "2♣": 9, "2♥": 9, "2♦": 9, "2♠": 9,
    "3♣": 10, "3♥": 10, "3♦": 10, "3♠": 10,
    "7♥": 11, "7♠": 12, "A♦": 13, "A♠": 14,
}

var Cards = []string{
    "4♣", "4♥", "4♦", "4♠", "5♣", "5♥", "5♦", "5♠", "6♣", "6♥", "6♦", "6♠",
    "7♣", "7♦", "J♣", "J♥", "J♦", "J♠", "Q♣", "Q♥", "Q♦", "Q♠", "K♣", "K♥", "K♦", "K♠",
    "A♣", "A♥", "2♣", "2♥", "2♦", "2♠", "3♣", "3♥", "3♦", "3♠", "7♥", "7♠", "A♦", "A♠",
}


func CreateTerminalRepr(Card string) []string{
    NewCard := make([]string, len(CardRepr))
    copy(NewCard, CardRepr)


    NewCard[1] = NewCard[1][0:4] + string(Card[0]) + NewCard[1][5:]
    NewCard[3] = NewCard[2][0:4] + string([]rune(Card)[1]) + NewCard[2][5:]
    NewCard[5] = NewCard[5][0:4] + string(Card[0]) + NewCard[5][5:]
    return NewCard
}