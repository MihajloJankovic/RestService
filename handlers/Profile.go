package handlers

import (
	"context"
	"errors"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type Porfilehendler struct {
	l  *log.Logger
	cc protos.ProfileClient
}

func NewPorfilehendler(l *log.Logger, cc protos.ProfileClient) *Porfilehendler {
	return &Porfilehendler{l, cc}

}

func (h *Porfilehendler) SetProfile(w http.ResponseWriter, r *http.Request) {

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
	_, err = h.cc.SetProfile(context.Background(), rt)
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *Porfilehendler) GetProfile(w http.ResponseWriter, r *http.Request) {
	emaila := mux.Vars(r)["email"]
	ee := new(protos.ProfileRequest)
	ee.Email = emaila
	response, err := h.cc.GetProfile(context.Background(), ee)
	if err != nil || response == nil {
		log.Println("RPC failed: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Profile not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
