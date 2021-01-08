package httplistener

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"github.com/plally/fox_pics_slash_commands/goslash"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// TODO test http handler
type Listener struct {
	PublicKey ed25519.PublicKey
	Handler   goslash.InteractionHandler
}

func NewHttpListener(publicKey string) *Listener {
	return &Listener{
		PublicKey: fromHex(publicKey),
	}
}

func (listener *Listener) SetHandler(handler goslash.InteractionHandler) {
	listener.Handler = handler
}

func (listener *Listener) ListenAndServe(addr, path string) {
	mux := http.NewServeMux()
	mux.HandleFunc(path, listener.InteractionEndpoint)

	http.ListenAndServe(addr, mux)
}

func (listener *Listener) InteractionEndpoint(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")

	message := append([]byte(timestamp), body...)
	if !ed25519.Verify(listener.PublicKey, message, fromHex(signature)) {
		w.WriteHeader(401)
		w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var interaction goslash.Interaction
	err := json.Unmarshal(body, &interaction)
	if err != nil {
		log.Error(err)
	}

	if interaction.Type == goslash.PING {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"type": 1}`))
		return
	}

	response := listener.Handler(&interaction)
	if response != nil {
		data, err := json.Marshal(response)
		if err != nil {
			log.Error(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

}

func fromHex(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
