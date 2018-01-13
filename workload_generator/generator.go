package main

import (
    "encoding/json"
    "bufio"
    "strings"
    "log"
    "os"
    "./commands"
)

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
        command := commands.ParseData(data)
        cmds = append(cmds, command)
        js, err := json.Marshal(command)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("%s", js)
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}