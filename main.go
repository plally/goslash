package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var PublicKey ed25519.PublicKey = fromHex("")

func InteractionEndpoint(w http.ResponseWriter, r *http.Request) {
	if !hasValidSignature(r) {
		w.WriteHeader(401)
		w.Write([]byte("no"))
		return
	}
	data := readAndReplaceBody(r)

	var interaction Interaction
	json.Unmarshal(data, &interaction)

	if interaction.Type == PING {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"type": 1}`))

		return
	}
}

func fromHex(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func readAndReplaceBody(req *http.Request) []byte {
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	_ = req.Body.Close() //  must close
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func hasValidSignature(req *http.Request) bool {
	body := readAndReplaceBody(req)

	signature := req.Header.Get("X-Signature-Ed25519")
	timestamp := req.Header.Get("X-Signature-Timestamp")

	message := append([]byte(timestamp), body...)
	return ed25519.Verify(PublicKey, message, fromHex(signature))
}

func main() {
	http.ListenAndServe(":6969", nil)
}
