package main

import (
    "bufio"
    "strings"
    "fmt"
    "log"
    "os"
    "./commands"
)


func ParseData(data []string) commands.Command {
    if len(data) < 2{
        log.Fatal("Length of entry too short")
    }
    var command commands.Command

    switch cmdName := data[1]; cmdName {
        case "ADD":
            command = commands.ParseAdd(data)

        case "QUOTE":
            command = commands.ParseQuote(data)

        case "BUY":
            command = commands.ParseBuy(data)

        case "SELL":
            command = commands.ParseSell(data)


        default:
            log.Fatal("Invalid command")
    }

    fmt.Println(command)
    return command
}

func main() {
    file, err := os.Open("workfiles/1userWorkLoad")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    cmds := make([]commands.Command, 0)

    replacer := strings.NewReplacer("[", "", "]", "", ",", " ")
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := replacer.Replace(scanner.Text())
        data := strings.Fields(line)
        command := ParseData(data)
        cmds = append(cmds, command)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}