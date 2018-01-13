package commands

import (
    "strconv"
    "log"
)

type Command struct {
    Name string
    Username string
    Symbol string
    Amount float64
}

func ParseAdd(data []string) Command {
    amount, err := strconv.ParseFloat(data[3], 64)
    if err != nil{
        log.Fatal("Could not parse Amount")
    }
    return Command{Name: data[1], Username: data[2], Amount: amount}
}


func ParseQuote(data []string) Command {
    return Command{Name: data[1], Username: data[2], Symbol: data[3]}
}


func ParseBuy(data []string) Command {
    amount, err := strconv.ParseFloat(data[4], 64)
    if err != nil{
        log.Fatal("Could not parse Amount")
    }
    return Command{Name: data[1], Username: data[2], Symbol: data[3], Amount: amount}
}


func ParseSell(data []string) Command {
    amount, err := strconv.ParseFloat(data[4], 64)
    if err != nil{
        log.Fatal("Could not parse Amount")
    }
    return Command{Name: data[1], Username: data[2], Symbol: data[3], Amount: amount}
}