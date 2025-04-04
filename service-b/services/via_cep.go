package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro,omitempty"`
}

func GetCityByCEP(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	tr := otel.Tracer("service-b")
	ctx, span := tr.Start(ctx, "ViaCEP Lookup")
	defer span.End()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to call ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep returned status: %d", resp.StatusCode)
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode ViaCEP response: %w", err)
	}

	if data.Erro {
		return "", errors.New("can not find zipcode")
	}

	if data.Localidade == "" {
		return "", errors.New("invalid zipcode")
	}

	return data.Localidade, nil
}
