package handlers

import (
	"context"
	"errors"
	protosAuth "github.com/MihajloJankovic/Auth-Service/protos/main"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/gorilla/mux"
	"log"
	"mime"
	"net/http"
)

type AuthHandler struct {
	l    *log.Logger
	cc   protosAuth.AuthClient
	hh   *Porfilehendler
	resh *ReservationHandler
	acch *AccommodationHandler
	avah *AvabilityHendler
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

func NewAuthHandler(l *log.Logger, cc protosAuth.AuthClient, hb *Porfilehendler, resh *ReservationHandler, acch *AccommodationHandler, avah *AvabilityHendler) *AuthHandler {
	return &AuthHandler{l, cc, hb, resh, acch, avah}

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
	out2.Username = rt.Username
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
	_, err = h.cc.Login(context.Background(), rt)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Login failed"))
		if err != nil {
			return
		}
		return
	}
	jwt := GenerateJwt(w, rt.GetEmail())
	w.WriteHeader(http.StatusOK)
	// adds token to request header
	RenderJSON(w, jwt)
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
func (h *AuthHandler) DeleteReservation(accid string) error {

	err := h.resh.DeleteByAccomndation(accid)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		return err
	}
	return nil
}

func (h *AuthHandler) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.cc.RequestPasswordReset(context.Background(), &protosAuth.AuthGet{Email: rt.Email})
	if err != nil {
		log.Println("RPC failed: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to request password reset"))
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, "Password reset requested successfully")
}
func (h *AuthHandler) DeleteHost(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	user, err := h.hh.GetProfileInner(email)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Couldn't delete host"))

		return
	}
	if user.GetRole() != "Host" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("You are not a host"))

		return
	}
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	if re.GetEmail() != email {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	accommodations, err := h.acch.GetAccommodationByEmail(email)
	for _, acc := range accommodations.Dummy {
		err = h.resh.CheckActiveReservation(acc.GetUid())
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
	}
	for _, acc := range accommodations.Dummy {
		err := h.resh.DeleteByAccomndation(acc.GetUid())
		if err != nil {
			log.Printf("RPC failed: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Couldn't delete host"))

			return
		}
		err = h.avah.DeleteByAccomndation(acc.GetUid())
		if err != nil {
			log.Printf("RPC failed: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Availability service unavaible"))

			return
		}
	}
	err = h.hh.DeleteProfile(email)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Couldn't delete account"))
		if err != nil {
			return
		}
		return
	}
	temp := new(protosAuth.AuthGet)
	temp.Email = email
	_, err = h.cc.Delete(context.Background(), temp)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Couldn't delete account"))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, "Account deleted successfully")
}
func (h *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	user, err := h.hh.GetProfileInner(email)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Couldn't delete guest"))

		return
	}
	if user.GetRole() != "Guest" {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("You are not a guest"))

		return
	}
	res := ValidateJwt(r, h.hh)
	if res == nil {
		err := errors.New("jwt error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	re := res
	if re.GetEmail() != email {
		err := errors.New("authorization error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// TODO Check if the user had any reservations active and delete if he doesn't have any KOPI PASTA SAMO EMAIL UMESTO ACCIDA OD HOSTA, CHECKACTIVERESERVATION(EMAIL)

	err = h.hh.DeleteProfile(email)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Couldn't delete account"))
		if err != nil {
			return
		}
		return
	}
	temp := new(protosAuth.AuthGet)
	temp.Email = email
	_, err = h.cc.Delete(context.Background(), temp)
	if err != nil {
		log.Printf("RPC failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Couldn't delete account"))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, "Account deleted successfully")
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
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
	rt, err := DecodeBodyReset(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusAccepted)
		return
	}

	_, err = h.cc.ResetPassword(context.Background(), rt)
	if err != nil {
		log.Println("RPC failed: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset password"))
		return
	}

	w.WriteHeader(http.StatusOK)
	RenderJSON(w, "Password reset successfully")
}
