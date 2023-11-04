package main

import (
	"context"
	"github.com/MihajloJankovic/RestService/handlers"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	l := log.New(os.Stdout, "standard-api", log.LstdFlags)
	conn, err := grpc.Dial("profile-service:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cc := protos.NewProfileClient(conn)
	hh := handlers.NewPorfilehendler(l, cc)
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/{email}", hh.GetProfile).Methods("GET")
	router.HandleFunc("/addprofile", hh.SetProfile).Methods("POST")
	// test ee := protos.ProfileRequest{Email: "pera@gmail.com"}
	// test cc.GetProfile(context.Background(),&ee)

	srv := &http.Server{Addr: ":9090", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
	<-quit
	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
