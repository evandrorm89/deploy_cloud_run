package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ViaCepResponse struct {
	Localidade string `json:localidade`
}

type WeatherReport struct {
	Current WeatherResponse `json:current`
}

type Current struct {
	Temp_c float64 `json:temp_c`
	Temp_f float64 `json:temp_f`
}

type WeatherResponse struct {
	Temp_c float64 `json:temp_C`
	Temp_f float64 `json:temp_F`
	Temp_k float64 `json:temp_K`
}

func isValidCep(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	_, err := strconv.Atoi(cep)
	return err == nil
}

func getTempCep(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	if !isValidCep(cep) {
		http.Error(w, `{"message: "invalid zipcode"}`, http.StatusUnprocessableEntity)
		return
	}
	res, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		http.Error(w, `{message: "can not find zipcode"}`, http.StatusNotFound)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	var c ViaCepResponse
	err = json.Unmarshal(body, &c)
	if err != nil {
		http.Error(w, `{message: "Erro interno"}`, http.StatusInternalServerError)
		return
	}

	if c.Localidade == "" {
		http.Error(w, `{message: "can not find zipcode"}`, http.StatusNotFound)
		return
	}

	res, err = http.Get("https://api.weatherapi.com/v1/current.json?key=602ac96551be4db2b0112256243006&q=" + c.Localidade + "&aqi=no")
	if err != nil {
		http.Error(w, `{message: "Erro ao achar o tempo atual para a localidade informada"}`, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)

	var t WeatherReport
	err = json.Unmarshal(body, &t)

	tempC := t.Current.Temp_c
	tempF := t.Current.Temp_f
	tempK := t.Current.Temp_c + 273.0

	response := WeatherResponse{
		Temp_c: tempC,
		Temp_f: tempF,
		Temp_k: tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/weather/{cep}", getTempCep)
	http.ListenAndServe(":8080", r)
}
