package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	protosava "github.com/MihajloJankovic/Aviability-Service/protos/main"
	"github.com/MihajloJankovic/RestService/handlers"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	protosRes "github.com/MihajloJankovic/reservation-service/protos/genfiles"
	habb "github.com/gorilla/handlers"
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
	defer func(connAcc *grpc.ClientConn) {
		err := connAcc.Close()
		if err != nil {

		}
	}(connAcc)
	connRes, err := grpc.Dial("reservation-service:9096", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(connRes *grpc.ClientConn) {
		err := connRes.Close()
		if err != nil {

		}
	}(connRes)
	resc := protosRes.NewReservationClient(connRes)
	cc := protos.NewProfileClient(conn)
	acc := protosAcc.NewAccommodationClient(connAcc)
	hh := handlers.NewPorfilehendler(l, cc)
	acch := handlers.NewAccommodationHandler(l, acc, hh)
	resh := handlers.NewReservationHandler(l, resc, hh)

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
	connAva, err := grpc.Dial("avaibility-service:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	la := log.New(os.Stdout, "standard-avaibility-api", log.LstdFlags)
	ccava := protosava.NewAccommodationAviabilityClient(connAva)
	hhava := handlers.NewAvabilityHendler(la, ccava, acc, hh)

	router := mux.NewRouter()
	router.StrictSlash(true)
	//profile
	router.HandleFunc("/profile/{email}", hh.GetProfile).Methods("GET")
	router.HandleFunc("/update-profile", hh.UpdateProfile).Methods("POST")
	//accommondation
	router.HandleFunc("/accommodation/{email}", acch.GetAccommodation).Methods("GET")
	router.HandleFunc("/accommodations", acch.GetAllAccommodation).Methods("GET")
	router.HandleFunc("/add-accommodation", acch.SetAccommodation).Methods("POST")
	router.HandleFunc("/update-accommodation", acch.UpdateAccommodation).Methods("POST")
	//reservation
	router.HandleFunc("/reservation/{id}", resh.GetReservation).Methods("GET")
	router.HandleFunc("/reservations", resh.GetAllReservation).Methods("GET")
	router.HandleFunc("/set-accommodation", resh.SetReservation).Methods("POST")
	router.HandleFunc("/update-accommodation", resh.UpdateReservation).Methods("POST")
	//auth
	router.HandleFunc("/register", hhAuth.Register).Methods("POST")
	router.HandleFunc("/login", hhAuth.Login).Methods("POST")
	router.HandleFunc("/getTicket/{email}", hhAuth.GetTicket).Methods("GET")
	router.HandleFunc("/activate/{email}/{ticket}", hhAuth.Activate).Methods("GET")
	router.HandleFunc("/change-password", hhAuth.ChangePassword).Methods("POST")
	router.HandleFunc("/request-reset", hhAuth.RequestPasswordReset).Methods("POST")
	router.HandleFunc("/reset", hhAuth.ResetPassword).Methods("POST")
	//avaibility
	router.HandleFunc("/set-avaibility", hhava.SetAvability).Methods("POST")
	router.HandleFunc("/get-all-avaibility", hhava.GetAllbyId).Methods("POST")
	router.HandleFunc("/check-avaibility", hhava.CheckAvaibility).Methods("POST")

	headersOk := habb.AllowedHeaders([]string{"Content-Type", "Authorization"})
	originsOk := habb.AllowedOrigins([]string{"http://localhost:4200"}) // Replace with your frontend origin
	methodsOk := habb.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Use the CORS middleware
	corsRouter := habb.CORS(originsOk, headersOk, methodsOk)(router)

	// Start the server
	srv := &http.Server{Addr: ":9090", Handler: corsRouter}
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
