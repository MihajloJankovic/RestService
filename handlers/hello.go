package handlers

import (
	"context"
	"errors"
	"fmt"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
)

type Hello struct {
	l  *log.Logger
	cc protos.ProfileClient
}

func NewHello(l *log.Logger, cc protos.ProfileClient) *Hello {
	return &Hello{l, cc}

}

func (h *Hello) SetProfile(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := DecodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	h.l.Println(rt)

	response, err := h.cc.SetProfile(context.Background(), rt)
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Printf("Response: %v", response)

	//fmt.Fprintf(w, "Hello %s", response.GetFirstname())
	//fmt.Fprintf(w, "Hello %s", response.String())
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s", d)
}
func (h *Hello) GetProfile(w http.ResponseWriter, r *http.Request) {
	emaila := mux.Vars(r)["email"]
	ee := new(protos.ProfileRequest)
	ee.Email = emaila
	response, err := h.cc.GetProfile(context.Background(), ee)
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Printf("Response: %v", response)

	//fmt.Fprintf(w, "Hello %s", response.GetFirstname())
	//fmt.Fprintf(w, "Hello %s", response.String())
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s", d)
}
