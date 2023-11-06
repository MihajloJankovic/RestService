package main

import (
	"context"
	"github.com/MihajloJankovic/RestService/handlers"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	//protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
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
	connAcc, err := grpc.Dial("accommodation-service:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cc := protos.NewProfileClient(conn)
	acc := protosAcc.NewAccommodationClient(connAcc)
	hh := handlers.NewPorfilehendler(l, cc)
	acch := handlers.NewAccommodationHandler(l, acc)
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/profile/{email}", hh.GetProfile).Methods("GET")
	router.HandleFunc("/add-profile", hh.SetProfile).Methods("POST")
	router.HandleFunc("/update-profile", hh.UpdateProfile).Methods("POST")
	router.HandleFunc("/accommodation/{email}", acch.GetAccommodation).Methods("GET")
	router.HandleFunc("/add-accommodation", acch.SetAccommodation).Methods("POST")
	router.HandleFunc("/update-accommodation", acch.UpdateAccommodation).Methods("POST")
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
