package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"hash/fnv"
	"os"
	"strings"
	"sync"
	"test/workload_generator/commands"
	"time"
	"runtime"
	"github.com/valyala/fasthttp"
)

const tlen time.Duration = time.Duration(time.Millisecond * 1000)

func postData(client *fasthttp.Client, url string) (err error, status int) {
	req := fasthttp.AcquireRequest()
    req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	err = client.Do(req, resp)
	if err != nil {
		log.Println(err.Error())
		return
	}
	status = 0 // resp.StatusCode TODO: get actual?
	// io.Copy(ioutil.Discard, resp.Body())
	// resp.Body.Close()
	return
}

func postUserData(wg *sync.WaitGroup, url string, cmds []commands.Command) {
	defer wg.Done()
	client := &fasthttp.Client{
		MaxIdleConnDuration: 30 * time.Second,
		MaxConnsPerHost: 1,
	}

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

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}

func main() {
	var host = flag.String("host", "transaction", "hostname of target")
	var port = flag.String("port", "8888", "port to target")
	var filename = flag.String("filepath", "workfiles/10userWorkLoad", "path to workload file")
	var max_goroutines = *flag.Int("go", 10, "number of go routines")
	flag.Parse()

	runtime.GOMAXPROCS(1)

	//log.SetOutput(ioutil.Discard)

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

	userMap := make(map[int][]commands.Command)
	for _, command := range allCmds {
		key := hash(command.Username) % max_goroutines
		if _, exist := userMap[key]; !exist {
			userMap[key] = make([]commands.Command, 0)
		}
		userMap[key] = append(userMap[key], command)
	}

	var wg sync.WaitGroup
	var start = time.Now()

	wg.Add(len(userMap))

	for _, cmds := range userMap {
		go postUserData(&wg, url, cmds)
	}
	wg.Wait()

	client := &fasthttp.Client{}
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

	log.Printf("%d in %s\n", len(allCmds), totalTime)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
