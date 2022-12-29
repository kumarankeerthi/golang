package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	port := "8080"

	// use the port value is set as  environment variable
	if envPort := os.Getenv("GOPORT"); envPort != "" {
		port = envPort

	}

	// create a server using theaddress and handler
	server := &http.Server{Addr: "0.0.0.0:" + port, Handler: service()}

	// set the context for the server start and stop

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// create a buffered channel of size 1 with os.Signal as type
	sig := make(chan os.Signal, 1)

	// signal.Notify will push the signal object to sig channel if one of this syscall is recieved.
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {

		// this goroutie will get triggered when sig channel gets a valye pushed.

		<-sig
		fmt.Println("inside first foroute--2")
		// Shutdown the server gracefully with say 30 sec dealy
		shutdownServerCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)
		go func() {
			<-shutdownServerCtx.Done()
			fmt.Println("inside first foroute")
			if shutdownServerCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit")
			}
		}()
		fmt.Println("inside first foroute- triggering now")
		//Trigger graceful shutdown
		err := server.Shutdown(shutdownServerCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	//Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()
}
func service() http.Handler {
	// create a handler  (http.handler) using chi
	r := chi.NewRouter()

	// add required middleware from chi library
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100000 * time.Millisecond)

		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte("Rest Api using go-chi libary"))
	})

	return r
}
