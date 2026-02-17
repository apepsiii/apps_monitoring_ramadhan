package services

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/ramadhan/amaliah-monitoring/internal/models"
	"github.com/ramadhan/amaliah-monitoring/internal/repository"
	"github.com/ramadhan/amaliah-monitoring/internal/utils"
	"github.com/xuri/excelize/v2"
)

type AdminService struct {
	UserRepo *repository.UserRepository
}

func NewAdminService(userRepo *repository.UserRepository) *AdminService {
	return &AdminService{
		UserRepo: userRepo,
	}
}

type ImportResult struct {
	Total   int      `json:"total"`
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

func (s *AdminService) ImportStudentsFromCSV(file multipart.File) (*ImportResult, error) {
	reader := csv.NewReader(file)
	result := &ImportResult{}

	// Skip header
	_, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	lineNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Error reading record", lineNum))
			lineNum++
			continue
		}
		lineNum++

		// Expected format: username, email, full_name, class, password
		if len(record) < 5 {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Insufficient columns", lineNum))
			continue
		}

		username := record[0]
		email := record[1]
		fullName := record[2]
		class := record[3]
		password := record[4]

		// Basic validation
		if username == "" || email == "" || fullName == "" || password == "" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Missing required fields", lineNum))
			continue
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Failed to hash password", lineNum))
			continue
		}

		user := &models.User{
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
			FullName:     fullName,
			Class:        class,
			Role:         "user",
			Points:       0,
		}

		err = s.UserRepo.Create(user)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Line %d: Failed to create user (%s) - possibly duplicate", lineNum, username))
			continue
		}

		result.Success++
		result.Total++
	}

	return result, nil
}

func (s *AdminService) ImportStudentsFromExcel(file multipart.File) (*ImportResult, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get first sheet name
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	result := &ImportResult{}

	// Skip header (start from index 1)
	for i, row := range rows {
		if i == 0 {
			continue
		}

		// Expected format: username, email, full_name, class, password
		if len(row) < 5 {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Insufficient columns", i+1))
			continue
		}

		username := row[0]
		email := row[1]
		fullName := row[2]
		class := row[3]
		password := row[4]

		// Basic validation
		if username == "" || email == "" || fullName == "" || password == "" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Missing required fields", i+1))
			continue
		}

		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to hash password", i+1))
			continue
		}

		user := &models.User{
			Username:     username,
			Email:        email,
			PasswordHash: hashedPassword,
			FullName:     fullName,
			Class:        class,
			Role:         "user",
			Points:       0,
		}

		err = s.UserRepo.Create(user)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("Row %d: Failed to create user (%s) - possibly duplicate", i+1, username))
			continue
		}

		result.Success++
		result.Total++
	}

	return result, nil
}
