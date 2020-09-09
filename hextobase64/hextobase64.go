package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Data struct {
	Hex    string `json:"hex"`
	Base64 string `json:"base64"`
}

func base64ToHex(base64Data []byte) []byte {
	hexData := make([]byte, hex.EncodedLen(len(base64Data)))
	hex.Encode(hexData, base64Data)
	return hexData

}

func hexToBase64(hexData []byte) []byte {
	base64Data := make([]byte, base64.StdEncoding.EncodedLen(len(hexData)))
	base64.StdEncoding.Encode(base64Data, hexData)
	return base64Data
}

func convertBase64ToHex(w http.ResponseWriter, r *http.Request) {
	var d Data
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	err := jsonDecoder.Decode(&d)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	base64Data, err := base64.StdEncoding.DecodeString(d.Base64)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.Hex = string(base64ToHex(base64Data))

	resp, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func convertHexToBase64(w http.ResponseWriter, r *http.Request) {
	var d Data
	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	err := jsonDecoder.Decode(&d)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hexData, err := hex.DecodeString(d.Hex)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.Base64 = string(hexToBase64(hexData))

	resp, _ := json.Marshal(d)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

}

func main() {
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/hex-to-base64", convertHexToBase64).Methods(http.MethodPost)
	api.HandleFunc("/base64-to-hex", convertBase64ToHex).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))

}
