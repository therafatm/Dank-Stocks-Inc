package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"test/workload_generator/commands"
	"time"
)

const tlen time.Duration = time.Duration(5 * time.Second)

func postData(client *http.Client, url string) (err error, status int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	req.Header.Set("Connection", "close")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	log.Println("Sent request %s\n", url)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	status = resp.StatusCode
	log.Println("Body: %s\n", body)
	return
}

func postUserData(wg *sync.WaitGroup, url string, cmds []commands.Command) {
	defer wg.Done()
	client := &http.Client{
		Timeout: tlen,
	}
	for _, command := range cmds {
		endpoint := commands.FormatCommandEndpoint(command)
		if endpoint != "" {
			err, status := postData(client, url+endpoint)
			for err != nil || status == 404 {
				err, status = postData(client, url+endpoint)
			}
		}
	}
}

func main() {
	var host = flag.String("host", "192.168.1.150", "hostname of target")
	var port = flag.String("port", "8888", "port to target")
	var filename = flag.String("filepath", "../workfiles/10userWorkLoad", "path to workload file")
	flag.Parse()

	log.SetOutput(ioutil.Discard)

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	url := fmt.Sprintf("http://%s:%s", *host, *port)
	allCmds := make([]commands.Command, 0)
	otherCmds := make([]commands.Command, 0)

	replacer := strings.NewReplacer("[", "", "]", "", ".", "", ",", " ")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := replacer.Replace(scanner.Text())
		data := strings.Fields(line)
		command := commands.ParseData(data)
		if len(command.Username) > 0 {
			allCmds = append(allCmds, command)
		} else {
			otherCmds = append(otherCmds, command)
		}
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

	client := &http.Client{
		Timeout: tlen,
	}
	for _, cmd := range otherCmds {
		endpoint := commands.FormatCommandEndpoint(cmd)
		if endpoint != "" {
			err, status := postData(client, url+endpoint)
			for err != nil || status == 404 {
				err, status = postData(client, url+endpoint)
			}
		}
	}

	totalTime := time.Since(start)

	fmt.Println("%d in %s\n", len(allCmds), totalTime)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
