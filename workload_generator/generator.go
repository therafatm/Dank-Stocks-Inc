package main

import (
    "bufio"
    "strings"
    "fmt"
    "log"
    "os"
)

func parseLine() {

}

func main() {
    file, err := os.Open("workfiles/1userWorkLoad")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    replacer := strings.NewReplacer("[", "", "]", "", ",", " ")
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := replacer.Replace(scanner.Text())
        data := strings.Split(line, " ")
        fmt.Println()
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}