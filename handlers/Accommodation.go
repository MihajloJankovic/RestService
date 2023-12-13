package handlers

import (
	"context"
	"errors"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
	"strconv"
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

func (h *AccommodationHandler) DeleteAccommodation(id string) error {
	req := &protosAcc.DeleteRequest{Uid: id}

	_, err := h.acc.DeleteAccommodation(context.Background(), req)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return err
	}
	return nil
}

func (h *AccommodationHandler) FilterByPriceRange(w http.ResponseWriter, r *http.Request) {
	//TODO GETALL ACCOMONDATIONS , THEN  check check avability service foreach
	//TODO accomondation(needs to add that method becuse current method checks only for date avaibility :) ,method shoud return all avaible objets for that
	//TODO price , so them in foreach add new foreach for avaibility list , and call reservation service checkactivereservation for dates and id of each avaibility
	//TODO (not user method in reservation ),if retturns error them pop it out of copy of avability list ,na kraju velikog foreacha ne unutrasnjeg ,
	//TODO proveris da li je ta kopija ava veca od nule ako je veca taj accomondation ostaje u kopiji liste accomondationa ako je nula znaci nema slobodnih
	//TODO  termina za taj smestaj i cenu i izbaci ga iz kopije liste . i posle foreacha kad prodje kroz sve smestaje vrati smestaje koji su ostali.
	//accs := h.acc.GetAllAccommodation()
	//for index, accs := range accs {
	//	lista_ava := h.ava.GetallbyIDandPrice() //NEEDS TO BE IMPLEMENTED
	//	//ako je lista prazna izvaci iz kopije liste accomondationa  i preskoci ostatak iteracije fora continue verovatno
	//
	//	for _, avab := range lista_ava {
	//		h.reservation.checkifThereisReservationfordate() //NEEDS TO BE IMPLEMENTED IN SERVICE CURRENCT METHODS DONT DO THE JOB
	//		if err != nil {
	//		//	izbaci iz kopije liste lista_ava
	//		}
	//	}
	//		// ako je lista_ava prazna izbaci iz kopije liste accomondationa trnutni smestaj vrv po indexu
	//}
	////vrati smestaje koji su ostali
	minPriceStr := mux.Vars(r)["min_price"]
	maxPriceStr := mux.Vars(r)["max_price"]

	// Konvertuj stringove u float64
	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		http.Error(w, "Invalid min_price parameter", http.StatusBadRequest)
		return
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		http.Error(w, "Invalid max_price parameter", http.StatusBadRequest)
		return
	}

	// Pozovi odgovarajuću metodu na AccommodationClient i obradi rezultat
	filteredAccommodations, err := h.acc.FilterByPriceRange(context.Background(), &protosAcc.PriceRangeRequest{
		MinPrice: float32(minPrice),
		MaxPrice: float32(maxPrice),
	})

	if err != nil {
		http.Error(w, "Error filtering accommodations", http.StatusInternalServerError)
		return
	}

	// Konvertuj rezultate u JSON i pošalji klijentu
	RenderJSON(w, filteredAccommodations.Dummy)
}
