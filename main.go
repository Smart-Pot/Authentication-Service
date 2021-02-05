package main

import (
	"authservice/cmd"
	"authservice/config"
	"authservice/data"
	"log"
	"os"
	"os/signal"
)

func main() {
	config.ReadConfig()

	data.DatabaseConnection()

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
