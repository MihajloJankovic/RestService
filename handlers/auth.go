package handlers

import (
	protos "github.com/MihajloJankovic/Auth-Service/protos/main"
	"log"
	"net/http"
)

type Authhendler struct {
	l  *log.Logger
	cc protos.AuthClient
}

func NewAuthhendler(l *log.Logger, cc protos.AuthClient) *Porfilehendler {
	return &Porfilehendler{l, cc}

}

func (h *Porfilehendler) Login(w http.ResponseWriter, r *http.Request) {

	return
}
func (h *Porfilehendler) Register(w http.ResponseWriter, r *http.Request) {
	return
}
