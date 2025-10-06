package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/d1vyanshu-kumar/students-api/internal/config"
	student "github.com/d1vyanshu-kumar/students-api/internal/http/handlers/students"
	"github.com/d1vyanshu-kumar/students-api/internal/storage/sqlite"
)

// setup coustome logger and we are going to use inbuilt log package so we dont need to setup any coustome logger
func main() {
	// load config
	// here is the first step we need to load the config

	cfg := config.MustLoad()

	// database setup
	

	storage, er:= sqlite.New(cfg)

	if er != nil {
		fmt.Println(er)
		return
	}

	slog.Info("database connected", slog.String("env", cfg.Env), slog.String("version", "1.0.0")) // we are using structured log here

	

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage)) // we make plural because we are going to have multiple students.
	// and in near future if we want to add a new dependecy  we can inject here inside a new function.
	 
	// setup server

	server := http.Server{
		Addr:cfg.Addr,
		Handler: router,

	}

	slog.Info("server started on port 8082", slog.String("addr", cfg.Addr))


	// make a channel before starting the server and the value is store inside this channel that will be a signal and its signal of Operating system.

	done := make(chan os.Signal, 1) // this is a buffere channel thus the size will be 1.

	// now we only need to find a way to how to send a signal to the above done channel so the it will gracefully shutdown the server.
	// channel -> pipeline and as we know  it is going to use to communicate between different goroutines. so here we go:

    signal.Notify(done, os.Interrupt, syscall.SIGINT,syscall.SIGTERM) // this will notify the done channel when we get an interrupt signal from the operating system. 

	go func() {
		err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to start server:", err.Error())
	}
	}()
	
	<- done // this will block the main thread until we get any signal from the operating system and now after this now from here we can gracefully shutdown the server.

	// now we can shutdown the server gracefully.
	// before that we are going to log first and sLog which is the structured Log

	slog.Info("server is shutting down...")

	// now look some time server is going to hang so for this we need no crate a timeout context. and give it a specific time limit. if it is not going to shout down then it will going to give us a report.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)

	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server stopped sucessfully")
		
}