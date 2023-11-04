package handlers

import (
	"bytes"
	"encoding/json"
	protos "github.com/MihajloJankovic/profile-service/protos/main"
	"github.com/google/uuid"
	"io"
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

//
//func decodeGroupBody(r io.Reader) (*ConfigGroup, error) {
//	dec := json.NewDecoder(r)
//	dec.DisallowUnknownFields()
//
//	var rt ConfigGroup
//	if err := json.Unmarshal(StreamToByte(r), &rt); err != nil {
//		return nil, err
//	}
//	return &rt, nil
//}
//
//func renderJSON(w http.ResponseWriter, v interface{}) {
//	js, err := json.Marshal(v)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write(js)
//}

func createId() string {
	return uuid.New().String()
}
