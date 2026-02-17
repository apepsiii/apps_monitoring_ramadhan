package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) DownloadReport(c echo.Context) error {
	reportType := c.QueryParam("type") // "daily" or "student"
	date := c.QueryParam("date")
	className := c.QueryParam("class")
	
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	if reportType == "daily" {
		f, err := h.ExportService.GenerateDailyReportExcel(date, className)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/reports?error="+err.Error())
		}
		
		c.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=Laporan_Harian_%s.xlsx", date))
		
		return f.Write(c.Response().Writer)
	}

	return c.Redirect(http.StatusSeeOther, "/admin/reports?error=Tipe laporan tidak valid")
}
