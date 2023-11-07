package main

import (
	"context"

	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	"github.com/MihajloJankovic/RestService/handlers"
	protos "github.com/MihajloJankovic/profile-service/protos/main"

	//protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	conna, err := grpc.Dial("auth-service:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	ccAuth := protosAuth.NewAuthClient(conna)
	hhAuth := handlers.NewAuthhendler(l, ccAuth)
	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/profile/{email}", hh.GetProfile).Methods("GET")
	router.HandleFunc("/addprofile", hh.SetProfile).Methods("POST")
	router.HandleFunc("/update-profile", hh.UpdateProfile).Methods("POST")
	// test ee := protos.ProfileRequest{Email: "pera@gmail.com"}
	// test cc.GetProfile(context.Background(),&ee)
	router.HandleFunc("/register", hhAuth.Register).Methods("POST")
	router.HandleFunc("/login", hhAuth.Login).Methods("POST")
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
