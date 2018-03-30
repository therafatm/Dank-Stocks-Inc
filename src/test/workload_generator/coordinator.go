package main

import (
	"os"
	//"fmt"
	"log"
	"flag"
	"bufio"
	"strings"
	"test/workload_generator/commands"
	"test/workload_generator/workcom"
)


func main() {
	var host = flag.String("host", "transaction", "hostname of target")
	var port = flag.String("port", "8888", "port to target")
	var filename = flag.String("filepath", "workfiles/10userWorkLoad", "path to workload file")
	flag.Parse()

	wconn := workcom.NewWorkloadConnection()

	// log.SetOutput(ioutil.Discard)

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//url := fmt.Sprintf("http://%s:%s", *host, *port)
	users := make(map[string]bool)
	otherCmds := make([]commands.Command, 0)

	replacer := strings.NewReplacer("[", "", "]", "", ".", "", ",", " ")
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := replacer.Replace(scanner.Text())
		data := strings.Fields(line)
		command := commands.ParseData(data)
		if len(command.Username) > 0 {
			if _, exist := users[command.Username]; !exist{
				users[command.Username] = true
				route := workcom.UserRoute{ Username: command.Username, Host: *host, Port: *port}
				wconn.PublishUserRoute(route)
			}
			wconn.PublishCommand(command)
		} else {
			otherCmds = append(otherCmds, command)
		}
	}
}