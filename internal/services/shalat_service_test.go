package services_test

import (
	"testing"
	"time"

	"github.com/ramadhan/amaliah-monitoring/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestGetShalat(t *testing.T) {
	s := services.NewShalatService()
	
	// Test case 1: Valid location and date
	provinsi := "Jawa Barat"
	kabkota := "Kota Bandung"
	now := time.Now()
	month := int(now.Month())
	year := now.Year()

	data, err := s.GetShalat(provinsi, kabkota, month, year)
	
	if assert.NoError(t, err) {
		assert.NotNil(t, data)
		if data != nil {
			assert.Equal(t, provinsi, data.Provinsi)
			assert.Equal(t, kabkota, data.Kabkota)
			assert.NotEmpty(t, data.Jadwal)
		}
	}

	// Verify schedule structure
	if len(data.Jadwal) > 0 {
		jadwal := data.Jadwal[0]
		assert.NotEmpty(t, jadwal.Subuh)
		assert.NotEmpty(t, jadwal.Maghrib)
		assert.NotEmpty(t, jadwal.Isya)
	}
}

func TestGetProvinsi(t *testing.T) {
	s := services.NewShalatService()

	provinsi, err := s.GetProvinsi()

	assert.NoError(t, err)
	assert.NotEmpty(t, provinsi)
	assert.Contains(t, provinsi, "Jawa Barat")
	assert.Contains(t, provinsi, "DKI Jakarta")
}

func TestGetKabkota(t *testing.T) {
	s := services.NewShalatService()

	// Test valid province
	kabkota, err := s.GetKabkota("Jawa Barat")
	assert.NoError(t, err)
	assert.NotEmpty(t, kabkota)
	assert.Contains(t, kabkota, "Kota Bandung")

	// Test invalid province (might return empty or error depending on API)
	kabkotaInvalid, _ := s.GetKabkota("InvalidProvince")
	assert.Empty(t, kabkotaInvalid)
}
