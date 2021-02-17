package main

import (
	"authservice/cmd"
	"authservice/service/oauth"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/Smart-Pot/pkg"
	"github.com/Smart-Pot/pkg/adapter/amqp"
	"github.com/Smart-Pot/pkg/db"
)

func main() {
	if err := pkg.Config.ReadConfig(); err != nil {
		log.Fatal(err)
	}
	log.Println("Configurations are set")
	
	if err := db.Connect(db.PkgConfig("users")); err != nil {
		log.Fatal(err)
	}
	log.Println("DB Connection established")

	if err := amqp.Set(pkg.Config.AMQPAddress); err != nil {
		log.Fatal(err)
	}
	log.Println("AMQP module is set")

	wd, _ := os.Getwd()
	if err := oauth.ReadConfig(filepath.Join(wd, "config")); err != nil {
		log.Fatal(err)
	}
	log.Println("OAuth module is set")

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
