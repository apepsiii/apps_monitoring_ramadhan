package services

import (
	"fmt"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/repository"
	"github.com/xuri/excelize/v2"
)

type ExportService struct {
	UserRepo    *repository.UserRepository
	PrayerRepo  *repository.PrayerRepository
	FastingRepo *repository.FastingRepository
	QuranRepo   *repository.QuranRepository
	AmaliahRepo *repository.AmaliahRepository
}

func NewExportService(
	userRepo *repository.UserRepository,
	prayerRepo *repository.PrayerRepository,
	fastingRepo *repository.FastingRepository,
	quranRepo *repository.QuranRepository,
	amaliahRepo *repository.AmaliahRepository,
) *ExportService {
	return &ExportService{
		UserRepo:    userRepo,
		PrayerRepo:  prayerRepo,
		FastingRepo: fastingRepo,
		QuranRepo:   quranRepo,
		AmaliahRepo: amaliahRepo,
	}
}

func (s *ExportService) GenerateDailyReportExcel(date string, className string) (*excelize.File, error) {
	f := excelize.NewFile()
	
	// Create Sheet for Prayer
	sheetName := "Shalat"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1") // Remove default sheet

	// Headers
	headers := []string{"Nama Siswa", "Kelas", "Subuh", "Dzuhur", "Ashar", "Maghrib", "Isya"}
	for i, h := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, h)
	}

	// Get Data
	var users []*models.User
	var err error

	if className != "" {
		users, err = s.UserRepo.GetByClass(className)
	} else {
		users, err = s.UserRepo.GetAll()
	}

	if err != nil {
		return nil, err
	}
	
	row := 2
	for _, user := range users {
		if user.Role == "admin" {
			continue
		}
		
		prayer, _ := s.PrayerRepo.GetByUserAndDate(user.ID, date)
		
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), user.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), user.Class)
		
		if prayer != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), prayer.Subuh)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), prayer.Dzuhur)
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), prayer.Ashar)
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), prayer.Maghrib)
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), prayer.Isya)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), "-")
		}
		row++
	}

	// Create Sheet for Fasting
	sheetName = "Puasa"
	f.NewSheet(sheetName)
	
	headers = []string{"Nama Siswa", "Kelas", "Status", "Alasan"}
	for i, h := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, h)
	}
	
	row = 2
	for _, user := range users {
		if user.Role == "admin" {
			continue
		}
		
		fasting, _ := s.FastingRepo.GetByUserAndDate(user.ID, date)
		
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), user.FullName)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), user.Class)
		
		if fasting != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), fasting.Status)
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fasting.Reason)
		} else {
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "-")
		}
		row++
	}

	return f, nil
}

func (s *ExportService) GenerateStudentReportExcel(userID int, startDate, endDate string) (*excelize.File, error) {
	user, err := s.UserRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	
	// Sheet 1: Summary
	sheetName := "Summary"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	f.SetCellValue(sheetName, "A1", "Laporan Amaliah Ramadhan")
	f.SetCellValue(sheetName, "A3", "Nama:")
	f.SetCellValue(sheetName, "B3", user.FullName)
	f.SetCellValue(sheetName, "A4", "Kelas:")
	f.SetCellValue(sheetName, "B4", user.Class)
	f.SetCellValue(sheetName, "A5", "Periode:")
	f.SetCellValue(sheetName, "B5", fmt.Sprintf("%s sd %s", startDate, endDate))

	// Sheet 2: Shalat
	sheetName = "Shalat"
	f.NewSheet(sheetName)
	
	headers := []string{"Tanggal", "Subuh", "Dzuhur", "Ashar", "Maghrib", "Isya"}
	for i, h := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, h)
	}

	prayers, _ := s.PrayerRepo.GetByUserAndDateRange(userID, startDate, endDate)
	row := 2
	for _, p := range prayers {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.Date)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), p.Subuh)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), p.Dzuhur)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), p.Ashar)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), p.Maghrib)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), p.Isya)
		row++
	}

	// Sheet 3: Quran
	sheetName = "Al-Quran"
	f.NewSheet(sheetName)
	
	headers = []string{"Tanggal", "Mulai", "Selesai", "Catatan"}
	for i, h := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, h)
	}

	readings, _ := s.QuranRepo.GetByDateRange(userID, startDate, endDate)
	row = 2
	for _, r := range readings {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), r.Date)
		start := fmt.Sprintf("%s: %d", r.StartSurahName, r.StartAyah)
		end := fmt.Sprintf("%s: %d", r.EndSurahName, r.EndAyah)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), start)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), end)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), r.Notes)
		row++
	}

	return f, nil
}
