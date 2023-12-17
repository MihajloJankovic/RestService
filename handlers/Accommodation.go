package handlers

import (
	"context"
	"errors"
	protosava "github.com/MihajloJankovic/Aviability-Service/protos/main"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type AccommodationHandler struct {
	l   *log.Logger
	acc protosAcc.AccommodationClient
	hh  *Porfilehendler
	ava *AvabilityHendler
}

func NewAccommodationHandler(l *log.Logger, acc protosAcc.AccommodationClient, hb *Porfilehendler, ava *AvabilityHendler) *AccommodationHandler {
	return &AccommodationHandler{l, acc, hb, ava}

}

func (h *AccommodationHandler) SetAccommodation(w http.ResponseWriter, r *http.Request) {

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
	rt.Uid = (uuid.New()).String()
	_, err = h.acc.SetAccommodation(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
	}
	w.WriteHeader(http.StatusCreated)
}
func (h *AccommodationHandler) GetOneAccommodation(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ee := new(protosAcc.AccommodationRequestOne)
	ee.Id = id

	response, err := h.acc.GetOneAccommodation(context.Background(), ee)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Accommodation not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
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
	re := res
	if re.GetRole() != "Host" {
		err := errors.New("you are not host")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if re.GetEmail() != ee.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.acc.GetAccommodation(context.Background(), ee)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Accommodation not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response.Dummy)
}
func (h *AccommodationHandler) GetAccommodationByEmail(email string) (*protosAcc.DummyList, error) {
	ee := new(protosAcc.AccommodationRequest)
	ee.Email = email

	response, err := h.acc.GetAccommodation(context.Background(), ee)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		return nil, err
	}
	return response, nil
}
func (h *AccommodationHandler) GetAllAccommodation(w http.ResponseWriter, r *http.Request) {

	emptyRequest := new(protosAcc.Emptya)
	response, err := h.acc.GetAllAccommodation(context.Background(), emptyRequest)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Accommodation not found"))
		if err != nil {
			return
		}
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
		err := errors.New("expect application/json Content-Type")
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
	_, err = h.acc.UpdateAccommodation(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Couldn't update accommodation"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully update accommodation"))
	if err != nil {
		return
	}
}

func (h *AccommodationHandler) FilterByPriceRange(w http.ResponseWriter, r *http.Request) {
	rt, err := DecodeBodyPriceAndId(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accs, err := h.acc.GetAllAccommodation(context.Background(), &protosAcc.Emptya{})
	if err != nil {
		http.Error(w, "Error getting accommodations", http.StatusInternalServerError)
		return
	}

	var availableAccommodations []*protosAcc.AccommodationResponse

	for _, acc := range accs.Dummy {
		// Dohvatanje liste dostupnosti i cena za trenutni smeštaj
		listaAva, err := h.ava.GetallbyIDandPrice(acc.Uid, rt.MinPrice, rt.MaxPrice)
		if err != nil {
			http.Error(w, "Error getting availability list...", http.StatusInternalServerError)
			return
		}

		copiedAva := make([]*protosava.CheckSet, len(listaAva))
		copy(copiedAva, listaAva)

		// Iteracija kroz kopiranu listu dostupnosti
		for _, avab := range copiedAva {
			// Provera dostupnosti za svaku stavku dostupnosti
			err := h.ava.CheckAvaibility(avab.Uid, avab.From, avab.To)
			if err != nil {
				// Ako postoji rezervacija, ukloni iz kopirane liste dostupnosti
				copiedAva = removeAvability(copiedAva, avab)
			}
		}

		// Ako je kopirana lista dostupnosti neprazna, smeštaj je dostupan
		if len(copiedAva) > 0 {
			availableAccommodations = append(availableAccommodations, acc)
		}
	}

	// Konvertuj rezultate u JSON i pošalji klijentu
	RenderJSON(w, availableAccommodations)
}

func removeAvability(slice []*protosava.CheckSet, item *protosava.CheckSet) []*protosava.CheckSet {
	var index int
	for i, v := range slice {
		if v == item {
			index = i
			break
		}
	}
	return append(slice[:index], slice[index+1:]...)
}
