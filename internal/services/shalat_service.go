package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// Nominatim Reverse Geocoding Response
type NominatimResponse struct {
	Address struct {
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		County      string `json:"county"`
		State       string `json:"state"` // Province
		Country     string `json:"country"`
	} `json:"address"`
	DisplayName string `json:"display_name"`
}

func (s *ShalatService) ReverseGeocode(lat, long float64) (string, string, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=10", lat, long)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}
	// Nominatim requires a User-Agent identify your application
	req.Header.Set("User-Agent", "AmaliahRamadhanApp/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var result NominatimResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", err
	}

	// Extract city/kabupaten and province
	city := result.Address.City
	if city == "" {
		city = result.Address.Town
	}
	if city == "" {
		city = result.Address.County // Kabupaten often maps to county
	}
	// Fallback to village if really needed, but usually not for Shalat API matching
	
	province := result.Address.State

	return province, city, nil
}

func (s *ShalatService) MatchLocation(detectedProv, detectedCity string) (string, string, error) {
	// 1. Get All Provinces
	provinces, err := s.GetProvinsi()
	if err != nil {
		return "", "", err
	}

	// 2. Find best matching province
	bestProv := ""
	bestProvScore := -1

	// Normalize detected province (remove "Provinsi", "Daerah Istimewa", etc if needed, but simple contains/levenshtein might work)
	normalizedDetectedProv := normalizeString(detectedProv)

	for _, p := range provinces {
		score := similarity(normalizedDetectedProv, normalizeString(p))
		if score > bestProvScore {
			bestProvScore = score
			bestProv = p
		}
	}

	// If score is too low, we might not have found a match (optional threshold check)
	if bestProv == "" {
		return "", "", fmt.Errorf("could not match province: %s", detectedProv)
	}

	// 3. Get Kabkota for the matched province
	kabkotas, err := s.GetKabkota(bestProv)
	if err != nil {
		return bestProv, "", err // At least we returned the province
	}

	// 4. Find best matching city
	bestCity := ""
	bestCityScore := -1
	normalizedDetectedCity := normalizeString(detectedCity)

	for _, k := range kabkotas {
		// API usually returns "KAB. X" or "KOTA Y"
		// Nominatim returns simple "X" or "Y" or "Y City"
		// We should match the core name.
		
		score := similarity(normalizedDetectedCity, normalizeString(k))
		if score > bestCityScore {
			bestCityScore = score
			bestCity = k
		}
	}

	if bestCity == "" {
		return bestProv, "", fmt.Errorf("could not match city: %s", detectedCity)
	}

	return bestProv, bestCity, nil
}

// Simple Levenshtein distance-based similarity (0-100)
// Higher is better
func similarity(s1, s2 string) int {
	if s1 == s2 {
		return 100
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0
	}
	
	// Check if one contains the other (very strong signal for City vs KOTA City)
	if contains(s1, s2) || contains(s2, s1) {
		return 90
	}

	d := levenshtein(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}
	
	return 100 - (d * 100 / maxLen)
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func normalizeString(s string) string {
	s = strings.ToLower(s)
	replacer := strings.NewReplacer(
		"provinsi ", "",
		"daerah istimewa ", "",
		"di ", "",
		"kabupaten ", "",
		"kab. ", "",
		"kota ", "",
		" city", "",
		" regency", "",
		" district", "",
	)
	return strings.TrimSpace(replacer.Replace(s))
}

// Levenshtein distance implementation
func levenshtein(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}
	matrix := make([][]int, n+1)
	for i := range matrix {
		matrix[i] = make([]int, m+1)
	}
	for i := 0; i <= n; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= m; j++ {
		matrix[0][j] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			matrix[i][j] = min(min(matrix[i-1][j]+1, matrix[i][j-1]+1), matrix[i-1][j-1]+cost)
		}
	}
	return matrix[n][m]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
