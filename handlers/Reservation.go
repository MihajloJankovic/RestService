package handlers

import (
	"context"
	"errors"
	protosRes "github.com/MihajloJankovic/reservation-service/protos/genfiles"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"hash/fnv"
	"log"
	"mime"
	"net/http"
	"strconv"
)

type ReservationHandler struct {
	l   *log.Logger
	acc protosRes.ReservationClient
	hh  *Porfilehendler
	ava *AvabilityHendler
}

func NewReservationHandler(l *log.Logger, acc protosRes.ReservationClient, hb *Porfilehendler, ava *AvabilityHendler) *ReservationHandler {
	return &ReservationHandler{l, acc, hb, ava}

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
	err = h.ava.CheckAvaibility(rt.Accid, rt.DateFrom, rt.DateTo)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	hash := fnv.New32a()
	hash.Write([]byte((uuid.New()).String()))
	rt.Id = int32(hash.Sum32())
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
func (h *ReservationHandler) CheckActiveReservationByEmail(email string) error {
	temp := new(protosRes.Emaill)
	temp.Email = email
	_, err := h.acc.CheckActiveReservationByEmail(context.Background(), temp)
	if err != nil {
		log.Println("Couldn't delete reservation because it's active")
		return err
	}
	return nil
}
func (h *ReservationHandler) DeleteReservationByEmail(email string) error {
	temp := new(protosRes.Emaill)
	temp.Email = email
	_, err := h.acc.DeleteReservationByEmail(context.Background(), temp)
	if err != nil {
		log.Println("Couldn't delete reservation because it's active")
		return err
	}
	return nil
}

func (h *ReservationHandler) CheckActiveReservation(accid string) error {
	temp := new(protosRes.DateFromDateTo)
	temp.DateFrom = getTodaysDateInLocal()
	temp.DateTo = getTodaysDateInLocal()
	temp.Accid = accid
	_, err := h.acc.CheckActiveReservation(context.Background(), temp)
	if err != nil {
		log.Println("Couldn't delete reservation because it's active")
		return err
	}
	return nil
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
