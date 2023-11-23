package handlers

import (
	"context"
	"errors"
	"log"
	"mime"
	"net/http"

	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	l  *log.Logger
	cc protosAuth.AuthClient
	hh *Porfilehendler
}
type RequestRegister struct {
	Email     string
	Firstname string
	Lastname  string
	Birthday  string
	Gender    bool
	Role      string
	Password  string
}

func NewAuthHandler(l *log.Logger, cc protosAuth.AuthClient, hb *Porfilehendler) *AuthHandler {
	return &AuthHandler{l, cc, hb}

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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
	out := new(protosAuth.AuthRequest)
	out.Email = rt.Email
	out.Password = rt.Password
	out2 := new(protos.ProfileResponse)
	out2.Email = rt.Email
	out2.Firstname = rt.Firstname
	out2.Lastname = rt.Lastname
	out2.Birthday = rt.Birthday
	out2.Gender = rt.Gender
	if rt.Role != "Guest" && rt.Role != "Host" {
		rt.Role = "Guest"
	}
	out2.Role = rt.Role
	payload, err := ToJSON(out2)

	val := h.hh.SetProfile(w, payload)
	if val == false {
		w.WriteHeader(http.StatusBadRequest)
		RenderJSON(w, "couldn't create user,something went wrong'")
	} else {
		_, err = h.cc.Register(context.Background(), out)
		if err != nil {
			log.Printf("RPC failed: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
		RenderJSON(w, "registered")
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	rt, err := DecodeBodyAuthLog(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	response, err := h.cc.Login(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Login failed"))
		if err != nil {
			return
		}
		return
	}
	GenerateJwt(w, rt.GetEmail())
	w.WriteHeader(http.StatusOK)
	// adds token to request header
	RenderJSON(w, response)
}

func (h *AuthHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	response, err := h.cc.GetTicket(context.Background(), &protosAuth.AuthGet{Email: email})
	if err != nil {
		log.Println("RPC failed: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get ticket"))
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}

func (h *AuthHandler) Activate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	email := params["email"]
	ticket := params["ticket"]

	_, err := h.cc.Activate(context.Background(), &protosAuth.ActivateRequest{Email: email, Ticket: ticket})
	if err != nil {
		log.Println("RPC failed:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to activate account"))
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, "Activated account")

}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
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
	rt, err := DecodeBodyPassword(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	out := new(protosAuth.ChangePasswordRequest)
	out.Email = rt.Email
	out.CurrentPassword = rt.CurrentPassword
	out.NewPassword = rt.NewPassword

	_, err = h.cc.ChangePassword(context.Background(), out)
	if err != nil {
		log.Println("RPC failed: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failure in password changing!"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	RenderJSON(w, "Password changed successfully!")

}
