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
type RequestRegister struct {
	Email     string
	Firstname string
	Lastname  string
	Birthday  string
	Gender    string
	Role      string
	Password  string
	Username  string
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
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	rt, err := DecodeBodyAuth(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	out := new(protos.ProfileSet)
	out.Email = rt.Email
	out.Password = rt.Password
	out.Firstname = rt.Firstname
	out.Lastname = rt.Lastname
	out.Birthday = rt.Birthday
	out.Gender = rt.Gender
	out.Username = rt.Username
	if rt.Role != "Guest" && rt.Role != "Host" {
		rt.Role = "Guest"
	}
	out.Role = rt.Role
	_, err = h.cc.SetProfile(context.Background(), out)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, "couldn't create user,something went wrong'")
	} else {
		w.WriteHeader(http.StatusCreated)
		RenderJSON(w, "registered")
	}
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
		_, err := w.Write([]byte("couldn't update profile"))
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("successfully update profile"))
	if err != nil {
		return
	}
}
