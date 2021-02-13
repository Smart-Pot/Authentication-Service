package main

import (
	"authservice/cmd"
	"authservice/data"
	"authservice/service/oauth"
	"log"
	"os"
	"os/signal"
	"path/filepath"

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

	wd,_ := os.Getwd()
	if err := oauth.ReadConfig(filepath.Join(wd,"config"));err != nil {
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
