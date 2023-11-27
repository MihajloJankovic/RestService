package handlers

import (
	"context"
	"errors"
	protosava "github.com/MihajloJankovic/Aviability-Service/protos/main"
	"log"
	"mime"
	"net/http"

	protos "github.com/MihajloJankovic/Aviability-Service/protos/main"
	protosac "github.com/MihajloJankovic/accommodation-service/protos/main"
)

type AvabilityHendler struct {
	l  *log.Logger
	cc protos.AccommodationAviabilityClient
	ac protosac.AccommodationClient
	pp *Porfilehendler
}

func NewAvabilityHendler(l *log.Logger, cc protos.AccommodationAviabilityClient, cca protosac.AccommodationClient, pp *Porfilehendler) *AvabilityHendler {
	return &AvabilityHendler{l, cc, cca, pp}

}

func (h *AvabilityHendler) SetAvability(w http.ResponseWriter, r *http.Request) {
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
	rt, err := DecodeBodyAva(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	res := ValidateJwt(r, h.pp)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, err = h.cc.SetAccommodationAviability(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Couldn't create avability"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully created"))
	if err != nil {
		return
	}

}
func (h *AvabilityHendler) CheckAvaibility(accid string, dateFrom string, dateTo string) error {
	temp := new(protos.CheckRequest)
	temp.Id = accid
	temp.From = dateFrom
	temp.To = dateTo

	_, err := h.cc.GetAccommodationCheck(context.Background(), temp)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return err
	}
	return nil
}
func (h *AvabilityHendler) GetAllbyId(w http.ResponseWriter, r *http.Request) {

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
	rt, err := DecodeBodyAva3(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	res := ValidateJwt(r, h.pp)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	s, err := h.cc.GetAllforAccomendation(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Not available"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, s)

}
func (h *AvabilityHendler) DeleteByAccomndation(accid string) error {
	req := new(protosava.DeleteRequestb)
	req.Uid = accid
	//this is deleted by accomndation uid
	_, err := h.cc.DeleteByUser(context.Background(), req)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return err
	}
	return nil
}
