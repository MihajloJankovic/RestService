package handlers

import (
	"context"
	"errors"
	"log"
	"mime"
	"net/http"

	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
)

type Porfilehendler struct {
	l  *log.Logger
	cc protos.ProfileClient
}

func NewPorfilehendler(l *log.Logger, cc protos.ProfileClient) *Porfilehendler {
	return &Porfilehendler{l, cc}

}

func (h *Porfilehendler) SetProfile(w http.ResponseWriter, sg string) bool {

	rt, err := DecodeBodyPorfileadd(sg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return false
	}
	_, err = h.cc.SetProfile(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return false
	}
	return true
}
func (h *Porfilehendler) GetProfile(w http.ResponseWriter, r *http.Request) {

	emaila := mux.Vars(r)["email"]
	ee := new(protos.ProfileRequest)
	ee.Email = emaila
	res := ValidateJwt(r, h)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	if re.GetEmail() != ee.GetEmail() {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	response, err := h.cc.GetProfile(context.Background(), ee)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusNotAcceptable)
		_, err := w.Write([]byte("Profile not found"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}
func (h *Porfilehendler) GetProfileInner(emaili string) (*protos.ProfileResponse, error) {

	emaila := emaili
	ee := new(protos.ProfileRequest)
	ee.Email = emaila

	response, err := h.cc.GetProfile(context.Background(), ee)
	if err != nil || response == nil {
		log.Printf("RPC failed: %v\n", err)
		err := errors.New("rpc failed to get inner profile")
		return nil, err
	}
	return response, nil
}
func (h *Porfilehendler) UpdateProfile(w http.ResponseWriter, r *http.Request) {

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
	rt, err := DecodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	res := ValidateJwt(r, h)
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
	_, err = h.cc.UpdateProfile(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("Couldn't update profile"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Successfully update profile"))
	if err != nil {
		return
	}
}
