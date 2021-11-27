package main

import (
	"amqprpc/amqprpc"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	client, err := amqprpc.NewClient(&amqprpc.Config{
		Dsn:               "amqp://admin:admin@localhost:5672/",
		ClientTimeout:     10,
		ReconnectInterval: 5,
		Log:               log.StandardLogger(),
	})

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		params := Args{A: 5, B: i}

		var result Result
		if err := client.Call("rpc.method.test__multiply", params, &result); err != nil {
			log.Fatal(err)
		}
		log.Println(result)
	}

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
