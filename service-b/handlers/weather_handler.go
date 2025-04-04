package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"service-b/services"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func HandleWeather(w http.ResponseWriter, r *http.Request) {
	var req CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	if !isValidCEP(req.CEP) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := services.GetCityByCEP(r.Context(), req.CEP)
	if err != nil {
		if err.Error() == "can not find zipcode" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		return
	}

	tempC, err := services.GetTemperatureByCity(r.Context(), city)
	if err != nil {
		http.Error(w, "failed to get temperature", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: services.CelsiusToFahrenheit(tempC),
		TempK: services.CelsiusToKelvin(tempC),
	})
}

func isValidCEP(cep string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, cep)
	return match
}
