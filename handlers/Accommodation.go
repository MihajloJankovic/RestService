package handlers

import (
	"context"
	"errors"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type AccommodationHandler struct {
	l   *log.Logger
	acc protosAcc.AccommodationClient
}

func NewAccommodationHandler(l *log.Logger, acc protosAcc.AccommodationClient) *AccommodationHandler {
	return &AccommodationHandler{l, acc}

}

func (h *AccommodationHandler) SetAccommodation(w http.ResponseWriter, r *http.Request) {

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
	rt, err := DecodeBodyAcc(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	_, err = h.acc.SetAccommodation(context.Background(), rt)
	if err != nil {
		log.Println("RPC failed: %v", err)
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *AccommodationHandler) GetAccommodation(w http.ResponseWriter, r *http.Request) {
	emaila := mux.Vars(r)["email"]
	ee := new(protosAcc.AccommodationRequest)
	ee.Email = emaila
	response, err := h.acc.GetAccommodation(context.Background(), ee)
	if err != nil || response == nil {
		log.Println("RPC failed: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Accommodation not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
func (h *AccommodationHandler) UpdateAccommodation(w http.ResponseWriter, r *http.Request) {

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
	rt, err := DecodeBodyAcc(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	_, err = h.acc.UpdateAccommodation(context.Background(), rt)
	if err != nil {
		log.Println("RPC failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Couldn't update accommodation"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully update accommodation"))
}
