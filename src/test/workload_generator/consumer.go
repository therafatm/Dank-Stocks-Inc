package main

import (
	"test/workload_generator/workcom"
)


func main() {
	wconn := workcom.NewWorkloadConnection()
	nameq := wconn.NameQueue()
	msgs, err := wconn.Consume(
		nameq.Name,
		"",
		true,
		false,
		false,
		false,
		nil)

	forever := make(chan bool)

    go func() {
            for d := range msgs {
                log.Printf(" [x] %s", d.Body)
            }
    }()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}