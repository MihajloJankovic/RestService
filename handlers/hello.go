package handlers

import (
	"context"
	"fmt"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l  *log.Logger
	cc protos.ProfileClient
}

func NewHello(l *log.Logger, cc protos.ProfileClient) *Hello {
	return &Hello{l, cc}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("hello world")
	ee := new(protos.ProfileRequest)
	ee.Email = "pera@gmail.com"
	response, err := h.cc.GetProfile(context.Background(), ee)
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Printf("Response: %v", response)

	fmt.Fprintf(w, "Hello %s", response.GetFirstname())
	fmt.Fprintf(w, "Hello %s", response.String())
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s", d)
}
