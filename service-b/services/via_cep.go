package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/otel"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

func GetCityByCEP(ctx context.Context, cep string) (string, error) {
	tr := otel.Tracer("service-b")
	ctx, span := tr.Start(ctx, "ViaCEP Lookup")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if containsErroField(bodyBytes) {
		return "", errors.New("can not find zipcode")
	}

	var data ViaCEPResponse
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return "", fmt.Errorf("failed to decode ViaCEP response: %w", err)
	}

	if data.Localidade == "" {
		return "", errors.New("invalid zipcode")
	}

	return data.Localidade, nil
}

func containsErroField(body []byte) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return false
	}
	_, exists := obj["erro"]
	return exists
}
