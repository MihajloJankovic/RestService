package handlers

import (
	"bytes"
	"encoding/json"
	protosAcc "github.com/MihajloJankovic/accommodation-service/protos/main"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"io"
	"net/http"
)

func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
func DecodeBody(r io.Reader) (*protos.ProfileResponse, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt protos.ProfileResponse
	if err := json.Unmarshal(StreamToByte(r), &rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func DecodeBodyAcc(r io.Reader) (*protosAcc.AccommodationResponse, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt protosAcc.AccommodationResponse
	if err := json.Unmarshal(StreamToByte(r), &rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func RenderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
