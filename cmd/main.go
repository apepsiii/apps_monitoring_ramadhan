package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ramadhan/amaliah-monitoring/internal/config"
	"github.com/ramadhan/amaliah-monitoring/internal/handlers"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Echo
	e := echo.New()

	// Setup Template Renderer
	e.Renderer = config.NewTemplateRenderer()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Static("web/static"))

	// Initialize Database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Run Migrations
	if err := config.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize Handlers
	h := handlers.NewHandler(db)

	// Routes
	e.GET("/", h.Home)

	// Public Routes - Jadwal Shalat & Imsakiyah
	e.GET("/jadwal", h.ShowJadwal)
	e.GET("/api/kabkota", h.GetKabkotaAPI)

	// Auth Routes
	e.GET("/login", h.ShowLogin)
	e.POST("/login", h.Login)
	e.GET("/register", h.ShowRegister)
	e.POST("/register", h.Register)
	e.GET("/logout", h.Logout)

	// Protected Routes Group
	user := e.Group("/user")
	user.Use(h.AuthMiddleware)
	user.GET("/dashboard", h.UserDashboard)
	user.GET("/prayers", h.ShowPrayers)
	user.POST("/prayers", h.SavePrayers)
	user.GET("/fasting", h.ShowFasting)
	user.POST("/fasting", h.SaveFasting)
	user.GET("/quran", h.ShowQuran)
	user.POST("/quran", h.SaveQuran)
	user.GET("/quran/delete/:id", h.DeleteQuran)
	user.GET("/amaliah", h.ShowAmaliah)
	user.POST("/amaliah", h.SaveAmaliah)

	// Islamic Content Routes
	user.GET("/doa", h.ShowDoa)
	user.GET("/hadits", h.ShowHadits)
	user.GET("/quran-indonesia", h.ShowQuranIndonesia)

	// Profile Routes
	user.GET("/profile", h.ShowProfile)
	user.POST("/profile", h.UpdateProfile)
	user.POST("/profile/change-password", h.ChangePassword)

	// API Routes (protected)
	user.GET("/api/kabkota", h.GetKabkotaAPI)
	user.GET("/api/imsakiyah", h.GetImsakiyahAPI)

	// Admin Routes Group
	admin := e.Group("/admin")
	admin.Use(h.AdminMiddleware)
	admin.GET("/dashboard", h.AdminDashboard)
	admin.GET("/users", h.ManageUsers)
	admin.POST("/users", h.CreateUser)
	admin.GET("/users/search", h.SearchUsers)
	admin.GET("/users/edit/:id", h.EditUser)
	admin.POST("/users/update/:id", h.UpdateUser)
	admin.GET("/users/delete/:id", h.DeleteUser)
	admin.GET("/users/detail/:id", h.ShowUserDetail)
	admin.GET("/reports", h.ShowReports)
	admin.GET("/reports/generate", h.GenerateReport)
	admin.GET("/statistics", h.ShowStatistics)

	// Start Server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
