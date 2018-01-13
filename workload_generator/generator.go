package main

import (
    "bufio"
    "strings"
    "fmt"
    "log"
    "os"
    "./commands"
)


func ParseData(data []string) interface{}{
    if len(data) < 2{
        log.Fatal("Length of entry too short")
    }
    var command interface{}

    switch cmdName := data[1]; cmdName {
        case "ADD":
            command = commands.ParseAdd(data)

        default:
            log.Fatal("Invalid command")
    }

    return command
}

func main() {
    file, err := os.Open("workfiles/1userWorkLoad")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    cmds := make([]interface{}, 0)

    replacer := strings.NewReplacer("[", "", "]", "", ",", " ")
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := replacer.Replace(scanner.Text())
        data := strings.Fields(line)
        command := ParseData(data)
        cmds = append(cmds, command)
        fmt.Println(cmds)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}