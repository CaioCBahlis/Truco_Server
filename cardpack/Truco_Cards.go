package cardpack




type Card struct{
	Name string
	Value int
}


var Values = map[string]int{
    "4C": 1, "4O": 1, "4P": 1, "4E": 1,
    "5C": 2, "5O": 2, "5P": 2, "5E": 2,
    "6C": 3, "6O": 3, "6P": 3, "6E": 3,
    "7C": 4, "7P": 4,
    "JC": 5, "JO": 5, "JP": 5, "JE": 5,
    "QC": 6, "QO": 6, "QP": 6, "QE": 6,
    "KC": 7, "KO": 7, "KP": 7, "KE": 7,
    "1C": 8, "1O": 8,
    "2C": 9, "2O": 9, "2P": 9, "2E": 9,
    "3C": 10, "3O": 10, "3P": 10, "3E": 10,
    "7O": 11, "7E": 12, "1P": 13, "1E": 14,
}


var Cards = []string{ "4C", "4O", "4P", "4E", "5C", "5O", "5P", "5E", "6C", "6O", "6P", "6E", 
				  "7C", "7P", "JC", "JO", "JP", "JE", "QC", "QO", "QP", "QE", "KC", "KO", "KP", "KE", 
                  "1C", "1O", "2C", "2O", "2P", "2E", "3C", "3O", "3P", "3E", "7O", "7E", "1P", "1E",
}

																			
