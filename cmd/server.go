package cmd

import (
	"authservice/endpoints"
	"authservice/service"
	"authservice/transport"
	"fmt"
	golog "log"
	"net/http"
	"os"
	"time"

	"github.com/Smart-Pot/pkg"
	"github.com/Smart-Pot/pkg/adapter/amqp"
	"github.com/go-kit/kit/log"
)

func startServer() error {
	// Defining Logger
	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	// Make producer for service
	verifyProducer, err := amqp.MakeProducer("VerificationMail")
	forgotProducer, err := amqp.MakeProducer("ForgotPassword")
	if err != nil {
		return err
	}

	service := service.NewService(logger, verifyProducer,forgotProducer)
	endpoint := endpoints.MakeEndpoints(service)
	handler := transport.MakeHTTPHandlers(endpoint, logger)

	fmt.Println("HELLO",pkg.Config.Server.Address)

	l := golog.New(os.Stdout, "AUTH-SERVICE", 0)
	// Set handler and listen given port
	s := http.Server{
		Addr:         pkg.Config.Server.Address, // configure the bind address
		Handler:      handler,                   // set the default handler
		ErrorLog:     l,                         // set the logger for the server
		ReadTimeout:  5 * time.Second,           // max time to read request from the client
		WriteTimeout: 10 * time.Second,          // max time to write response to the client
		IdleTimeout:  120 * time.Second,         // max time for connections using TCP Keep-Alive
	}

	return s.ListenAndServe()
}
