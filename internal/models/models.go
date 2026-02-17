package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Class        string    `json:"class"`
	Role         string    `json:"role"`
	Points       int       `json:"points"`
	Avatar       string    `json:"avatar"`
	Bio          string    `json:"bio"`
	Theme        string    `json:"theme"`
	TargetKhatam int       `json:"target_khatam"`
	Provinsi     string    `json:"provinsi"`
	Kabkota      string    `json:"kabkota"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ProfileUpdateRequest struct {
	FullName     string `json:"full_name" form:"full_name"`
	Email        string `json:"email" form:"email"`
	Class        string `json:"class" form:"class"`
	Bio          string `json:"bio" form:"bio"`
	Avatar       string `json:"avatar" form:"avatar"`
	Theme        string `json:"theme" form:"theme"`
	TargetKhatam int    `json:"target_khatam" form:"target_khatam"`
	Provinsi     string `json:"provinsi" form:"provinsi"`
	Kabkota      string `json:"kabkota" form:"kabkota"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" form:"current_password"`
	NewPassword     string `json:"new_password" form:"new_password"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
}

type Prayer struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Date      string    `json:"date"`
	Subuh     string    `json:"subuh"`
	Dzuhur    string    `json:"dzuhur"`
	Ashar     string    `json:"ashar"`
	Maghrib   string    `json:"maghrib"`
	Isya      string    `json:"isya"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Fasting struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Date      string    `json:"date"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type QuranReading struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Date           string    `json:"date"`
	StartSurahID   int       `json:"start_surah_id"`
	StartSurahName string    `json:"start_surah_name"`
	StartAyah      int       `json:"start_ayah"`
	EndSurahID     int       `json:"end_surah_id"`
	EndSurahName   string    `json:"end_surah_name"`
	EndAyah        int       `json:"end_ayah"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
}

type AmaliahType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Points      int       `json:"points"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type DailyAmaliah struct {
	ID            int         `json:"id"`
	UserID        int         `json:"user_id"`
	AmaliahTypeID int         `json:"amaliah_type_id"`
	Date          string      `json:"date"`
	Notes         string      `json:"notes"`
	AmaliahType   AmaliahType `json:"amaliah_type,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	FullName string `json:"full_name" form:"full_name"`
	Class    string `json:"class" form:"class"`
}

type Streak struct {
	UserID        int `json:"user_id"`
	PrayerStreak  int `json:"prayer_streak"`
	FastingStreak int `json:"fasting_streak"`
	QuranStreak   int `json:"quran_streak"`
	BestPrayer    int `json:"best_prayer"`
	BestFasting   int `json:"best_fasting"`
	BestQuran     int `json:"best_quran"`
	AmaliahStreak int `json:"amaliah_streak"`
	BestAmaliah   int `json:"best_amaliah"`
}

type ImsakiyahSchedule struct {
	Tanggal int    `json:"tanggal"`
	Imsak   string `json:"imsak"`
	Subuh   string `json:"subuh"`
	Terbit  string `json:"terbit"`
	Dhuha   string `json:"dhuha"`
	Dzuhur  string `json:"dzuhur"`
	Ashar   string `json:"ashar"`
	Maghrib string `json:"maghrib"`
	Isya    string `json:"isya"`
}

type ImsakiyahData struct {
	Provinsi  string              `json:"provinsi"`
	Kabkota   string              `json:"kabkota"`
	Hijriah   string              `json:"hijriah"`
	Masehi    string              `json:"masehi"`
	Imsakiyah []ImsakiyahSchedule `json:"imsakiyah"`
}

type ImsakiyahResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    ImsakiyahData `json:"data"`
}

type ShalatSchedule struct {
	Tanggal        int    `json:"tanggal"`
	TanggalLengkap string `json:"tanggal_lengkap"`
	Hari           string `json:"hari"`
	Imsak          string `json:"imsak"`
	Subuh          string `json:"subuh"`
	Terbit         string `json:"terbit"`
	Dhuha          string `json:"dhuha"`
	Dzuhur         string `json:"dzuhur"`
	Ashar          string `json:"ashar"`
	Maghrib        string `json:"maghrib"`
	Isya           string `json:"isya"`
}

type ShalatData struct {
	Provinsi  string           `json:"provinsi"`
	Kabkota   string           `json:"kabkota"`
	Bulan     int              `json:"bulan"`
	Tahun     int              `json:"tahun"`
	BulanNama string           `json:"bulan_nama"`
	Jadwal    []ShalatSchedule `json:"jadwal"`
}

type ShalatResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    ShalatData `json:"data"`
}

type ProvinsiResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type KabkotaResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}
