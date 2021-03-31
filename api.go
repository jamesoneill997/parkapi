package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jamesoneill997/parkapi/handlers"
	"github.com/jamesoneill997/parkapi/initialise"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var port = os.Getenv("PORT")

func init() {
	client = initialise.Client
	ctx = initialise.Ctx
}

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	ctx = context.Background()
	userHandler := handlers.NewUser(l)
	loginHandler := handlers.NewLogin(ctx, l, client)
	logoutHandler := handlers.NewLogout(ctx, l, client)
	carparkHandler := handlers.NewCarPark(ctx, l, client)
	vehicleHandler := handlers.NewVehicle(ctx, l, client)
	accessHandler := handlers.NewAccess(ctx, l, client)

	//creation of serve multiplexer
	sm := http.NewServeMux()
	sm.Handle("/users", userHandler)
	sm.Handle("/login", loginHandler)
	sm.Handle("/logout", logoutHandler)
	sm.Handle("/carparks", carparkHandler)
	sm.Handle("/vehicles", vehicleHandler)
	sm.Handle("/access", accessHandler)

	//server config
	server := &http.Server{
		Addr:         port,
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 4 * time.Second,
	}

	//goroutine to listen the server
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//make a channel that receives a signal
	sigChan := make(chan os.Signal)

	//notify keyboard interrupts and crashes
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	l.Println("Termination, graceful shutdown.", sig)
	timeoutCtx, err := context.WithTimeout(context.Background(), 30*time.Second)
	if err != nil {
		l.Fatal(err)
	}

	server.Shutdown(timeoutCtx)
}
