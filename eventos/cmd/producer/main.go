package main

import (
	"fmt"

	"github.com/samuelmarscos/eventos/pkg/rabbitmq"
)

func main() {

	ch, err := rabbitmq.OpenChannel()

	if err != nil {
		panic(err)
	}
	defer ch.Close()

	for i := 0; i < 100000; i++ {
		msg := fmt.Sprintf("Esse é o %d Hello World !", i)
		rabbitmq.Publish(ch, msg, "amq.direct")
	}

}
