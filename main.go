package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var PublicKey ed25519.PublicKey = fromHex("5edbddae50f6351a2a7bd049c2daf071ea4c67503ff8a980125c2172562957cb")

func InteractionEndpoint(w http.ResponseWriter, r *http.Request) {
	if !hasValidSignature(r) {
		w.WriteHeader(401)
		w.Write([]byte("does this satisfy you *discord*"))
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

	foxUrl := getFoxUrl()

	w.WriteHeader(200)
	respString := fmt.Sprintf(`{
		"type": 4,
		"data": {
			"tts": false,
			"content": "%v",
			"embeds": [],
			"allowed_mentions": []
		}
	}`, foxUrl)
	w.Write([]byte(respString))
}

func getFoxUrl() string {
	resp, err := http.Get("https://api.foxorsomething.net/fox/random.json")
	if err != nil {
		return "not found"
	}
	defer resp.Body.Close()

	var foxInfo map[string]string
	foxData, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(foxData, &foxInfo)
	id := foxInfo["id"]
	return fmt.Sprintf("https://api.foxorsomething.net/fox/%v.png", id)
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
	http.HandleFunc("/log/interaction", InteractionEndpoint)

	http.ListenAndServe(":6969", nil)
}
