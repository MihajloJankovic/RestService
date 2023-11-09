package handlers

import (
	"context"
	"errors"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/glavno"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type AccommodationHandler struct {
	l   *log.Logger
	acc protosAcc.AccommodationClient
	hh  *Porfilehendler
}

func NewAccommodationHandler(l *log.Logger, acc protosAcc.AccommodationClient, hb *Porfilehendler) *AccommodationHandler {
	return &AccommodationHandler{l, acc, hb}

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
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := *res
	if re.GetEmail() != rt.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
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
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := *res
	if re.GetEmail() != ee.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.acc.GetAccommodation(context.Background(), ee)
	if err != nil || response == nil {
		log.Println("RPC failed: %v", err)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Accommodation not found"))
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response.Dummy)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := *res
	if re.GetEmail() != rt.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
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
