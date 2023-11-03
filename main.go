package main

import (
	"context"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"servis1/handlers"
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
	hh := handlers.NewHello(l, cc)
	gb := handlers.NewGoodBye(l)
	sm := http.NewServeMux()

	// test ee := protos.ProfileRequest{Email: "pera@gmail.com"}
	// test cc.GetProfile(context.Background(),&ee)
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gb)
	srv := &http.Server{Addr: ":9090", Handler: sm}
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
