package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/utils"
)

// ─── Public Admin Registration ────────────────────────────────────────────────

type AdminRegisterFormData struct {
	FullName      string
	Phone         string
	SchoolName    string
	SchoolAddress string
	SchoolLevel   string
	StudentCount  string
	Username      string
	Email         string
}

// ShowAdminRegister renders the public admin registration form
func (h *Handler) ShowAdminRegister(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/register_admin.html", map[string]interface{}{
		"Title":    "Daftar Admin Sekolah",
		"FormData": AdminRegisterFormData{},
	})
}

// AdminRegister handles the public admin registration form submission
func (h *Handler) AdminRegister(c echo.Context) error {
	formData := AdminRegisterFormData{
		FullName:      c.FormValue("full_name"),
		Phone:         c.FormValue("phone"),
		SchoolName:    c.FormValue("school_name"),
		SchoolAddress: c.FormValue("school_address"),
		SchoolLevel:   c.FormValue("school_level"),
		StudentCount:  c.FormValue("student_count"),
		Username:      c.FormValue("username"),
		Email:         c.FormValue("email"),
	}
	password := c.FormValue("password")

	renderErr := func(msg string) error {
		return c.Render(http.StatusOK, "auth/register_admin.html", map[string]interface{}{
			"Title":    "Daftar Admin Sekolah",
			"Error":    msg,
			"FormData": formData,
		})
	}

	// Validate required fields
	if formData.FullName == "" || formData.Phone == "" || formData.SchoolName == "" ||
		formData.SchoolLevel == "" || formData.Username == "" || formData.Email == "" || password == "" {
		return renderErr("Semua field bertanda * wajib diisi")
	}
	if len(password) < 6 {
		return renderErr("Password minimal 6 karakter")
	}

	// Check username uniqueness in users table
	var uCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", formData.Username).Scan(&uCount)
	if uCount > 0 {
		return renderErr("Username sudah digunakan, pilih username lain")
	}

	// Check username uniqueness in pending requests
	var rCount int
	h.DB.QueryRow("SELECT COUNT(*) FROM admin_requests WHERE username = ? AND status = 'pending'", formData.Username).Scan(&rCount)
	if rCount > 0 {
		return renderErr("Username sudah ada dalam daftar pengajuan yang sedang diproses")
	}

	// Hash password
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return renderErr("Terjadi kesalahan, coba lagi")
	}

	studentCount := 0
	if formData.StudentCount != "" {
		studentCount, _ = strconv.Atoi(formData.StudentCount)
	}

	// Insert into admin_requests
	_, err = h.DB.Exec(`
		INSERT INTO admin_requests (full_name, phone, school_name, school_address, school_level, student_count, username, email, password_hash, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending')`,
		formData.FullName, formData.Phone, formData.SchoolName, formData.SchoolAddress,
		formData.SchoolLevel, studentCount, formData.Username, formData.Email, hashed,
	)
	if err != nil {
		return renderErr("Gagal menyimpan pendaftaran, coba lagi")
	}

	return c.Redirect(http.StatusSeeOther, "/register-admin/thanks")
}

// AdminRegisterThanks shows the thank you page after registration
func (h *Handler) AdminRegisterThanks(c echo.Context) error {
	return c.Render(http.StatusOK, "auth/register_admin_thanks.html", map[string]interface{}{
		"Title": "Pendaftaran Berhasil",
	})
}

// ─── School Approve (Superadmin) ─────────────────────────────────────────────

// SchoolApprove approves a pending admin_request: creates user + school, activates account
func (h *Handler) SchoolApprove(c echo.Context) error {
	reqID, _ := strconv.Atoi(c.Param("id"))
	if reqID <= 0 {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=ID tidak valid")
	}

	// Fetch the request
	var req struct {
		FullName      string
		Phone         string
		SchoolName    string
		SchoolAddress string
		SchoolLevel   string
		StudentCount  int
		Username      string
		Email         string
		PasswordHash  string
	}
	err := h.DB.QueryRow(`
		SELECT full_name, phone, school_name, school_address, school_level, student_count, username, email, password_hash
		FROM admin_requests WHERE id = ? AND status = 'pending'`, reqID,
	).Scan(&req.FullName, &req.Phone, &req.SchoolName, &req.SchoolAddress, &req.SchoolLevel,
		&req.StudentCount, &req.Username, &req.Email, &req.PasswordHash)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=Pengajuan tidak ditemukan atau sudah diproses")
	}

	// Generate school code
	schoolCode := utils.GenerateRandomString(8)

	// Create school
	schoolRes, err := h.DB.Exec(
		"INSERT INTO schools (name, code, address, status) VALUES (?, ?, ?, 'active')",
		req.SchoolName, schoolCode, req.SchoolAddress,
	)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=Gagal membuat sekolah")
	}
	schoolID, _ := schoolRes.LastInsertId()

	// Create user with role = admin
	userRes, err := h.DB.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, class, role, points, avatar, bio, theme, school_id)
		VALUES (?, ?, ?, ?, '', 'admin', 0, 'default', '', 'emerald', ?)`,
		req.Username, req.Email, req.PasswordHash, req.FullName, schoolID,
	)
	if err != nil {
		// Rollback school creation
		h.DB.Exec("DELETE FROM schools WHERE id = ?", schoolID)
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=Gagal membuat akun admin (username/email mungkin sudah ada)")
	}
	userID, _ := userRes.LastInsertId()

	// Set admin_id on school
	h.DB.Exec("UPDATE schools SET admin_id = ? WHERE id = ?", userID, schoolID)

	// Mark request as approved
	h.DB.Exec("UPDATE admin_requests SET status = 'approved' WHERE id = ?", reqID)

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard?success=Akun admin berhasil diaktifkan untuk "+req.SchoolName)
}

// ─── School Reject (Superadmin) ──────────────────────────────────────────────

// SchoolReject rejects a pending admin_request
func (h *Handler) SchoolReject(c echo.Context) error {
	reqID, _ := strconv.Atoi(c.Param("id"))
	if reqID <= 0 {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=ID tidak valid")
	}

	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM admin_requests WHERE id = ? AND status = 'pending'", reqID).Scan(&count)
	if count == 0 {
		return c.Redirect(http.StatusSeeOther, "/admin/dashboard?error=Pengajuan tidak ditemukan atau sudah diproses")
	}

	h.DB.Exec("UPDATE admin_requests SET status = 'rejected' WHERE id = ?", reqID)

	return c.Redirect(http.StatusSeeOther, "/admin/dashboard?success=Pengajuan telah ditolak")
}

// ─── School Admin Dashboard ───────────────────────────────────────────────────

// SchoolAdminDashboard shows the school management page for school admins
func (h *Handler) SchoolAdminDashboard(c echo.Context) error {
	user := c.Get("user").(*models.User)

	if user.Role != "admin" {
		return c.Redirect(http.StatusSeeOther, "/user/dashboard")
	}

	var school models.School
	err := h.DB.QueryRow(
		"SELECT id, name, code, address, admin_id FROM schools WHERE id = ?",
		user.SchoolID,
	).Scan(&school.ID, &school.Name, &school.Code, &school.Address, &school.AdminID)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/user/dashboard?error=Sekolah tidak ditemukan")
	}

	rows, err := h.DB.Query(
		"SELECT id, full_name, class, points, avatar, role FROM users WHERE school_id = ? ORDER BY role ASC, full_name ASC",
		user.SchoolID,
	)
	if err != nil {
		return c.Render(http.StatusOK, "school/admin_dashboard.html", map[string]interface{}{
			"Title":  "Kelola Sekolah",
			"Error":  "Gagal mengambil data anggota",
			"School": school,
			"User":   user,
		})
	}
	defer rows.Close()

	var members []models.User
	for rows.Next() {
		var m models.User
		if err := rows.Scan(&m.ID, &m.FullName, &m.Class, &m.Points, &m.Avatar, &m.Role); err != nil {
			continue
		}
		members = append(members, m)
	}

	return c.Render(http.StatusOK, "school/admin_dashboard.html", map[string]interface{}{
		"Title":   "Kelola Sekolah",
		"School":  school,
		"Members": members,
		"User":    user,
		"Success": c.QueryParam("success"),
		"Error":   c.QueryParam("error"),
	})
}

// ─── School Update ────────────────────────────────────────────────────────────

func (h *Handler) SchoolUpdate(c echo.Context) error {
	user := c.Get("user").(*models.User)
	if user.Role != "admin" {
		return c.Redirect(http.StatusSeeOther, "/user/dashboard")
	}
	newName := c.FormValue("name")
	if newName == "" {
		return c.Redirect(http.StatusSeeOther, "/school/admin?error=Nama sekolah tidak boleh kosong")
	}
	h.DB.Exec("UPDATE schools SET name = ? WHERE id = ?", newName, user.SchoolID)
	return c.Redirect(http.StatusSeeOther, "/school/admin?success=Nama sekolah berhasil diperbarui")
}

// ─── School Remove Member ─────────────────────────────────────────────────────

func (h *Handler) SchoolRemoveMember(c echo.Context) error {
	user := c.Get("user").(*models.User)
	if user.Role != "admin" {
		return c.Redirect(http.StatusSeeOther, "/user/dashboard")
	}
	memberID, _ := strconv.Atoi(c.Param("id"))
	if memberID <= 0 || memberID == user.ID {
		return c.Redirect(http.StatusSeeOther, "/school/admin?error=Tidak dapat mengeluarkan anggota ini")
	}
	var count int
	h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND school_id = ?", memberID, user.SchoolID).Scan(&count)
	if count == 0 {
		return c.Redirect(http.StatusSeeOther, "/school/admin?error=Anggota tidak ditemukan di sekolah ini")
	}
	h.DB.Exec("UPDATE users SET school_id = 0 WHERE id = ? AND role != 'admin'", memberID)
	return c.Redirect(http.StatusSeeOther, "/school/admin?success=Anggota berhasil dikeluarkan")
}

// ─── Unused stubs kept for route compatibility ────────────────────────────────

func (h *Handler) SchoolSetup(c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/register-admin")
}

// suppress unused import warnings
var _ = time.Now
var _ = os.Getenv
