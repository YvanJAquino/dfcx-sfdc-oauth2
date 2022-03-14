package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/YvanJAquino/dfcx-sfdc-oauth2/handlers"
	"github.com/YvanJAquino/dfcx-sfdc-oauth2/helpers"
	redis "github.com/go-redis/redis/v8"
)

var (
	parent = context.Background()
	PORT   = helpers.GetEnvDefault("PORT", "8081")
	opts   = &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	rdb = redis.NewClient(opts)
)

func main() {

	notify, stop := signal.NotifyContext(parent, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/health-check", handlers.HealthCheckHandler)
	mux.HandleFunc("/generate-login", handlers.GenerateLoginHandler)
	mux.HandleFunc("/callback", handlers.CallbackHandler)
	mux.Handle("/alt-generate-login",
		handlers.NewGenerateLoginHandle(rdb))

	server := &http.Server{
		Addr:        ":" + PORT,
		Handler:     mux,
		BaseContext: func(net.Listener) context.Context { return parent },
	}
	fmt.Println("Listening and serving on :" + PORT)
	go server.ListenAndServe()
	<-notify.Done()
	fmt.Println("Gracefully shutting down the HTTP/S server")
	shutdown, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()
	server.Shutdown(shutdown)
}
