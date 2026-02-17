package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const BaseURL = "https://muslim-api-three.vercel.app"

type MuslimAPIService struct {
	client *http.Client
}

func NewMuslimAPIService() *MuslimAPIService {
	return &MuslimAPIService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// API Response types
type APIResponse struct {
	Status int             `json:"status"`
	Data   json.RawMessage `json:"data"`
}

// Doa types
type Doa struct {
	ID     string `json:"_id"`
	Arab   string `json:"arab"`
	Indo   string `json:"indo"`
	Judul  string `json:"judul"`
	Source string `json:"source"`
}

// Hadits types
type Hadits struct {
	ID    string `json:"_id"`
	Arab  string `json:"arab"`
	Indo  string `json:"indo"`
	Judul string `json:"judul"`
	No    string `json:"no"`
}

// Quran types
type Surah struct {
	ID             string `json:"_id"`
	NameShort      string `json:"name_short"`
	NameLong       string `json:"name_long"`
	NameID         string `json:"name_id"`
	NameEN         string `json:"name_en"`
	Number         string `json:"number"`
	NumberOfVerses string `json:"number_of_verses"`
	RevelationID   string `json:"revelation_id"`
	TranslationID  string `json:"translation_id"`
	AudioURL       string `json:"audio_url"`
}

type Ayah struct {
	ID    string `json:"_id"`
	Arab  string `json:"arab"`
	Latin string `json:"latin"`
	Text  string `json:"text"`
	Ayah  string `json:"ayah"`
	Surah string `json:"surah"`
	Juz   string `json:"juz"`
	Page  string `json:"page"`
	Audio string `json:"audio"`
}

// Doa methods
func (s *MuslimAPIService) GetAllDoa() ([]Doa, error) {
	resp, err := s.client.Get(BaseURL + "/v1/doa")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var doaList []Doa
	if err := json.Unmarshal(apiResp.Data, &doaList); err != nil {
		return nil, err
	}

	return doaList, nil
}

func (s *MuslimAPIService) GetDoaBySource(source string) ([]Doa, error) {
	url := BaseURL + "/v1/doa"
	if source != "" {
		url = fmt.Sprintf("%s?source=%s", url, source)
	}

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var doaList []Doa
	if err := json.Unmarshal(apiResp.Data, &doaList); err != nil {
		return nil, err
	}

	return doaList, nil
}

func (s *MuslimAPIService) SearchDoa(query string) ([]Doa, error) {
	url := fmt.Sprintf("%s/v1/doa/find?query=%s", BaseURL, query)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var doaList []Doa
	if err := json.Unmarshal(apiResp.Data, &doaList); err != nil {
		return nil, err
	}

	return doaList, nil
}

// Hadits methods
func (s *MuslimAPIService) GetAllHadits() ([]Hadits, error) {
	resp, err := s.client.Get(BaseURL + "/v1/hadits")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var haditsList []Hadits
	if err := json.Unmarshal(apiResp.Data, &haditsList); err != nil {
		return nil, err
	}

	return haditsList, nil
}

func (s *MuslimAPIService) GetHaditsByNumber(nomor int) (*Hadits, error) {
	url := fmt.Sprintf("%s/v1/hadits?nomor=%d", BaseURL, nomor)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var hadits Hadits
	if err := json.Unmarshal(apiResp.Data, &hadits); err != nil {
		return nil, err
	}

	return &hadits, nil
}

func (s *MuslimAPIService) SearchHadits(query string) ([]Hadits, error) {
	url := fmt.Sprintf("%s/v1/hadits/find?query=%s", BaseURL, query)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var haditsList []Hadits
	if err := json.Unmarshal(apiResp.Data, &haditsList); err != nil {
		return nil, err
	}

	return haditsList, nil
}

// Quran methods
func (s *MuslimAPIService) GetAllSurah() ([]Surah, error) {
	resp, err := s.client.Get(BaseURL + "/v1/quran/surah")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var surahList []Surah
	if err := json.Unmarshal(apiResp.Data, &surahList); err != nil {
		return nil, err
	}

	return surahList, nil
}

func (s *MuslimAPIService) GetSurahByID(id int) (*Surah, error) {
	url := fmt.Sprintf("%s/v1/quran/surah?id=%d", BaseURL, id)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var surah Surah
	if err := json.Unmarshal(apiResp.Data, &surah); err != nil {
		return nil, err
	}

	return &surah, nil
}

func (s *MuslimAPIService) GetAyahBySurah(surahID int) ([]Ayah, error) {
	url := fmt.Sprintf("%s/v1/quran/ayah/surah?id=%d", BaseURL, surahID)

	resp, err := s.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var ayahList []Ayah
	if err := json.Unmarshal(apiResp.Data, &ayahList); err != nil {
		return nil, err
	}

	return ayahList, nil
}
