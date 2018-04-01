package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"test/workload_generator/commands"
	"time"
	"runtime"
)

const tlen time.Duration = time.Duration(time.Millisecond * 1000)

func postData(client *http.Client, url string) (err error, status int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	req.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	status = resp.StatusCode
	io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	return
}

func postUserData(wg *sync.WaitGroup, client *http.Client, url string, cmds []commands.Command) {
	defer wg.Done()
	for _, command := range cmds {
		endpoint := commands.FormatCommandEndpoint(command)
		if endpoint != "" {
			_, status := postData(client, url+endpoint)
			for status == 404 {
				_, status = postData(client, url+endpoint)
			}
		}
	}
}

func main() {
	var host = flag.String("host", "transaction", "hostname of target")
	var port = flag.String("port", "8888", "port to target")
	var filename = flag.String("filepath", "workfiles/10userWorkLoad", "path to workload file")
	flag.Parse()

	runtime.GOMAXPROCS(4)

	// log.SetOutput(ioutil.Discard)

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

	tr := &http.Transport{
		MaxIdleConns:       50,
		IdleConnTimeout:    30 * time.Second,
	}

	client := &http.Client{
		Timeout: tlen,
		Transport: tr,
	}

	for _, cmds := range userMap {
		go postUserData(&wg, client, url, cmds)
	}
	wg.Wait()

	for _, cmd := range otherCmds {
		endpoint := commands.FormatCommandEndpoint(cmd)
		if endpoint != "" {
			_, status := postData(client, url+endpoint)
			for status == 404 {
				_, status = postData(client, url+endpoint)
			}
		}
	}

	totalTime := time.Since(start)

	log.Print("%d in %s\n", len(allCmds), totalTime)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
