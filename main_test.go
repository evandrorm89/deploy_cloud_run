package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockClient struct{}

func (m *mockClient) Get(url string) (*http.Response, error) {
	var body io.ReadCloser

	if url == "https://viacep.com.br/ws/01001000/json/" {
		body = io.NopCloser(strings.NewReader(`{"localidade": "São Paulo"}`))
	} else if url == "https://viacep.com.br/ws/00000000/json/" {
		body = io.NopCloser(strings.NewReader(`{}`))
	} else {
		body = io.NopCloser(strings.NewReader(`{}`))
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil
}

var httpGet = http.Get

func TestIsValidCep(t *testing.T) {
	validCep := "11111111"
	invalidCep := "asdf123"
	shortCep := "123"

	assert.True(t, isValidCep(validCep))
	assert.False(t, isValidCep(invalidCep))
	assert.False(t, isValidCep(shortCep))
}

func TestGetTempCepBatch(t *testing.T) {
	httpGet = (&mockClient{}).Get
	r := chi.NewRouter()
	r.Get("/weather/{cep}", getTempCep)

	tests := []struct {
		name       string
		cep        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Valid CEP",
			cep:        "01001000",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid CEP format",
			cep:        "1234abcd",
			wantStatus: http.StatusUnprocessableEntity,
			wantBody:   `{"message": "invalid zipcode"}`,
		},
		{
			name:       "Non-existent CEP",
			cep:        "00000000",
			wantStatus: http.StatusNotFound,
			wantBody:   `{"message": "can not find zip code"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/weather/"+tt.cep, nil)
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, rr.Body.String())
			} else {
				var response WeatherResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotZero(t, response.Temp_c)
				assert.NotZero(t, response.Temp_f)
				assert.NotZero(t, response.Temp_k)
			}
		})
	}
}
