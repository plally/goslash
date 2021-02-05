package lambda

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/plally/goslash/goslash"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Listener struct {
	PublicKey ed25519.PublicKey
	Handler   goslash.InteractionHandler
}

func (listener *Listener) SetHandler(handler goslash.InteractionHandler) {
	listener.Handler = handler
}

func NewListener(publicKey string) *Listener {
	return &Listener{
		PublicKey: fromHex(publicKey),
	}
}

// Handler for GET {date}
func (listener *Listener) lambdaHandler(req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat:   "",
		DisableTimestamp:  false,
		DisableHTMLEscape: false,
		DataKey:           "",
		FieldMap:          nil,
		CallerPrettyfier:  nil,
		PrettyPrint:       false,
	})
	logger := log.WithField("listener_type", "lambda")

	signature := req.Headers["x-signature-ed25519"]
	timestamp := req.Headers["x-signature-timestamp"]

	body := []byte(req.Body)
	message := append([]byte(timestamp), body...)

	if !ed25519.Verify(listener.PublicKey, message, fromHex(signature)) {
		logger.Info("Invalid request signature returning 401")
		return statusResponse(http.StatusUnauthorized), nil
	}

	var interaction goslash.Interaction
	err := json.Unmarshal(body, &interaction)
	if err != nil {
		logger.WithField("error", err).Warn("error unmarshalling interaction")
		return statusResponse(http.StatusInternalServerError), err
	}

	logger = logger.WithField("interaction", interaction)

	if interaction.Type == goslash.PING {
		logger.Info("received ping interaction, responding with pong")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"type": 1}`,
		}, nil
	}

	logger.Info("sending interaction to handler")

	response := listener.Handler(&interaction)
	if response == nil {
		logger.Info("handler did not return a response, setting response to ACK")
		response = &goslash.InteractionResponse{
			Type: goslash.ACK,
		}
	}

	logger = logger.WithField("response", response)
	data, err := json.Marshal(response)

	if err != nil {
		log.WithField("error", err).Errorf("Could not marshal response")
		return statusResponse(http.StatusInternalServerError), err
	}

	logger.WithField("response", response).Info("Returning response from interaction")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		MultiValueHeaders: nil,
		Body:              string(data),
		IsBase64Encoded:   false,
	}, nil

}

func (listener *Listener) Start() {
	lambda.Start(listener.lambdaHandler)
}

// helpers
func statusResponse(status int) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}
}

func fromHex(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
