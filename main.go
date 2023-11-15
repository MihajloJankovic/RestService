package main

import (
	"context"
	"errors"

	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	"github.com/MihajloJankovic/RestService/handlers"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/glavno"
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
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	connAcc, err := grpc.Dial("accommodation-service:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	cc := protos.NewProfileClient(conn)
	acc := protosAcc.NewAccommodationClient(connAcc)
	hh := handlers.NewPorfilehendler(l, cc)
	acch := handlers.NewAccommodationHandler(l, acc, hh)

	connAuth, err := grpc.Dial("auth-service:9094", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)

	ccAuth := protosAuth.NewAuthClient(connAuth)
	hhAuth := handlers.NewAuthHandler(l, ccAuth, hh)

	router := mux.NewRouter()
	router.StrictSlash(true)

	router.HandleFunc("/profile/{email}", hh.GetProfile).Methods("GET")
	router.HandleFunc("/update-profile", hh.UpdateProfile).Methods("POST")
	router.HandleFunc("/accommodation/{email}", acch.GetAccommodation).Methods("GET")
	router.HandleFunc("/accommodations", acch.GetAllAccommodation).Methods("GET")
	router.HandleFunc("/add-accommodation", acch.SetAccommodation).Methods("POST")
	router.HandleFunc("/update-accommodation", acch.UpdateAccommodation).Methods("POST")
	router.HandleFunc("/register", hhAuth.Register).Methods("POST")
	router.HandleFunc("/login", hhAuth.Login).Methods("POST")
	router.HandleFunc("/getTicket/{email}", hhAuth.GetTicket).Methods("GET")
	router.HandleFunc("/activate/{email}/{ticket}", hhAuth.Activate).Methods("GET")

	srv := &http.Server{Addr: ":9090", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
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
