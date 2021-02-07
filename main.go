package main

import (
	"authservice/cmd"
	"authservice/data"
	"log"
	"os"
	"os/signal"

	"github.com/Smart-Pot/pkg"
	"github.com/Smart-Pot/pkg/adapter/amqp"
)

func main() {
	if err := pkg.Config.ReadConfig(); err != nil {
		log.Fatal(err)
	}

	data.DatabaseConnection()

	if err := amqp.Set(pkg.Config.AMQPAddress); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		if err := cmd.Execute(); err != nil {
			log.Fatal(err)
		}
	}()
	sig := <-c
	log.Println("GOT SIGNAL: " + sig.String())
}
