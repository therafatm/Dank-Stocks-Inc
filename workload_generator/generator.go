package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "bufio"
    "bytes"
    "log"
    "os"
    "./commands"
)

func postData(client *http.Client, url string, data []byte){
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}

func main() {
    file, err := os.Open("workfiles/100User_testWorkLoad/data")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    url := "http://127.0.0.1:8080"
    client := &http.Client{}
    cmds := make([]commands.Command, 0)

    replacer := strings.NewReplacer("[", "", "]", "", ",", " ")
    scanner := bufio.NewScanner(file)


    for scanner.Scan() {
        line := replacer.Replace(scanner.Text())
        data := strings.Fields(line)
        command := commands.ParseData(data)
        cmds = append(cmds, command)
    }

    for _, command := range cmds {
        js, err := json.Marshal(command)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("%s", js)
        postData(client, url, js)
    }


    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}