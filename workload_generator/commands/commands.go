package commands

import (
    "strconv"
    "log"
)


type Add struct {
    Command string
    Username string
    Amount float64
}

func ParseAdd(data []string) Add {
    amount, err := strconv.ParseFloat(data[3], 64)
    if err != nil{
        log.Fatal("Could not parse Amount")
    }
    return Add{Command: data[1], Username: data[2], Amount: amount}
}

type Add struct {
    Command string
    Username string
    Amount float64
}