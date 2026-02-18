package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/services"
)

func (h *Handler) ImportUsers(c echo.Context) error {
	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "File not found"})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	ext := strings.ToLower(filepath.Ext(file.Filename))
	
	var result *services.ImportResult

	if ext == ".csv" {
		result, err = h.AdminService.ImportStudentsFromCSV(src)
	} else if ext == ".xlsx" {
		result, err = h.AdminService.ImportStudentsFromExcel(src)
	} else {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error=Format file tidak valid. Harap upload .csv atau .xlsx")
	}

	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/users?error="+err.Error())
	}

	message := fmt.Sprintf("Berhasil import %d siswa", result.Success)
	if result.Failed > 0 {
		message += fmt.Sprintf(", %d gagal", result.Failed)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/users?success="+message)
}

func (h *Handler) ManageClasses(c echo.Context) error {
	classes, err := h.ClassRepo.GetAll()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "errors/500.html", nil)
	}

	return c.Render(http.StatusOK, "admin/classes.html", map[string]interface{}{
		"Title":   "Manajemen Kelas",
		"Classes": classes,
		"Success": c.QueryParam("success"),
		"Error":   c.QueryParam("error"),
	})
}

func (h *Handler) CreateClass(c echo.Context) error {
	name := c.FormValue("name")
	level := c.FormValue("level")
	description := c.FormValue("description")

	err := h.ClassRepo.Create(&models.Class{
		Name:        name,
		Level:       level,
		Description: description,
	})
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/classes?error=Gagal membuat kelas: "+err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/admin/classes?success=Kelas berhasil dibuat")
}

func (h *Handler) UpdateClass(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	name := c.FormValue("name")
	level := c.FormValue("level")
	description := c.FormValue("description")

	err := h.ClassRepo.Update(&models.Class{
		ID:          id,
		Name:        name,
		Level:       level,
		Description: description,
	})
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/classes?error=Gagal update kelas: "+err.Error())
	}

	return c.Redirect(http.StatusSeeOther, "/admin/classes?success=Kelas berhasil diupdate")
}

func (h *Handler) DeleteClass(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.ClassRepo.Delete(id)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/classes?error=Gagal menghapus kelas")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/classes?success=Kelas berhasil dihapus")
}
