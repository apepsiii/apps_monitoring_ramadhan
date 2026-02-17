package installer

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	AppName      = "amaliah-ramadhan"
	ServiceName  = "amaliah-ramadhan"
	InstallDir   = "/opt/amaliah-ramadhan"
	ServiceFile  = "/etc/systemd/system/amaliah-ramadhan.service"
	DefaultPort  = "8080"
	ConfigFile   = "/opt/amaliah-ramadhan/.env"
	DatabaseFile = "/opt/amaliah-ramadhan/amaliah.db"
)

type Installer struct {
	Port         string
	IsNewInstall bool
	Reader       *bufio.Reader
}

func NewInstaller() *Installer {
	return &Installer{
		Port:   DefaultPort,
		Reader: bufio.NewReader(os.Stdin),
	}
}

// ShowBanner menampilkan banner aplikasi
func (i *Installer) ShowBanner() {
	banner := `
╔══════════════════════════════════════════════════════════╗
║                                                          ║
║        🌙 AMALIAH RAMADHAN - SMK NIBA INSTALLER 🌙      ║
║                                                          ║
║              Monitoring Ibadah Harian Ramadhan          ║
║                   SMK NIBA Isep Misbah                  ║
║                                                          ║
╚══════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
	fmt.Printf("Version: 1.0.0\n")
	fmt.Printf("Platform: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
}

// ShowMenu menampilkan menu utama
func (i *Installer) ShowMenu() int {
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    PILIH AKSI                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  1. 🆕 Install Baru")
	fmt.Println("  2. 🔄 Update/Upgrade")
	fmt.Println("  3. 🗑️  Uninstall")
	fmt.Println("  4. ℹ️  Informasi Status")
	fmt.Println("  5. ❌ Keluar")
	fmt.Println()
	fmt.Print("Pilih (1-5): ")

	choice, _ := i.Reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	num, err := strconv.Atoi(choice)
	if err != nil {
		return 0
	}
	return num
}

// AskPort meminta input port
func (i *Installer) AskPort() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  KONFIGURASI PORT                        ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("Port default: %s\n", DefaultPort)
	fmt.Print("Gunakan port lain? (tekan Enter untuk default, atau ketik port): ")

	port, _ := i.Reader.ReadString('\n')
	port = strings.TrimSpace(port)

	if port != "" {
		i.Port = port
	}

	fmt.Printf("✓ Port dipilih: %s\n", i.Port)
}

// Confirm meminta konfirmasi
func (i *Installer) Confirm(message string) bool {
	fmt.Println()
	fmt.Printf("%s (y/n): ", message)
	response, _ := i.Reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// ShowProgress menampilkan progress dengan animasi
func (i *Installer) ShowProgress(message string, task func() error) error {
	fmt.Printf("\n⏳ %s", message)

	done := make(chan bool)
	errChan := make(chan error)

	// Animasi loading
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r⏳ %s %s", message, spinner[i%len(spinner)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// Jalankan task
	go func() {
		errChan <- task()
	}()

	err := <-errChan
	done <- true

	if err != nil {
		fmt.Printf("\r❌ %s [GAGAL]\n", message)
		return err
	}

	fmt.Printf("\r✅ %s [SUKSES]\n", message)
	return nil
}

// CheckPermission mengecek apakah running sebagai root
func (i *Installer) CheckPermission() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("installer harus dijalankan sebagai root (gunakan sudo)")
	}
	return nil
}

// CheckExistingInstallation mengecek apakah sudah ada instalasi
func (i *Installer) CheckExistingInstallation() bool {
	if _, err := os.Stat(InstallDir); !os.IsNotExist(err) {
		return true
	}
	return false
}

// InstallNew melakukan instalasi baru
func (i *Installer) InstallNew() error {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  INSTALASI BARU                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	// Check existing installation
	if i.CheckExistingInstallation() {
		fmt.Println()
		fmt.Println("⚠️  Instalasi sudah ada!")
		if !i.Confirm("Hapus instalasi lama dan install ulang?") {
			return fmt.Errorf("instalasi dibatalkan")
		}
		if err := i.Uninstall(); err != nil {
			return err
		}
	}

	// Ask for port
	i.AskPort()

	// Confirmation
	if !i.Confirm("Lanjutkan instalasi?") {
		return fmt.Errorf("instalasi dibatalkan oleh user")
	}

	fmt.Println()
	fmt.Println("🚀 Memulai instalasi...")
	fmt.Println()

	// Step 1: Create directory
	if err := i.ShowProgress("Membuat direktori instalasi", func() error {
		return os.MkdirAll(InstallDir, 0755)
	}); err != nil {
		return err
	}

	// Step 2: Copy binary
	if err := i.ShowProgress("Menyalin aplikasi ke sistem", func() error {
		return i.copyBinary()
	}); err != nil {
		return err
	}

	// Step 3: Copy static files
	if err := i.ShowProgress("Menyalin file statis (templates, css, js)", func() error {
		return i.copyStaticFiles()
	}); err != nil {
		return err
	}

	// Step 4: Create config
	if err := i.ShowProgress("Membuat file konfigurasi", func() error {
		return i.createConfig()
	}); err != nil {
		return err
	}

	// Step 5: Setup database
	if err := i.ShowProgress("Inisialisasi database", func() error {
		return i.setupDatabase()
	}); err != nil {
		return err
	}

	// Step 6: Create systemd service
	if err := i.ShowProgress("Membuat systemd service", func() error {
		return i.createService()
	}); err != nil {
		return err
	}

	// Step 7: Set permissions
	if err := i.ShowProgress("Mengatur hak akses file", func() error {
		return i.setPermissions()
	}); err != nil {
		return err
	}

	// Step 8: Enable and start service
	if err := i.ShowProgress("Mengaktifkan dan memulai service", func() error {
		return i.enableAndStartService()
	}); err != nil {
		return err
	}

	// Success message
	i.showSuccessMessage()

	return nil
}

// copyBinary menyalin binary ke install directory
func (i *Installer) copyBinary() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	destPath := filepath.Join(InstallDir, AppName)

	input, err := os.ReadFile(execPath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destPath, input, 0755)
	if err != nil {
		return err
	}

	return nil
}

// copyStaticFiles menyalin file statis
func (i *Installer) copyStaticFiles() error {
	// Files will be embedded in binary, so we extract them
	staticDir := filepath.Join(InstallDir, "web")

	// Create web directory structure
	dirs := []string{
		filepath.Join(staticDir, "static"),
		filepath.Join(staticDir, "static/css"),
		filepath.Join(staticDir, "static/js"),
		filepath.Join(staticDir, "static/images"),
		filepath.Join(staticDir, "static/fonts"),
		filepath.Join(staticDir, "templates"),
		filepath.Join(staticDir, "templates/layouts"),
		filepath.Join(staticDir, "templates/auth"),
		filepath.Join(staticDir, "templates/user"),
		filepath.Join(staticDir, "templates/admin"),
		filepath.Join(staticDir, "templates/errors"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Static files will be extracted from embedded data
	// For now, copy from current directory if exists
	srcWeb := "web"
	if _, err := os.Stat(srcWeb); err == nil {
		cmd := exec.Command("cp", "-r", srcWeb, InstallDir)
		return cmd.Run()
	}

	return nil
}

// createConfig membuat file .env
func (i *Installer) createConfig() error {
	config := fmt.Sprintf(`# Amaliah Ramadhan Configuration
APP_NAME=Amaliah Ramadhan
APP_PORT=%s
APP_ENV=production

# Database
DB_PATH=%s

# Security
JWT_SECRET=%s
`, i.Port, DatabaseFile, generateRandomKey())

	return os.WriteFile(ConfigFile, []byte(config), 0644)
}

// setupDatabase menjalankan migrasi database
func (i *Installer) setupDatabase() error {
	// Database akan dibuat otomatis saat aplikasi pertama kali jalan
	// Kita bisa jalankan aplikasi dalam mode migrate
	return nil
}

// createService membuat systemd service file
func (i *Installer) createService() error {
	service := fmt.Sprintf(`[Unit]
Description=Amaliah Ramadhan - Monitoring Ibadah Harian
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=%s
ExecStart=%s/amaliah-ramadhan
Restart=always
RestartSec=5
StandardOutput=append:/var/log/amaliah-ramadhan.log
StandardError=append:/var/log/amaliah-ramadhan.error.log

[Install]
WantedBy=multi-user.target
`, InstallDir, InstallDir)

	return os.WriteFile(ServiceFile, []byte(service), 0644)
}

// setPermissions mengatur hak akses file
func (i *Installer) setPermissions() error {
	// Set owner
	cmd := exec.Command("chown", "-R", "root:root", InstallDir)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Set executable permission
	binaryPath := filepath.Join(InstallDir, AppName)
	return os.Chmod(binaryPath, 0755)
}

// enableAndStartService mengaktifkan dan memulai systemd service
func (i *Installer) enableAndStartService() error {
	// Reload systemd
	cmd := exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Enable service
	cmd = exec.Command("systemctl", "enable", ServiceName)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Start service
	cmd = exec.Command("systemctl", "start", ServiceName)
	return cmd.Run()
}

// Update melakukan update aplikasi
func (i *Installer) Update() error {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  UPDATE/UPGRADE                          ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	if !i.CheckExistingInstallation() {
		return fmt.Errorf("tidak ada instalasi yang ditemukan. Gunakan 'Install Baru' terlebih dahulu")
	}

	if !i.Confirm("Lanjutkan update?") {
		return fmt.Errorf("update dibatalkan")
	}

	fmt.Println()
	fmt.Println("🔄 Memulai update...")
	fmt.Println()

	// Step 1: Stop service
	if err := i.ShowProgress("Menghentikan service", func() error {
		cmd := exec.Command("systemctl", "stop", ServiceName)
		return cmd.Run()
	}); err != nil {
		return err
	}

	// Step 2: Backup database
	if err := i.ShowProgress("Backup database", func() error {
		return i.backupDatabase()
	}); err != nil {
		return err
	}

	// Step 3: Update binary
	if err := i.ShowProgress("Update aplikasi", func() error {
		return i.copyBinary()
	}); err != nil {
		return err
	}

	// Step 4: Update static files
	if err := i.ShowProgress("Update file statis", func() error {
		return i.copyStaticFiles()
	}); err != nil {
		return err
	}

	// Step 5: Start service
	if err := i.ShowProgress("Memulai service kembali", func() error {
		cmd := exec.Command("systemctl", "start", ServiceName)
		return cmd.Run()
	}); err != nil {
		return err
	}

	i.showUpdateSuccessMessage()

	return nil
}

// Uninstall menghapus instalasi
func (i *Installer) Uninstall() error {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                    UNINSTALL                             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")

	if !i.CheckExistingInstallation() {
		fmt.Println()
		fmt.Println("ℹ️  Tidak ada instalasi yang ditemukan")
		return nil
	}

	fmt.Println()
	fmt.Println("⚠️  PERINGATAN: Semua data akan dihapus!")
	if !i.Confirm("Yakin ingin uninstall?") {
		return fmt.Errorf("uninstall dibatalkan")
	}

	fmt.Println()
	fmt.Println("🗑️  Memulai uninstall...")
	fmt.Println()

	// Step 1: Stop and disable service
	if err := i.ShowProgress("Menghentikan service", func() error {
		cmd := exec.Command("systemctl", "stop", ServiceName)
		cmd.Run() // Ignore error
		cmd = exec.Command("systemctl", "disable", ServiceName)
		return cmd.Run()
	}); err != nil {
		// Continue even if error
	}

	// Step 2: Remove service file
	if err := i.ShowProgress("Menghapus service file", func() error {
		return os.Remove(ServiceFile)
	}); err != nil {
		// Continue even if error
	}

	// Step 3: Remove installation directory
	if err := i.ShowProgress("Menghapus direktori instalasi", func() error {
		return os.RemoveAll(InstallDir)
	}); err != nil {
		return err
	}

	// Step 4: Reload systemd
	if err := i.ShowProgress("Reload systemd", func() error {
		cmd := exec.Command("systemctl", "daemon-reload")
		return cmd.Run()
	}); err != nil {
		// Continue even if error
	}

	fmt.Println()
	fmt.Println("✅ Uninstall selesai!")
	fmt.Println()

	return nil
}

// ShowStatus menampilkan status instalasi
func (i *Installer) ShowStatus() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  STATUS APLIKASI                         ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Check installation
	installed := i.CheckExistingInstallation()
	if installed {
		fmt.Println("📦 Instalasi      : ✅ Terinstall")
		fmt.Printf("📁 Direktori      : %s\n", InstallDir)
	} else {
		fmt.Println("📦 Instalasi      : ❌ Belum terinstall")
		return
	}

	// Check service status
	cmd := exec.Command("systemctl", "is-active", ServiceName)
	output, _ := cmd.Output()
	status := strings.TrimSpace(string(output))

	if status == "active" {
		fmt.Println("🟢 Service Status : ✅ Berjalan")
	} else {
		fmt.Println("🔴 Service Status : ❌ Tidak berjalan")
	}

	// Read config
	if configData, err := os.ReadFile(ConfigFile); err == nil {
		lines := strings.Split(string(configData), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "APP_PORT=") {
				port := strings.TrimPrefix(line, "APP_PORT=")
				fmt.Printf("🌐 Port           : %s\n", port)
				fmt.Printf("🔗 URL            : http://localhost:%s\n", port)
			}
		}
	}

	// Check database
	if _, err := os.Stat(DatabaseFile); err == nil {
		info, _ := os.Stat(DatabaseFile)
		fmt.Printf("💾 Database       : ✅ (%d KB)\n", info.Size()/1024)
	}

	// Show logs command
	fmt.Println()
	fmt.Println("📝 Lihat log:")
	fmt.Printf("   sudo journalctl -u %s -f\n", ServiceName)
	fmt.Println()
}

// backupDatabase membuat backup database
func (i *Installer) backupDatabase() error {
	if _, err := os.Stat(DatabaseFile); os.IsNotExist(err) {
		return nil // No database to backup
	}

	backupFile := fmt.Sprintf("%s.backup.%s", DatabaseFile, time.Now().Format("20060102-150405"))

	input, err := os.ReadFile(DatabaseFile)
	if err != nil {
		return err
	}

	return os.WriteFile(backupFile, input, 0644)
}

// generateRandomKey membuat random key untuk JWT
func generateRandomKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// showSuccessMessage menampilkan pesan sukses instalasi
func (i *Installer) showSuccessMessage() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║              ✅ INSTALASI BERHASIL! ✅                   ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("🎉 Amaliah Ramadhan telah terinstall dengan sukses!")
	fmt.Println()
	fmt.Printf("🌐 Aplikasi berjalan di: http://localhost:%s\n", i.Port)
	fmt.Println()
	fmt.Println("📝 Perintah berguna:")
	fmt.Printf("   • Status  : sudo systemctl status %s\n", ServiceName)
	fmt.Printf("   • Stop    : sudo systemctl stop %s\n", ServiceName)
	fmt.Printf("   • Start   : sudo systemctl start %s\n", ServiceName)
	fmt.Printf("   • Restart : sudo systemctl restart %s\n", ServiceName)
	fmt.Printf("   • Log     : sudo journalctl -u %s -f\n", ServiceName)
	fmt.Println()
	fmt.Println("📖 Login default:")
	fmt.Println("   Username: admin")
	fmt.Println("   Password: admin123")
	fmt.Println()
	fmt.Println("⚠️  Jangan lupa ubah password default setelah login!")
	fmt.Println()
}

// showUpdateSuccessMessage menampilkan pesan sukses update
func (i *Installer) showUpdateSuccessMessage() {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║                ✅ UPDATE BERHASIL! ✅                    ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("🎉 Aplikasi telah diupdate ke versi terbaru!")
	fmt.Println()
	fmt.Printf("🌐 Aplikasi berjalan di: http://localhost:%s\n", i.Port)
	fmt.Println()
	fmt.Println("💾 Database backup tersimpan di:")
	fmt.Printf("   %s.backup.*\n", DatabaseFile)
	fmt.Println()
}

// Run menjalankan installer wizard
func (i *Installer) Run() error {
	// Check permission
	if err := i.CheckPermission(); err != nil {
		fmt.Println()
		fmt.Printf("❌ Error: %s\n", err)
		fmt.Println()
		return err
	}

	// Show banner
	i.ShowBanner()

	// Main loop
	for {
		choice := i.ShowMenu()

		var err error
		switch choice {
		case 1:
			err = i.InstallNew()
		case 2:
			err = i.Update()
		case 3:
			err = i.Uninstall()
		case 4:
			i.ShowStatus()
			continue
		case 5:
			fmt.Println()
			fmt.Println("👋 Terima kasih!")
			fmt.Println()
			return nil
		default:
			fmt.Println()
			fmt.Println("❌ Pilihan tidak valid!")
			continue
		}

		if err != nil {
			fmt.Println()
			fmt.Printf("❌ Error: %s\n", err)
			fmt.Println()
		}

		// Ask to continue or exit
		fmt.Println()
		if !i.Confirm("Kembali ke menu utama?") {
			fmt.Println()
			fmt.Println("👋 Terima kasih!")
			fmt.Println()
			break
		}
	}

	return nil
}
