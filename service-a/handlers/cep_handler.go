package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"

	"go.opentelemetry.io/otel"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func HandleCEP(w http.ResponseWriter, r *http.Request) {
	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if !isValidCEP(req.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	serviceBURL := os.Getenv("SERVICE_B_URL")
	if serviceBURL == "" {
		serviceBURL = "http://localhost:8081"
	}

	jsonBody, _ := json.Marshal(req)

	tr := otel.Tracer("service-a")
	ctx, span := tr.Start(r.Context(), "Call Service B")
	defer span.End()

	reqToB, err := http.NewRequestWithContext(ctx, http.MethodPost, serviceBURL+"/weather", bytes.NewBuffer(jsonBody))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	reqToB.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(reqToB)
	if err != nil {
		http.Error(w, "failed to contact service B", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	w.Write(body)
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
