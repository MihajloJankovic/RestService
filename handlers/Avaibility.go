package handlers

import (
	"context"
	"errors"
	protos "github.com/MihajloJankovic/Aviability-Service/protos/main"
	protosava "github.com/MihajloJankovic/Aviability-Service/protos/main"
	protosac "github.com/MihajloJankovic/accommodation-service/protos/main"
	"log"
	"mime"
	"net/http"
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
		_, err := w.Write([]byte("couldn't create avability"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)

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

func (h *AvabilityHendler) GetallbyIDandPrice(id string, minPrice int32, maxPrice int32) ([]*protosava.CheckSet, error) {
	// Kreiraj PriceAndIdRequest objekat sa prosleđenim vrednostima
	request := &protosava.PriceAndIdRequest{
		Id:       id,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
	}

	// Pozovi odgovarajuću metodu iz availability servisa
	listaAva, err := h.cc.GetallbyIDandPrice(context.Background(), request)
	if err != nil {
		// Obradi grešku, loguj je ili vrati odgovarajući error
		return nil, err
	}

	// Ako je sve u redu, možeš dalje raditi sa listom dostupnosti ili vratiti korisniku
	return listaAva.Dummy, nil
}
