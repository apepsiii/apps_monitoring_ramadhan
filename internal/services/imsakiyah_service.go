package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

const baseURL = "https://equran.id/api/v2/imsakiyah"

type ImsakiyahService struct {
	client *http.Client
}

func NewImsakiyahService() *ImsakiyahService {
	return &ImsakiyahService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ImsakiyahService) GetProvinsi() ([]string, error) {
	resp, err := s.client.Get(baseURL + "/provinsi")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.ProvinsiResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (s *ImsakiyahService) GetKabkota(provinsi string) ([]string, error) {
	payload := map[string]string{"provinsi": provinsi}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := s.client.Post(baseURL+"/kabkota", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.KabkotaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (s *ImsakiyahService) GetImsakiyah(provinsi, kabkota string) (*models.ImsakiyahData, error) {
	payload := map[string]string{
		"provinsi": provinsi,
		"kabkota":  kabkota,
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := s.client.Post(baseURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.ImsakiyahResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return &result.Data, nil
}

func (s *ImsakiyahService) GetTodaySchedule(provinsi, kabkota string, dayOfMonth int) (*models.ImsakiyahSchedule, error) {
	data, err := s.GetImsakiyah(provinsi, kabkota)
	if err != nil {
		return nil, err
	}

	for _, schedule := range data.Imsakiyah {
		if schedule.Tanggal == dayOfMonth {
			return &schedule, nil
		}
	}

	return nil, fmt.Errorf("schedule not found for day %d", dayOfMonth)
}
