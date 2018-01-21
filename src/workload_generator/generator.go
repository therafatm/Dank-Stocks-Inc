package main

import (
    "net/http"
    "strings"
    "io/ioutil"
    "bufio"
    "sync"
    "time"
    "log"
    "fmt"
    "os"
    "./commands"
)

func postData(client *http.Client, url string){
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Connection", "keep-alive")
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Sent request %s\n", url)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Body: %s\n", body)

}

func postUserData(wg *sync.WaitGroup, url string, cmds []commands.Command){
    defer wg.Done()
    client := &http.Client{}
    for _, command := range cmds {
        endpoint := commands.FormatCommandEndpoint(command)
        if endpoint != "" {
            postData(client, url + endpoint)
        }
    }
}

func main() {
    file, err := os.Open("../../workfiles/10userWorkLoad")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    url := "http://0.0.0.0:8888"
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