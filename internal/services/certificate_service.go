package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/ramadhan/amaliah-monitoring/internal/models"
)

type CertificateService struct {
	// We might need config or other services here later
}

func NewCertificateService() *CertificateService {
	return &CertificateService{}
}

func (s *CertificateService) Generate(user *models.User, stats map[string]interface{}) ([]byte, error) {
	pdf := fpdf.New("L", "mm", "A4", "") // Landscape, A4
	pdf.AddPage()

	// Colors
	primaryColor := []int{99, 102, 241} // Indigo
	secondaryColor := []int{79, 70, 229} // Darker Indigo
	// 1. Background Template
	// Load the full page background image (A4 Landscape: 297x210 mm)
	pdf.ImageOptions("web/static/images/sertifikat_bg.png", 0, 0, 297, 210, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// 2. Load Fonts
	pdf.AddUTF8Font("Poppins", "", "web/static/fonts/Poppins-Regular.ttf")
	pdf.AddUTF8Font("Poppins", "B", "web/static/fonts/Poppins-Bold.ttf")

	// 3. Logo (Center Aligned)
	pdf.ImageOptions("web/static/images/logoniba.png", 135, 20, 30, 0, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")

	// 4. Header
	pdf.SetY(60)
	pdf.SetFont("Poppins", "B", 30)
	pdf.SetTextColor(primaryColor[0], primaryColor[1], primaryColor[2])
	pdf.CellFormat(0, 15, "SERTIFIKAT PENGHARGAAN", "", 1, "C", false, 0, "")

	pdf.SetFont("Poppins", "", 16)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 10, "Diberikan kepada:", "", 1, "C", false, 0, "")

	// 5. Student Details
	pdf.SetY(95)
	pdf.SetFont("Poppins", "BU", 24)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(0, 12, user.FullName, "", 1, "C", false, 0, "")

	pdf.SetFont("Poppins", "", 14)
	pdf.SetTextColor(80, 80, 80)
	if user.Class != "" {
		pdf.CellFormat(0, 8, "Kelas "+user.Class, "", 1, "C", false, 0, "")
	}

	// 6. Achievement Text
	pdf.SetY(120)
	pdf.SetFont("Poppins", "", 12)
	pdf.SetTextColor(50, 50, 50)
	pdf.MultiCell(0, 6, "Atas partisipasi aktif dan pencapaian amaliah selama bulan suci Ramadhan 1447 H / 2026 M\ndalam kegiatan Smartren dan Program Monitoring Ibadah Ramadhan.", "", "C", false)

	// 7. Stats Summary Box
	statsY := 145.0
	pdf.SetY(statsY)
	
	// Helper for stats
	drawStat := func(label, value string, x float64) {
		pdf.SetXY(x, statsY)
		pdf.SetFont("Poppins", "B", 14)
		pdf.SetTextColor(secondaryColor[0], secondaryColor[1], secondaryColor[2])
		pdf.CellFormat(50, 8, value, "", 2, "C", false, 0, "")
		
		pdf.SetFont("Poppins", "", 10)
		pdf.SetTextColor(100, 100, 100)
		pdf.CellFormat(50, 5, label, "", 0, "C", false, 0, "")
	}

	totalPoints := "0"
	if tp, ok := stats["total_points"].(int); ok {
		totalPoints = fmt.Sprintf("%d", tp)
	}

	totalPages := "0"
	if tp, ok := stats["total_pages"].(int); ok {
		totalPages = fmt.Sprintf("%d Hal", tp)
	}
	
	khatamProgress := "0%"
	if kp, ok := stats["khatam_percent"].(int); ok {
		khatamProgress = fmt.Sprintf("%d%%", kp)
	}

	// Center stats: Amaliah Points, Total Pages, Quran Progress
	drawStat("Total Poin", totalPoints, 65)
	drawStat("Total Halaman", totalPages, 125)
	drawStat("Progress Khatam", khatamProgress, 185)

	// 7. Footer / Signatures
	pdf.SetY(175)
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(0, 0, 0)
	
	// Date
	now := time.Now()
	dateStr := fmt.Sprintf("Diberikan pada tanggal: %d %s %d", now.Day(), now.Month().String(), now.Year())
	pdf.CellFormat(0, 6, dateStr, "", 1, "C", false, 0, "")

	// 8. Output
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
