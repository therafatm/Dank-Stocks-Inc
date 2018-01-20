package main

import (
    "encoding/json"
    "net/http"
    "strings"
    "bufio"
    "bytes"
    "sync"
    "time"
    "log"
    "fmt"
    "os"
    "./commands"
)

func postData(client *http.Client, url string, data []byte){
    req, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Content-Type", "application/json")
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }else{
        defer resp.Body.Close()
    }
}

func postUserData(wg *sync.WaitGroup, url string, cmds []commands.Command){
    defer wg.Done()
    client := &http.Client{}
    for _, command := range cmds {
        js, err := json.Marshal(command)
        if err != nil {
            log.Fatal(err)
        }
        postData(client, url, js)
    }
}

func main() {
    file, err := os.Open("workfiles/10userWorkLoad")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    url := "http://0.0.0.0:8080"
    allCmds := make([]commands.Command, 0)

    replacer := strings.NewReplacer("[", "", "]", "", ",", " ")
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := replacer.Replace(scanner.Text())
        data := strings.Fields(line)
        command := commands.ParseData(data)
        allCmds = append(allCmds, command)
    }

    userMap := make(map[string][]commands.Command)
    for _, command := range allCmds {
        if _, exist := userMap[command.Username]; !exist {
            userMap[command.Username] = make([]commands.Command, 0)
        }
        userMap[command.Username] = append(userMap[command.Username], command)
    }

    var wg sync.WaitGroup
    var start = time.Now()

    wg.Add(len(userMap))
    for _, cmds := range userMap {
        go postUserData(&wg, url, cmds)
    }
    
    wg.Wait()
    totalTime := time.Since(start)
    fmt.Printf("%d in %s\n", len(allCmds), totalTime)

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}