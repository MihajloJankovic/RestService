package handlers

import (
	"context"
	"errors"
	protosRes "github.com/MihajloJankovic/reservation-service/protos/genfiles"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"strconv"
)

type ReservationHandler struct {
	l   *log.Logger
	acc protosRes.ReservationClient
	hh  *Porfilehendler
}

func NewReservationHandler(l *log.Logger, acc protosRes.ReservationClient, hb *Porfilehendler) *ReservationHandler {
	return &ReservationHandler{l, acc, hb}

}

func (h *ReservationHandler) SetReservation(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBodyRes(r.Body)
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
	re := res
	if re.GetEmail() != rt.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	//TODO Add call to availability service for check if date is available
	_, err = h.acc.SetReservation(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *ReservationHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	ida := mux.Vars(r)["id"]
	ee := new(protosRes.ReservationRequest)
	vv, _ := strconv.ParseInt(ida, 10, 32)
	ee.Id = int32(vv)
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	response, err := h.acc.GetReservation(context.Background(), ee)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Accommodation not found"))
		if err != nil {
			return
		}
		return
	}
	if re.GetEmail() != response.GetDummy()[0].GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response.Dummy)
}
func (h *ReservationHandler) GetAllReservation(w http.ResponseWriter, r *http.Request) {
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	if re.GetRole() != "Host" {
		err := errors.New("you are not host")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	emptyRequest := new(protosRes.Emptyaa)
	response, err := h.acc.GetAllReservations(context.Background(), emptyRequest)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Reservation not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response.Dummy)
}
func (h *ReservationHandler) DeleteByAccomndation(accid string) error {
	req := new(protosRes.DeleteRequestaa)
	req.Uid = accid
	_, err := h.acc.DeleteByAccomnendation(context.Background(), req)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return err
	}
	return nil
}
func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {

	contentType := r.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBodyRes(r.Body)
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
	re := res
	if re.GetRole() != "Host" {
		err := errors.New("you are not host")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if re.GetEmail() != rt.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, err = h.acc.UpdateReservation(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Couldn't update reservation"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully update reservation"))
	if err != nil {
		return
	}
}
