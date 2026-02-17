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

const shalatBaseURL = "https://equran.id/api/v2/shalat"

type ShalatService struct {
	client *http.Client
}

func NewShalatService() *ShalatService {
	return &ShalatService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ShalatService) GetProvinsi() ([]string, error) {
	resp, err := s.client.Get(shalatBaseURL + "/provinsi")
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

func (s *ShalatService) GetKabkota(provinsi string) ([]string, error) {
	payload := map[string]string{"provinsi": provinsi}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := s.client.Post(shalatBaseURL+"/kabkota", "application/json", bytes.NewBuffer(jsonPayload))
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

func (s *ShalatService) GetShalat(provinsi, kabkota string, bulan, tahun int) (*models.ShalatData, error) {
	payload := map[string]interface{}{
		"provinsi": provinsi,
		"kabkota":  kabkota,
	}
	if bulan > 0 {
		payload["bulan"] = bulan
	}
	if tahun > 0 {
		payload["tahun"] = tahun
	}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := s.client.Post(shalatBaseURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result models.ShalatResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != 200 {
		return nil, fmt.Errorf("API error: %s", result.Message)
	}

	return &result.Data, nil
}

func (s *ShalatService) GetTodaySchedule(provinsi, kabkota string) (*models.ShalatSchedule, error) {
	now := time.Now()
	data, err := s.GetShalat(provinsi, kabkota, int(now.Month()), now.Year())
	if err != nil {
		return nil, err
	}

	dayOfMonth := now.Day()
	for _, schedule := range data.Jadwal {
		if schedule.Tanggal == dayOfMonth {
			return &schedule, nil
		}
	}

	return nil, fmt.Errorf("schedule not found for day %d", dayOfMonth)
}
