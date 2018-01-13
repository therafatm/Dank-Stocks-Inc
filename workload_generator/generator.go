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
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    req.Header.Set("Connection", "close")
    req.Header.Set("Content-Type", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("%s\n", err)
    }
    defer resp.Body.Close()
}

func main() {
    file, err := os.Open("workfiles/45User_testWorkLoad")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    url := "http://localhost:8080"
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
        go postData(client, url, js)
    }


    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}