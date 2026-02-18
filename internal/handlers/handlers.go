package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/repository"
	"github.com/ramadhan/amaliah-monitoring/internal/services"
	"github.com/ramadhan/amaliah-monitoring/internal/utils"
	"github.com/xuri/excelize/v2"
)

type Handler struct {
	DB               *sql.DB
	UserRepo         *repository.UserRepository
	PrayerRepo       *repository.PrayerRepository
	FastingRepo      *repository.FastingRepository
	QuranRepo        *repository.QuranRepository
	AmaliahRepo      *repository.AmaliahRepository
	MuslimAPI        *services.MuslimAPIService
	ImsakiyahService *services.ImsakiyahService
	ShalatService    *services.ShalatService
	AdminService     *services.AdminService
	ExportService    *services.ExportService
	BadgeRepo        *repository.BadgeRepository
	BadgeService     *services.BadgeService
	StatisticsService *services.StatisticsService
	CertificateService *services.CertificateService
	ClassRepo        *repository.ClassRepository
}

func NewHandler(db *sql.DB) *Handler {
	userRepo := repository.NewUserRepository(db)
	prayerRepo := repository.NewPrayerRepository(db)
	fastingRepo := repository.NewFastingRepository(db)
	quranRepo := repository.NewQuranRepository(db)
	amaliahRepo := repository.NewAmaliahRepository(db)
	badgeRepo := repository.NewBadgeRepository(db)
	classRepo := repository.NewClassRepository(db)

	return &Handler{
		DB:               db,
		UserRepo:         userRepo,
		PrayerRepo:       prayerRepo,
		FastingRepo:      fastingRepo,
		QuranRepo:        quranRepo,
		AmaliahRepo:      amaliahRepo,
		MuslimAPI:        services.NewMuslimAPIService(),
		ImsakiyahService: services.NewImsakiyahService(),
		ShalatService:    services.NewShalatService(),
		AdminService:     services.NewAdminService(userRepo),
		ExportService:    services.NewExportService(userRepo, prayerRepo, fastingRepo, quranRepo, amaliahRepo),
		BadgeRepo:        badgeRepo,
		BadgeService:     services.NewBadgeService(badgeRepo, prayerRepo, amaliahRepo, quranRepo),
		StatisticsService: services.NewStatisticsService(prayerRepo, amaliahRepo, fastingRepo, userRepo),
		CertificateService: services.NewCertificateService(),
		ClassRepo:        classRepo,
	}
}

// Home - Landing Page
func (h *Handler) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"Title": "Selamat Datang",
	})
}

// Jadwal Shalat & Imsakiyah - Public Page
func (h *Handler) ShowJadwal(c echo.Context) error {
	provinsi := c.QueryParam("provinsi")
	kabkota := c.QueryParam("kabkota")
	tab := c.QueryParam("tab")
	if tab == "" {
		tab = "shalat"
	}

	provinsiList, _ := h.ShalatService.GetProvinsi()

	var kabkotaList []string
	if provinsi != "" {
		kabkotaList, _ = h.ShalatService.GetKabkota(provinsi)
	}

	var shalatData *models.ShalatData
	var imsakiyahData *models.ImsakiyahData
	var todayShalat *models.ShalatSchedule
	var todayImsakiyah *models.ImsakiyahSchedule

	if provinsi != "" && kabkota != "" {
		now := time.Now()
		shalatData, _ = h.ShalatService.GetShalat(provinsi, kabkota, int(now.Month()), now.Year())
		imsakiyahData, _ = h.ImsakiyahService.GetImsakiyah(provinsi, kabkota)

		if shalatData != nil {
			dayOfMonth := now.Day()
			for _, s := range shalatData.Jadwal {
				if s.Tanggal == dayOfMonth {
					todayShalat = &s
					break
				}
			}
		}

		if imsakiyahData != nil {
			dayOfMonth := now.Day()
			for _, s := range imsakiyahData.Imsakiyah {
				if s.Tanggal == dayOfMonth {
					todayImsakiyah = &s
					break
				}
			}
		}
	}

	data := map[string]interface{}{
		"Title":          "Jadwal Shalat & Imsakiyah",
		"Provinsi":       provinsi,
		"Kabkota":        kabkota,
		"ProvinsiList":   provinsiList,
		"KabkotaList":    kabkotaList,
		"Tab":            tab,
		"ShalatData":     shalatData,
		"ImsakiyahData":  imsakiyahData,
		"TodayShalat":    todayShalat,
		"TodayImsakiyah": todayImsakiyah,
	}

	user := c.Get("user")
	if user != nil {
		data["User"] = user
	}

	return c.Render(http.StatusOK, "jadwal.html", data)
}

// Jadwal Shalat & Imsakiyah - Protected Page for Logged-in Users
func (h *Handler) ShowUserJadwal(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Get location from query params or use user's saved location
	provinsi := c.QueryParam("provinsi")
	kabkota := c.QueryParam("kabkota")
	tab := c.QueryParam("tab")
	if tab == "" {
		tab = "shalat"
	}

	// Use user's saved location if not provided in query
	if provinsi == "" && user.Provinsi != "" {
		provinsi = user.Provinsi
	}
	if kabkota == "" && user.Kabkota != "" {
		kabkota = user.Kabkota
	}

	provinsiList, _ := h.ShalatService.GetProvinsi()

	var kabkotaList []string
	if provinsi != "" {
		kabkotaList, _ = h.ShalatService.GetKabkota(provinsi)
	}

	var shalatData *models.ShalatData
	var imsakiyahData *models.ImsakiyahData
	var todayShalat *models.ShalatSchedule
	var todayImsakiyah *models.ImsakiyahSchedule

	if provinsi != "" && kabkota != "" {
		now := time.Now()
		shalatData, _ = h.ShalatService.GetShalat(provinsi, kabkota, int(now.Month()), now.Year())
		imsakiyahData, _ = h.ImsakiyahService.GetImsakiyah(provinsi, kabkota)

		if shalatData != nil {
			dayOfMonth := now.Day()
			for _, s := range shalatData.Jadwal {
				if s.Tanggal == dayOfMonth {
					todayShalat = &s
					break
				}
			}
		}

		if imsakiyahData != nil {
			dayOfMonth := now.Day()
			for _, s := range imsakiyahData.Imsakiyah {
				if s.Tanggal == dayOfMonth {
					todayImsakiyah = &s
					break
				}
			}
		}
	}

	data := map[string]interface{}{
		"Title":          "Jadwal Shalat & Imsakiyah",
		"User":           user,
		"Provinsi":       provinsi,
		"Kabkota":        kabkota,
		"ProvinsiList":   provinsiList,
		"KabkotaList":    kabkotaList,
		"Tab":            tab,
		"ShalatData":     shalatData,
		"ImsakiyahData":  imsakiyahData,
		"TodayShalat":    todayShalat,
		"TodayImsakiyah": todayImsakiyah,
	}

	return c.Render(http.StatusOK, "user/jadwal.html", data)
}

// Auth Handlers
func (h *Handler) ShowLogin(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/login.html", map[string]interface{}{
		"Title": "Masuk",
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.Render(http.StatusOK, "auth/login.html", map[string]interface{}{
			"Title": "Masuk",
			"Error": "Invalid request",
		})
	}

	user, err := h.UserRepo.GetByUsername(req.Username)
	if err != nil {
		return c.Render(http.StatusOK, "auth/login.html", map[string]interface{}{
			"Title": "Masuk",
			"Error": "Username atau password salah",
		})
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return c.Render(http.StatusOK, "auth/login.html", map[string]interface{}{
			"Title": "Masuk",
			"Error": "Username atau password salah",
		})
	}

	// Generate JWT
	token, err := utils.GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	// Set cookie
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
	if os.Getenv("APP_ENV") == "production" {
		cookie.Secure = true
	}
	c.SetCookie(cookie)

	// superadmin = system admin, goes to system admin panel
	// admin = school admin, goes to user dashboard (where school management is)
	// user = regular user, goes to user dashboard
	if user.Role == "superadmin" {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}
	return c.Redirect(http.StatusSeeOther, "/user/dashboard")
}

func (h *Handler) ShowRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/register.html", map[string]interface{}{
		"Title": "Daftar",
	})
}

func (h *Handler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.Render(http.StatusOK, "auth/register.html", map[string]interface{}{
			"Title": "Daftar",
			"Error": "Invalid request",
		})
	}

	if len(req.Password) < 6 {
		return c.Render(http.StatusOK, "auth/register.html", map[string]interface{}{
			"Title": "Daftar",
			"Error": "Password minimal 6 karakter",
		})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	role := "user"
	schoolID := 0

	// Handle School Logic
	if req.NewSchoolName != "" {
		// Create New School
		// Generate simple code: First word + random number? Or just random 6 chars
		// For simplicity/robustness: Random 6 chars
		schoolCode := utils.GenerateRandomString(6)

		// Create School Entry
		res, err := h.DB.Exec("INSERT INTO schools (name, code, address) VALUES (?, ?, ?)", req.NewSchoolName, schoolCode, "-")
		if err != nil {
			return c.Render(http.StatusOK, "auth/register.html", map[string]interface{}{
				"Title": "Daftar",
				"Error": "Gagal membuat sekolah baru",
			})
		}
		sid, _ := res.LastInsertId()
		schoolID = int(sid)
		role = "admin" // Creator becomes admin

	} else if req.SchoolCode != "" {
		// Join Existing School
		var sid int
		err := h.DB.QueryRow("SELECT id FROM schools WHERE code = ?", req.SchoolCode).Scan(&sid)
		if err != nil {
			return c.Render(http.StatusOK, "auth/register.html", map[string]interface{}{
				"Title": "Daftar",
				"Error": "Kode sekolah tidak valid",
			})
		}
		schoolID = sid
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Class:        req.Class,
		Role:         role,
		Points:       0,
		SchoolID:     schoolID,
	}

	if err := h.UserRepo.Create(user); err != nil {
		return c.Render(http.StatusOK, "auth/register.html", map[string]interface{}{
			"Title": "Daftar",
			"Error": "Gagal mendaftar: username atau email sudah terdaftar",
		})
	}

	// If Created School, update AdminID
	if req.NewSchoolName != "" {
		createdUser, _ := h.UserRepo.GetByUsername(req.Username)
		if createdUser != nil {
			_, _ = h.DB.Exec("UPDATE schools SET admin_id = ? WHERE id = ?", createdUser.ID, schoolID)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (h *Handler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}

// User Handlers
func (h *Handler) UserDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Superadmin should never be here — redirect to admin panel
	if user.Role == "superadmin" {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	}

	today := time.Now().Format("2006-01-02")
	prayer, _ := h.PrayerRepo.GetByUserAndDate(user.ID, today)

	prayerCompleted := 0
	if prayer != nil {
		if prayer.Subuh == "jamaah" || prayer.Subuh == "sendiri" {
			prayerCompleted++
		}
		if prayer.Dzuhur == "jamaah" || prayer.Dzuhur == "sendiri" {
			prayerCompleted++
		}
		if prayer.Ashar == "jamaah" || prayer.Ashar == "sendiri" {
			prayerCompleted++
		}
		if prayer.Maghrib == "jamaah" || prayer.Maghrib == "sendiri" {
			prayerCompleted++
		}
		if prayer.Isya == "jamaah" || prayer.Isya == "sendiri" {
			prayerCompleted++
		}
	}

	fasting, _ := h.FastingRepo.GetTodayFasting(user.ID)

	todayAmaliah, _ := h.AmaliahRepo.GetDailyAmaliah(user.ID, today)
	todayCompleted := len(todayAmaliah)

	todayPoints, _ := h.AmaliahRepo.GetTodayPoints(user.ID)

	totalReadings, _ := h.QuranRepo.GetTotalReadings(user.ID)
	
	// NEW: Get Total Fasting
	totalFasting, _ := h.FastingRepo.GetTotalFasting(user.ID)

	prayerStreak, bestPrayer, _ := h.PrayerRepo.GetPrayerStreak(user.ID)
	fastingStreak, bestFasting, _ := h.FastingRepo.GetFastingStreak(user.ID)
	quranStreak, bestQuran, _ := h.QuranRepo.GetQuranStreak(user.ID)
	amaliahStreak, bestAmaliah, _ := h.AmaliahRepo.GetAmaliahStreak(user.ID)

	streak := &models.Streak{
		UserID:        user.ID,
		PrayerStreak:  prayerStreak,
		FastingStreak: fastingStreak,
		QuranStreak:   quranStreak,
		BestPrayer:    bestPrayer,
		BestFasting:   bestFasting,
		BestQuran:     bestQuran,
		AmaliahStreak: amaliahStreak,
		BestAmaliah:   bestAmaliah,
	}

	var todaySchedule *models.ImsakiyahSchedule
	var imsakiyahData *models.ImsakiyahData
	if user.Provinsi != "" && user.Kabkota != "" {
		dayOfMonth := time.Now().Day()
		todaySchedule, _ = h.ImsakiyahService.GetTodaySchedule(user.Provinsi, user.Kabkota, dayOfMonth)
		imsakiyahData, _ = h.ImsakiyahService.GetImsakiyah(user.Provinsi, user.Kabkota)
	}

	// Check and get badges
	newBadges, _ := h.BadgeService.CheckAndAwardBadges(user.ID)
	userBadges, _ := h.BadgeRepo.GetUserBadges(user.ID)

	// NEW: Prepare Chart Data (Quran Monthly Progress)
	type DailyQuranStat struct {
		Day    int
		Target int
		Actual int
	}
	var quranChartData []DailyQuranStat
	
	now := time.Now()
	// Get start and end of current month
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	nextMonth := now.AddDate(0, 1, 0)
	lastDayOfMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 0, 0, 0, 0, nextMonth.Location())
	
	startOfMonthStr := startOfMonth.Format("2006-01-02") 
	endOfMonthStr := lastDayOfMonth.Format("2006-01-02")

	// Get readings for the month
	monthReadings, _ := h.QuranRepo.GetByDateRange(user.ID, startOfMonthStr, endOfMonthStr)
	
	// Map readings to date
	readingsMap := make(map[int]int) // day -> pages
	for _, r := range monthReadings {
		t, _ := time.Parse("2006-01-02", r.Date)
		readingsMap[t.Day()] += r.Pages
	}

	// Calculate Daily Target
	targetKhatam := user.TargetKhatam
	if targetKhatam <= 0 {
		targetKhatam = 30 // Default 30 days
	}
	dailyTarget := int(math.Ceil(604.0 / float64(targetKhatam)))

	// Fill chart data for all days in month
	daysInMonth := lastDayOfMonth.Day()
	for d := 1; d <= daysInMonth; d++ {
		quranChartData = append(quranChartData, DailyQuranStat{
			Day:    d,
			Target: dailyTarget,
			Actual: readingsMap[d],
		})
	}

	// Format Date for Dashboard header
	indonesianMonths := []string{"", "Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	currentDate := time.Now()
	formattedDate := fmt.Sprintf("%d %s %d", currentDate.Day(), indonesianMonths[currentDate.Month()], currentDate.Year())
	
	location := user.Kabkota
	if location == "" {
		location = user.Provinsi
	}
	if location == "" {
		location = "Indonesia"
	}
	dashboardDate := fmt.Sprintf("%s, %s", location, formattedDate)

	// Get School Name, Code, and Status
	schoolName := ""
	schoolCode := ""
	schoolPending := false
	if user.SchoolID > 0 {
		var name, code, status string
		err := h.DB.QueryRow("SELECT name, code, status FROM schools WHERE id = ?", user.SchoolID).Scan(&name, &code, &status)
		if err == nil {
			schoolName = name
			schoolCode = code
			schoolPending = status == "pending"
		}
	}


	return c.Render(http.StatusOK, "user/dashboard.html", map[string]interface{}{
		"Title":           "Dashboard",
		"User":            user,
		"Prayer":          prayer,
		"PrayerCompleted": prayerCompleted,
		"Fasting":         fasting,
		"TodayPoints":     todayPoints,
		"TodayCompleted":  todayCompleted,
		"TotalReadings":   totalReadings,
		"TotalFasting":    totalFasting,
		"Streak":          streak,
		"TodaySchedule":   todaySchedule,
		"ImsakiyahData":   imsakiyahData,
		"UserBadges":      userBadges,
		"NewBadges":       newBadges,
		"QuranChartData":  quranChartData,
		"DashboardDate":   dashboardDate,
		"SchoolName":      schoolName,
		"SchoolCode":      schoolCode,
		"SchoolPending":   schoolPending,
	})
}

func (h *Handler) ShowPrayers(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Get today's prayer or create default
	today := time.Now()
	todayStr := today.Format("2006-01-02")
	prayer, err := h.PrayerRepo.GetByUserAndDate(user.ID, todayStr)
	if err != nil {
		// Create default prayer entry
		prayer = &models.Prayer{
			UserID:  user.ID,
			Date:    todayStr,
			Subuh:   "belum",
			Dzuhur:  "belum",
			Ashar:   "belum",
			Maghrib: "belum",
			Isya:    "belum",
		}
	}

	// Get last 7 days prayer stats
	weekAgo := today.AddDate(0, 0, -6).Format("2006-01-02")
	weekPrayers, _ := h.PrayerRepo.GetByUserAndDateRange(user.ID, weekAgo, todayStr)

	// Get month stats
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
	monthStats, _ := h.PrayerRepo.GetPrayerStats(user.ID, startOfMonth, todayStr)

	// Format today's date for display
	months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	todayFormatted := days[int(today.Weekday())] + ", " + strconv.Itoa(today.Day()) + " " + months[today.Month()-1] + " " + strconv.Itoa(today.Year())

	return c.Render(http.StatusOK, "user/prayers.html", map[string]interface{}{
		"Title":        "Shalat",
		"User":         user,
		"Prayer":       prayer,
		"WeekPrayers":  weekPrayers,
		"MonthStats":   monthStats,
		"TodayDate":    todayFormatted,
		"TodayDateISO": todayStr,
		"Error":        c.QueryParam("error"),
		"Success":      c.QueryParam("success"),
	})
}

func (h *Handler) SavePrayers(c echo.Context) error {
	user := c.Get("user").(*models.User)

	date := c.FormValue("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	subuh := c.FormValue("subuh")
	dzuhur := c.FormValue("dzuhur")
	ashar := c.FormValue("ashar")
	maghrib := c.FormValue("maghrib")
	isya := c.FormValue("isya")

	err := h.PrayerRepo.CreateOrUpdate(user.ID, date, subuh, dzuhur, ashar, maghrib, isya)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save prayer"})
	}

	return c.Redirect(http.StatusSeeOther, "/user/prayers")
}

func (h *Handler) ShowFasting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	today := time.Now()
	todayStr := today.Format("2006-01-02")

	// Get or create today's fasting
	fasting, err := h.FastingRepo.GetTodayFasting(user.ID)
	if err != nil {
		fasting = &models.Fasting{
			UserID: user.ID,
			Date:   todayStr,
			Status: "puasa",
			Reason: "",
		}
	}

	// Get current month fasting data
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	startOfMonthStr := startOfMonth.Format("2006-01-02")
	endOfMonthStr := utils.GetEndOfMonth()

	monthFastings, _ := h.FastingRepo.GetByUserAndDateRange(user.ID, startOfMonthStr, endOfMonthStr)

	// Create fasting map for quick lookup
	fastingMap := make(map[string]*models.Fasting)
	for _, f := range monthFastings {
		fastingMap[f.Date] = f
	}

	// Build calendar days
	// Build calendar days
	type CalendarDay struct {
		Day     int
		Date    string
		HasData bool
		Status  string
		Reason  string
		IsToday bool
		IsPast  bool
	}

	var calendarDays []CalendarDay
	emptyDays := int(startOfMonth.Weekday())

	// Get last day of month
	nextMonth := today.AddDate(0, 1, 0)
	lastDay := time.Date(nextMonth.Year(), nextMonth.Month(), 0, 0, 0, 0, 0, nextMonth.Location()).Day()

    // Calculate stats manually for "Assumed Puasa"
    notFastingCount := 0
    daysPassed := 0
    if today.Month() == startOfMonth.Month() {
        daysPassed = today.Day()
    } else if today.After(startOfMonth) {
         // If today is after this month (viewing past month?), handle logic. 
         // But here startOfMonth is based on today. So we are viewing current month.
         daysPassed = today.Day()
    }

	for day := 1; day <= lastDay; day++ {
		dateStr := time.Date(today.Year(), today.Month(), day, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
        isPast := dateStr < todayStr
        isToday := dateStr == todayStr // Logic check: today is dateStr?

		dayData := CalendarDay{
			Day:     day,
			Date:    dateStr,
			HasData: false,
			IsToday: isToday,
			IsPast:  isPast,
		}

		if f, exists := fastingMap[dateStr]; exists {
			dayData.HasData = true
			dayData.Status = f.Status
			dayData.Reason = f.Reason
            
            if f.Status == "tidak" && day <= daysPassed {
                notFastingCount++
            }
		} else if isPast {
             // AUTO-FILL 'Puasa' for past days
             dayData.HasData = true
             dayData.Status = "puasa"
             dayData.Reason = ""
        }

		calendarDays = append(calendarDays, dayData)
	}

	// Override DB stats with "Assumed" stats
    // Fasting = DaysPassed - NotFasting
    fastingCount := daysPassed - notFastingCount
    if fastingCount < 0 { fastingCount = 0 }
    
    // Total days to show in stats (Progress)
    // Usually total days of Ramadhan (30)? Or days passed?
    // Template shows "Stats.fasting / Stats.total_days". 
    // If total_days is 30, then progress bar works.
    // If total_days is daysPassed, it's 100%.
    // DB GetFastingStats returned "total_days" as count(*) of rows.
    
    // Let's use 30 (or lastDay) as total RAMADHAN days for context, 
    // OR just days passed for accurate "Progress so far".
    // User request: "Progress Puasa... belum nandain".
    // Let's update stats to reflect "Days Passed".
    
	stats := map[string]int{
        "fasting": fastingCount,
        "not_fasting": notFastingCount,
        "total_days": daysPassed, 
    }

	// Format today's date for display
	months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	todayFormatted := days[int(today.Weekday())] + ", " + strconv.Itoa(today.Day()) + " " + months[today.Month()-1] + " " + strconv.Itoa(today.Year())

	return c.Render(http.StatusOK, "user/fasting.html", map[string]interface{}{
		"Title":        "Puasa",
		"User":         user,
		"Fasting":      fasting,
		"CalendarDays": calendarDays,
		"EmptyDays":    emptyDays,
		"Stats":        stats,
		"TodayDate":    todayFormatted,
		"TodayDateISO": todayStr,
		"Error":        c.QueryParam("error"),
		"Success":      c.QueryParam("success"),
	})
}

func (h *Handler) SaveFasting(c echo.Context) error {
	user := c.Get("user").(*models.User)

	date := c.FormValue("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	status := c.FormValue("status")
	reason := c.FormValue("reason")

	err := h.FastingRepo.CreateOrUpdate(user.ID, date, status, reason)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save fasting"})
	}

	return c.Redirect(http.StatusSeeOther, "/user/fasting")
}

func (h *Handler) ShowQuran(c echo.Context) error {
	user := c.Get("user").(*models.User)

	today := time.Now()
	todayStr := today.Format("2006-01-02")

	// Get recent readings
	readings, _ := h.QuranRepo.GetByUser(user.ID, 10)

	// Get total readings count
	totalReadings, _ := h.QuranRepo.GetTotalReadings(user.ID)
	totalPages, _ := h.QuranRepo.GetTotalPagesRead(user.ID)

	// Calculate Progress & Target
	targetKhatamDays := user.TargetKhatam
	if targetKhatamDays <= 0 {
		targetKhatamDays = 30
	}
	targetDailyPages := float64(604) / float64(targetKhatamDays)
	progressPercent := float64(totalPages) / 604.0 * 100

	// Get today's readings
	todayReadings, _ := h.QuranRepo.GetByUserAndDate(user.ID, todayStr)
	todayPages := 0
	for _, r := range todayReadings {
		todayPages += r.Pages
	}

	// Get all surah from API for dropdown
	surahList, _ := h.MuslimAPI.GetAllSurah()

	// Format today's date
	months := []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	days := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}
	todayFormatted := days[int(today.Weekday())] + ", " + strconv.Itoa(today.Day()) + " " + months[today.Month()-1] + " " + strconv.Itoa(today.Year())

	return c.Render(http.StatusOK, "user/quran.html", map[string]interface{}{
		"Title":            "Al-Quran",
		"User":             user,
		"Readings":         readings,
		"TotalReadings":    totalReadings,
		"TotalPages":       totalPages,
		"ProgressPercent":  int(progressPercent),
		"TargetDailyPages": int(math.Ceil(targetDailyPages)),
		"TodayPages":       todayPages,
		"TodayReadings":    todayReadings,
		"TodayDate":        todayFormatted,
		"SurahList":        surahList,
		"Error":            c.QueryParam("error"),
		"Success":          c.QueryParam("success"),
	})
}

func (h *Handler) SaveQuran(c echo.Context) error {
	user := c.Get("user").(*models.User)

	startSurahID, _ := strconv.Atoi(c.FormValue("start_surah_id"))
	startSurahName := c.FormValue("start_surah_name")
	startAyah, _ := strconv.Atoi(c.FormValue("start_ayah"))
	endSurahID, _ := strconv.Atoi(c.FormValue("end_surah_id"))
	endSurahName := c.FormValue("end_surah_name")
	endAyah, _ := strconv.Atoi(c.FormValue("end_ayah"))
	pages, _ := strconv.Atoi(c.FormValue("pages"))
	notes := c.FormValue("notes")

	// Validation
	if startSurahID < 1 || startSurahID > 114 {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Surah awal harus antara 1-114")
	}
	if endSurahID < 1 || endSurahID > 114 {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Surah akhir harus antara 1-114")
	}
	if startAyah < 1 {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Ayat awal minimal 1")
	}
	if endAyah < 1 {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Ayat akhir minimal 1")
	}
	if startSurahID > endSurahID {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Surah awal tidak boleh lebih besar dari surah akhir")
	}
	if startSurahID == endSurahID && startAyah > endAyah {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Ayat awal tidak boleh lebih besar dari ayat akhir pada surah yang sama")
	}

	reading := &models.QuranReading{
		UserID:         user.ID,
		Date:           time.Now().Format("2006-01-02"),
		StartSurahID:   startSurahID,
		StartSurahName: startSurahName,
		StartAyah:      startAyah,
		EndSurahID:     endSurahID,
		EndSurahName:   endSurahName,
		EndAyah:        endAyah,
		Pages:          pages,
		Notes:          notes,
	}

	err := h.QuranRepo.Create(reading)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Gagal menyimpan bacaan")
	}

	return c.Redirect(http.StatusSeeOther, "/user/quran?success=Bacaan berhasil disimpan")
}

func (h *Handler) DeleteQuran(c echo.Context) error {
	user := c.Get("user").(*models.User)

	readingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=ID tidak valid")
	}

	// Get the reading to verify ownership
	readings, _ := h.QuranRepo.GetByUser(user.ID, 100)
	var readingToDelete *models.QuranReading
	for _, r := range readings {
		if r.ID == readingID {
			readingToDelete = r
			break
		}
	}

	if readingToDelete == nil {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Data tidak ditemukan")
	}

	err = h.QuranRepo.Delete(readingID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/quran?error=Gagal menghapus bacaan")
	}

	return c.Redirect(http.StatusSeeOther, "/user/quran?success=Bacaan berhasil dihapus")
}

func (h *Handler) ShowAmaliah(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Get all amaliah types
	types, _ := h.AmaliahRepo.GetAllTypes()

	// Get today's amaliah
	today := time.Now().Format("2006-01-02")
	todayAmaliah, _ := h.AmaliahRepo.GetDailyAmaliah(user.ID, today)

	// Create completed map
	completedMap := make(map[int]bool)
	for _, ta := range todayAmaliah {
		completedMap[ta.AmaliahTypeID] = true
	}

	// Count completed
	completedCount := len(todayAmaliah)

	// Get today's points
	todayPoints, _ := h.AmaliahRepo.GetTodayPoints(user.ID)

	// Get leaderboard
	leaderboard, _ := h.AmaliahRepo.GetLeaderboard(10)

	return c.Render(http.StatusOK, "user/amaliah.html", map[string]interface{}{
		"Title":          "Amaliah",
		"User":           user,
		"Types":          types,
		"TodayAmaliah":   todayAmaliah,
		"CompletedMap":   completedMap,
		"CompletedCount": completedCount,
		"TodayPoints":    todayPoints,
		"Leaderboard":    leaderboard,
		"Error":          c.QueryParam("error"),
		"Success":        c.QueryParam("success"),
	})
}

func (h *Handler) SaveAmaliah(c echo.Context) error {
	user := c.Get("user").(*models.User)

	amaliahTypeID, _ := strconv.Atoi(c.FormValue("amaliah_type_id"))
	notes := c.FormValue("notes")
	action := c.FormValue("action") // "add" or "remove"

	if action == "remove" {
		// Remove amaliah
		today := time.Now().Format("2006-01-02")
		item, err := h.AmaliahRepo.GetDailyAmaliahByType(user.ID, amaliahTypeID, today)
		if err == nil {
			h.AmaliahRepo.DeleteDailyAmaliah(item.ID)

			// Update user points
			amaliahType, _ := h.AmaliahRepo.GetTypeByID(amaliahTypeID)
			if amaliahType != nil {
				h.UserRepo.UpdatePoints(user.ID, -amaliahType.Points)
			}
		}
	} else {
		// Add amaliah
		// Check if already exists for today
		today := time.Now().Format("2006-01-02")
		_, err := h.AmaliahRepo.GetDailyAmaliahByType(user.ID, amaliahTypeID, today)
		if err == nil {
			// Already exists, do not add again
			return c.Redirect(http.StatusSeeOther, "/user/amaliah")
		}

		da := &models.DailyAmaliah{
			UserID:        user.ID,
			AmaliahTypeID: amaliahTypeID,
			Date:          today,
			Notes:         notes,
		}

		err = h.AmaliahRepo.CreateDailyAmaliah(da)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save amaliah"})
		}

		// Update user points
		amaliahType, _ := h.AmaliahRepo.GetTypeByID(amaliahTypeID)
		if amaliahType != nil {
			h.UserRepo.UpdatePoints(user.ID, amaliahType.Points)
		}
	}

	return c.Redirect(http.StatusSeeOther, "/user/amaliah")
}

// Admin Handlers
func (h *Handler) AdminDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Get real stats from database
	userStats, _ := h.UserRepo.GetStats()
	today := time.Now().Format("2006-01-02")
	activeUsers, _ := h.UserRepo.GetActiveUsersCount(today)

	// Get today's activity stats
	prayerStats, _ := h.PrayerRepo.GetTodayStats(today)
	fastingStats, _ := h.FastingRepo.GetTodayStats(today)
	quranStats, _ := h.QuranRepo.GetTodayStats(today)
	amaliahStats, _ := h.AmaliahRepo.GetTodayStats(today)

	// Get top users
	topUsers, _ := h.AmaliahRepo.GetLeaderboard(5)

	// Calculate percentages
	totalUsers := userStats["user_count"].(int)
	var prayerPercentage, fastingPercentage int
	if totalUsers > 0 {
		prayerPercentage = (prayerStats["total_users"] * 100) / totalUsers
		fastingPercentage = (fastingStats["fasting"] * 100) / totalUsers
	}

	// Get charts data
	dashboardStats, _ := h.StatisticsService.GetDashboardStats()

	// Get school list with member count
	type SchoolStat struct {
		ID      int
		Name    string
		Code    string
		Members int
	}
	var schools []SchoolStat
	rows, err := h.DB.Query(`
		SELECT s.id, s.name, s.code, COUNT(u.id) as member_count
		FROM schools s
		LEFT JOIN users u ON u.school_id = s.id
		GROUP BY s.id
		ORDER BY s.name ASC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var sc SchoolStat
			if err := rows.Scan(&sc.ID, &sc.Name, &sc.Code, &sc.Members); err == nil {
				schools = append(schools, sc)
			}
		}
	}

	// Get pending admin registration requests
	type PendingSchool struct {
		SchoolID   int
		SchoolName string
		AdminName  string
		Phone      string
		SchoolLevel string
		StudentCount int
	}
	var pendingSchools []PendingSchool
	pendingRows, err := h.DB.Query(`
		SELECT id, school_name, full_name, phone, school_level, student_count
		FROM admin_requests
		WHERE status = 'pending'
		ORDER BY created_at ASC
	`)
	if err == nil {
		defer pendingRows.Close()
		for pendingRows.Next() {
			var ps PendingSchool
			if err := pendingRows.Scan(&ps.SchoolID, &ps.SchoolName, &ps.AdminName, &ps.Phone, &ps.SchoolLevel, &ps.StudentCount); err == nil {
				pendingSchools = append(pendingSchools, ps)
			}
		}
	}

	return c.Render(http.StatusOK, "admin/dashboard.html", map[string]interface{}{
		"Title":             "Admin Dashboard",
		"User":              user,
		"UserStats":         userStats,
		"ActiveUsers":       activeUsers,
		"PrayerStats":       prayerStats,
		"FastingStats":      fastingStats,
		"QuranStats":        quranStats,
		"AmaliahStats":      amaliahStats,
		"TopUsers":          topUsers,
		"PrayerPercentage":  prayerPercentage,
		"FastingPercentage": fastingPercentage,
		"DashboardStats":    dashboardStats,
		"Schools":           schools,
		"PendingSchools":    pendingSchools,
		"Success":           c.QueryParam("success"),
		"Error":             c.QueryParam("error"),
	})
}

func (h *Handler) ManageUsers(c echo.Context) error {
	user := c.Get("user").(*models.User)

	users, _ := h.UserRepo.GetAll()
	classes, _ := h.ClassRepo.GetAll()

	return c.Render(http.StatusOK, "admin/users.html", map[string]interface{}{
		"Title":   "Kelola Siswa",
		"User":    user,
		"Users":   users,
		"Classes": classes,
		"Success": c.QueryParam("success"),
		"Error":   c.QueryParam("error"),
	})
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	newUser := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Class:        req.Class,
		Role:         "user",
	}

	if err := h.UserRepo.Create(newUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.Redirect(http.StatusSeeOther, "/admin/users")
}

// Download Template Excel untuk Import User
func (h *Handler) DownloadUserTemplate(c echo.Context) error {
	f := excelize.NewFile()
	defer f.Close()

	// Create sheet
	sheetName := "Template User"
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	// Set headers
	headers := []string{"Username", "Email", "Password", "Nama Lengkap", "Kelas"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// Add sample data
	f.SetCellValue(sheetName, "A2", "ahmad123")
	f.SetCellValue(sheetName, "B2", "ahmad@example.com")
	f.SetCellValue(sheetName, "C2", "password123")
	f.SetCellValue(sheetName, "D2", "Ahmad Fauzi")
	f.SetCellValue(sheetName, "E2", "XII IPA 1")

	// Style header
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"0D7E5E"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheetName, "A1", "E1", style)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 15)
	f.SetColWidth(sheetName, "B", "B", 25)
	f.SetColWidth(sheetName, "C", "C", 15)
	f.SetColWidth(sheetName, "D", "D", 25)
	f.SetColWidth(sheetName, "E", "E", 15)

	// Set content type and send file
	c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=template_user.xlsx")

	return f.Write(c.Response().Writer)
}



func (h *Handler) ShowReports(c echo.Context) error {
	user := c.Get("user").(*models.User)
	return c.Render(http.StatusOK, "admin/reports.html", map[string]interface{}{
		"Title": "Laporan",
		"User":  user,
	})
}

func (h *Handler) ShowStatistics(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Get date range from query params or use current date
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")

	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Get statistics
	userStats, _ := h.UserRepo.GetStats()
	classes, _ := h.UserRepo.GetAllClasses()

	return c.Render(http.StatusOK, "admin/statistics.html", map[string]interface{}{
		"Title":     "Statistik",
		"User":      user,
		"UserStats": userStats,
		"Classes":   classes,
		"StartDate": startDate,
		"EndDate":   endDate,
	})
}

// Admin User Management - Edit/Update/Delete

func (h *Handler) EditUser(c echo.Context) error {
	user := c.Get("user").(*models.User)

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=ID tidak valid")
	}

	targetUser, err := h.UserRepo.GetByID(userID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=User tidak ditemukan")
	}

	classes, _ := h.ClassRepo.GetAll()

	return c.Render(http.StatusOK, "admin/user_edit.html", map[string]interface{}{
		"Title":      "Edit Siswa",
		"User":       user,
		"TargetUser": targetUser,
		"Classes":    classes,
		"Error":      c.QueryParam("error"),
		"Success":    c.QueryParam("success"),
	})
}

func (h *Handler) UpdateUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=ID tidak valid")
	}

	targetUser, err := h.UserRepo.GetByID(userID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=User tidak ditemukan")
	}

	// Get form values
	fullName := c.FormValue("full_name")
	email := c.FormValue("email")
	class := c.FormValue("class")
	role := c.FormValue("role")

	// Validate email uniqueness if changed
	if email != targetUser.Email {
		existingUser, _ := h.UserRepo.GetByEmail(email)
		if existingUser != nil && existingUser.ID != userID {
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/users/edit/%d?error=Email sudah digunakan", userID))
		}
	}

	// Update user
	targetUser.FullName = fullName
	targetUser.Email = email
	targetUser.Class = class
	targetUser.Role = role

	if err := h.UserRepo.Update(targetUser); err != nil {
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/users/edit/%d?error=Gagal memperbarui user", userID))
	}

	// Handle password reset if provided
	newPassword := c.FormValue("new_password")
	if newPassword != "" {
		if len(newPassword) < 6 {
			return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/users/edit/%d?error=Password minimal 6 karakter", userID))
		}
		hashedPassword, _ := utils.HashPassword(newPassword)
		h.UserRepo.UpdatePassword(userID, hashedPassword)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/users?success=User berhasil diperbarui")
}

func (h *Handler) DeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=ID tidak valid")
	}

	// Prevent deleting yourself
	currentUser := c.Get("user").(*models.User)
	if userID == currentUser.ID {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=Tidak dapat menghapus akun sendiri")
	}

	if err := h.UserRepo.Delete(userID); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=Gagal menghapus user")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/users?success=User berhasil dihapus")
}

func (h *Handler) ShowUserDetail(c echo.Context) error {
	user := c.Get("user").(*models.User)

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=ID tidak valid")
	}

	targetUser, err := h.UserRepo.GetByID(userID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=User tidak ditemukan")
	}

	// Get user statistics
	today := time.Now().Format("2006-01-02")
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location()).Format("2006-01-02")

	prayerStats, _ := h.PrayerRepo.GetPrayerStats(userID, startOfMonth, today)
	fastingStats, _ := h.FastingRepo.GetFastingStats(userID, startOfMonth, today)
	totalReadings, _ := h.QuranRepo.GetTotalReadings(userID)
	totalPoints, _ := h.AmaliahRepo.GetTotalPoints(userID, startOfMonth, today)

	// Get recent activities
	recentPrayers, _ := h.PrayerRepo.GetByUser(userID, 7)
	recentFastings, _ := h.FastingRepo.GetByUser(userID, 7)
	recentQuran, _ := h.QuranRepo.GetByUser(userID, 5)
	recentAmaliah, _ := h.AmaliahRepo.GetByUser(userID, 10)

	return c.Render(http.StatusOK, "admin/user_detail.html", map[string]interface{}{
		"Title":          "Detail Siswa",
		"User":           user,
		"TargetUser":     targetUser,
		"PrayerStats":    prayerStats,
		"FastingStats":   fastingStats,
		"TotalReadings":  totalReadings,
		"TotalPoints":    totalPoints,
		"RecentPrayers":  recentPrayers,
		"RecentFastings": recentFastings,
		"RecentQuran":    recentQuran,
		"RecentAmaliah":  recentAmaliah,
	})
}

func (h *Handler) SearchUsers(c echo.Context) error {
	user := c.Get("user").(*models.User)

	query := c.QueryParam("q")
	users, _ := h.UserRepo.SearchUsers(query)

	return c.Render(http.StatusOK, "admin/users.html", map[string]interface{}{
		"Title":  "Kelola Siswa",
		"User":   user,
		"Users":  users,
		"Search": query,
	})
}

// Admin Reports

func (h *Handler) GenerateReport(c echo.Context) error {
	user := c.Get("user").(*models.User)

	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	reportType := c.QueryParam("type")

	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	var data map[string]interface{}

	switch reportType {
	case "daily":
		prayerStats, _ := h.PrayerRepo.GetTodayStats(startDate)
		fastingStats, _ := h.FastingRepo.GetTodayStats(startDate)
		quranStats, _ := h.QuranRepo.GetTodayStats(startDate)
		amaliahStats, _ := h.AmaliahRepo.GetTodayStats(startDate)
		data = map[string]interface{}{
			"prayer":  prayerStats,
			"fasting": fastingStats,
			"quran":   quranStats,
			"amaliah": amaliahStats,
		}
	default:
		data = map[string]interface{}{}
	}

	return c.Render(http.StatusOK, "admin/reports.html", map[string]interface{}{
		"Title":     "Laporan",
		"User":      user,
		"Data":      data,
		"StartDate": startDate,
		"EndDate":   endDate,
		"Type":      reportType,
	})
}

// Islamic Content Handlers - Doa
func (h *Handler) ShowDoa(c echo.Context) error {
	user := c.Get("user").(*models.User)

	source := c.QueryParam("source")
	search := c.QueryParam("search")

	var doaList []services.Doa
	var err error

	if search != "" {
		doaList, err = h.MuslimAPI.SearchDoa(search)
	} else {
		doaList, err = h.MuslimAPI.GetDoaBySource(source)
	}

	if err != nil {
		doaList = []services.Doa{}
	}

	return c.Render(http.StatusOK, "user/doa.html", map[string]interface{}{
		"Title":   "Doa Harian",
		"User":    user,
		"DoaList": doaList,
		"Source":  source,
		"Search":  search,
		"Error":   c.QueryParam("error"),
	})
}

// Islamic Content Handlers - Hadits
func (h *Handler) ShowHadits(c echo.Context) error {
	user := c.Get("user").(*models.User)

	search := c.QueryParam("search")
	nomorStr := c.QueryParam("nomor")

	var haditsList []services.Hadits
	var selectedHadits *services.Hadits
	var err error

	if nomorStr != "" {
		nomor, _ := strconv.Atoi(nomorStr)
		selectedHadits, err = h.MuslimAPI.GetHaditsByNumber(nomor)
		if err == nil && selectedHadits != nil {
			haditsList = []services.Hadits{*selectedHadits}
		}
	} else if search != "" {
		haditsList, err = h.MuslimAPI.SearchHadits(search)
	} else {
		haditsList, err = h.MuslimAPI.GetAllHadits()
	}

	if err != nil {
		haditsList = []services.Hadits{}
	}

	return c.Render(http.StatusOK, "user/hadits.html", map[string]interface{}{
		"Title":         "Hadits Arbain",
		"User":          user,
		"HaditsList":    haditsList,
		"Search":        search,
		"SelectedNomor": nomorStr,
		"Error":         c.QueryParam("error"),
	})
}

// Islamic Content Handlers - Quran Indonesia
func (h *Handler) ShowQuranIndonesia(c echo.Context) error {
	user := c.Get("user").(*models.User)

	surahIDStr := c.QueryParam("surah")

	var surahList []services.Surah
	var ayahList []services.Ayah
	var selectedSurah *services.Surah
	var err error

	surahList, err = h.MuslimAPI.GetAllSurah()
	if err != nil {
		surahList = []services.Surah{}
	}

	if surahIDStr != "" {
		surahID, _ := strconv.Atoi(surahIDStr)
		selectedSurah, _ = h.MuslimAPI.GetSurahByID(surahID)
		ayahList, err = h.MuslimAPI.GetAyahBySurah(surahID)
		if err != nil {
			ayahList = []services.Ayah{}
		}
	}

	return c.Render(http.StatusOK, "user/quran_indonesia.html", map[string]interface{}{
		"Title":         "Al-Quran Indonesia",
		"User":          user,
		"SurahList":     surahList,
		"AyahList":      ayahList,
		"SelectedSurah": selectedSurah,
		"SurahID":       surahIDStr,
		"Error":         c.QueryParam("error"),
	})
}

// Profile Handlers
func (h *Handler) ShowProfile(c echo.Context) error {
	user := c.Get("user").(*models.User)

	totalReadings, _ := h.QuranRepo.GetTotalReadings(user.ID)

	today := time.Now()
	startOfMonth := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location()).Format("2006-01-02")
	prayerStats, _ := h.PrayerRepo.GetPrayerStats(user.ID, startOfMonth, today.Format("2006-01-02"))

	fastingStats, _ := h.FastingRepo.GetFastingStats(user.ID, startOfMonth, today.Format("2006-01-02"))

	provinsiList, _ := h.ImsakiyahService.GetProvinsi()

	var kabkotaList []string
	if user.Provinsi != "" {
		kabkotaList, _ = h.ImsakiyahService.GetKabkota(user.Provinsi)
	}

	return c.Render(http.StatusOK, "user/profile.html", map[string]interface{}{
		"Title":         "Profil Saya",
		"User":          user,
		"TotalReadings": totalReadings,
		"PrayerStats":   prayerStats,
		"FastingStats":  fastingStats,
		"ProvinsiList":  provinsiList,
		"KabkotaList":   kabkotaList,
		"Error":         c.QueryParam("error"),
		"Success":       c.QueryParam("success"),
	})
}

func (h *Handler) UpdateProfile(c echo.Context) error {
	user := c.Get("user").(*models.User)

	req := &models.ProfileUpdateRequest{
		FullName:     c.FormValue("full_name"),
		Email:        c.FormValue("email"),
		Class:        c.FormValue("class"),
		Bio:          c.FormValue("bio"),
		Avatar:       c.FormValue("avatar"),
		Theme:        c.FormValue("theme"),
		TargetKhatam: 30,
		Provinsi:     c.FormValue("provinsi"),
		Kabkota:      c.FormValue("kabkota"),
	}

	targetKhatam, err := strconv.Atoi(c.FormValue("target_khatam"))
	if err == nil && targetKhatam >= 1 && targetKhatam <= 30 {
		req.TargetKhatam = targetKhatam
	}

	if req.Email != user.Email {
		existingUser, _ := h.UserRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != user.ID {
			return c.Redirect(http.StatusSeeOther, "/user/profile?error=Email sudah digunakan")
		}
	}

	err = h.UserRepo.UpdateProfile(user.ID, req)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal memperbarui profil")
	}

	updatedUser, _ := h.UserRepo.GetByID(user.ID)
	c.Set("user", updatedUser)

	return c.Redirect(http.StatusSeeOther, "/user/profile?success=Profil berhasil diperbarui")
}

func (h *Handler) UpdateAvatar(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Source
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal mengupload avatar")
	}
	src, err := file.Open()
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal membuka file")
	}
	defer src.Close()

	// Destination
	// Ensure directory exists
	uploadDir := "web/static/uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal membuat direktori upload")
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%d%s", user.ID, time.Now().Unix(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal menyimpan file")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal menyalin file")
	}

	// Update User Avatar in DB
	// We need to update the user struct and save it.
	// Assuming UserRepo has Update method or we handle it here.
	// Since UserRepo.Update updates everything, let's just update the field manually or use a specific method.
	// For simplicity, let's update the specific field via DB exec or UserRepo update.
	// Let's check UserRepo.Update implementation later or just do a direct exec for now to be safe/fast.
	_, err = h.DB.Exec("UPDATE users SET avatar = ? WHERE id = ?", filename, user.ID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal memperbarui database")
	}

	return c.Redirect(http.StatusSeeOther, "/user/profile?success=Avatar berhasil diperbarui")
}

func (h *Handler) ChangePassword(c echo.Context) error {
	user := c.Get("user").(*models.User)

	currentPassword := c.FormValue("current_password")
	newPassword := c.FormValue("new_password")
	confirmPassword := c.FormValue("confirm_password")

	// Validate current password
	if !utils.CheckPassword(currentPassword, user.PasswordHash) {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Password saat ini salah")
	}

	// Validate new password
	if len(newPassword) < 6 {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Password baru minimal 6 karakter")
	}

	// Validate password match
	if newPassword != confirmPassword {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Password baru dan konfirmasi tidak cocok")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal mengubah password")
	}

	// Update password
	err = h.UserRepo.UpdatePassword(user.ID, hashedPassword)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/profile?error=Gagal mengubah password")
	}

	return c.Redirect(http.StatusSeeOther, "/user/profile?success=Password berhasil diubah")
}

func (h *Handler) GetKabkotaAPI(c echo.Context) error {
	provinsi := c.QueryParam("provinsi")
	if provinsi == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Provinsi is required"})
	}

	kabkotaList, err := h.ImsakiyahService.GetKabkota(provinsi)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"provinsi": provinsi,
		"kabkota":  kabkotaList,
	})
}

func (h *Handler) GetImsakiyahAPI(c echo.Context) error {
	user := c.Get("user").(*models.User)

	provinsi := c.QueryParam("provinsi")
	kabkota := c.QueryParam("kabkota")

	if provinsi == "" || kabkota == "" {
		if user.Provinsi != "" && user.Kabkota != "" {
			provinsi = user.Provinsi
			kabkota = user.Kabkota
		} else {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Provinsi dan Kabkota diperlukan"})
		}
	}

	data, err := h.ImsakiyahService.GetImsakiyah(provinsi, kabkota)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, data)
}

func (h *Handler) AutoDetectLocation(c echo.Context) error {
	latStr := c.FormValue("lat")
	longStr := c.FormValue("long")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid latitude"})
	}
	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid longitude"})
	}

	// 1. Reverse Geocode
	detectedProv, detectedCity, err := h.ShalatService.ReverseGeocode(lat, long)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mendeteksi lokasi: " + err.Error()})
	}

	// 2. Match with API Data
	matchedProv, matchedCity, err := h.ShalatService.MatchLocation(detectedProv, detectedCity)
	
	// Prepare response
	response := map[string]string{
		"detected_prov": detectedProv,
		"detected_city": detectedCity,
		"provinsi":     matchedProv,
		"kabkota":      matchedCity,
	}

	if err != nil {
		response["status"] = "partial_match"
		response["message"] = "Lokasi terdeteksi namun tidak cocok sempurna dengan database. Silakan sesuaikan manual."
	} else {
		response["status"] = "success"
	}

	return c.JSON(http.StatusOK, response)
}

// Middleware
func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		token, err := utils.ValidateToken(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		userID := int(claims["user_id"].(float64))
		user, err := h.UserRepo.GetByID(userID)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		c.Set("user", user)
		return next(c)
	}
}

func (h *Handler) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("token")
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		token, err := utils.ValidateToken(cookie.Value)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		userID := int(claims["user_id"].(float64))
		user, err := h.UserRepo.GetByID(userID)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if user.Role != "superadmin" {
			return c.Redirect(http.StatusSeeOther, "/user/dashboard")
		}

		c.Set("user", user)
		return next(c)
	}
}

// Error Handlers
func (h *Handler) NotFound(c echo.Context) error {
	return c.Render(http.StatusNotFound, "errors/404.html", map[string]interface{}{
		"Title": "Halaman Tidak Ditemukan",
	})
}

func (h *Handler) Forbidden(c echo.Context) error {
	return c.Render(http.StatusForbidden, "errors/403.html", map[string]interface{}{
		"Title": "Akses Ditolak",
	})
}

func (h *Handler) DownloadCertificate(c echo.Context) error {
	user := c.Get("user").(*models.User)

	// Gather Stats
	totalPoints := user.Points
	totalPages, _ := h.QuranRepo.GetTotalPagesRead(user.ID)
	
	progressPercent := float64(totalPages) / 604.0 * 100
	
	stats := map[string]interface{}{
		"total_points": totalPoints,
		"total_pages":  totalPages,
		"khatam_percent": int(progressPercent),
	}

	pdfBytes, err := h.CertificateService.Generate(user, stats)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to generate certificate: "+err.Error())
	}

	// Serve content
	// Format: Sertifikat Amaliah Ramadhan_NAMASISWA
	filename := fmt.Sprintf("Sertifikat Amaliah Ramadhan_%s.pdf", user.FullName)
	c.Response().Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	return c.Blob(http.StatusOK, "application/pdf", pdfBytes)
}
