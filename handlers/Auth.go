package handlers

import (
	"context"
	"errors"
	"log"
	"mime"
	"net/http"

	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
)

type AuthHandler struct {
	l  *log.Logger
	cc protosAuth.AuthClient
}

func NewAuthHandler(l *log.Logger, cc protosAuth.AuthClient) *AuthHandler {
	return &AuthHandler{l, cc}

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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
	rt, err := DecodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	response, err := h.cc.Register(context.Background(), rt)
	if err != nil {
		log.Println("RPC failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Registration failed"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	RenderJSON(w, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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
	rt, err := DecodeBody(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}
	response, err := h.cc.Login(context.Background(), rt)
	if err != nil {
		log.Println("RPC failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Login failed"))
		return
	}
	w.WriteHeader(http.StatusOK)
	RenderJSON(w, response)
}

func (h *Porfilehendler) GetAuth(w http.ResponseWriter, r *http.Request) {
	// Your logic to fetch authentication data
	auths, err := h.cc.GetAuth(context.Background(), &protos.AuthRequest{})
	if err != nil {
		log.Println("Failed to get authentication data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get authentication data"))
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, auths)
}
