package config

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(100) NOT NULL,
			class VARCHAR(50),
			role VARCHAR(10) DEFAULT 'user',
			points INTEGER DEFAULT 0,
			avatar VARCHAR(255) DEFAULT 'default',
			bio TEXT,
			theme VARCHAR(20) DEFAULT 'emerald',
			target_khatam INTEGER DEFAULT 30,
			provinsi VARCHAR(100) DEFAULT '',
			kabkota VARCHAR(100) DEFAULT '',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS prayers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			subuh VARCHAR(20) DEFAULT 'belum',
			dzuhur VARCHAR(20) DEFAULT 'belum',
			ashar VARCHAR(20) DEFAULT 'belum',
			maghrib VARCHAR(20) DEFAULT 'belum',
			isya VARCHAR(20) DEFAULT 'belum',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS fastings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			status VARCHAR(20) DEFAULT 'puasa',
			reason VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS quran_readings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			date DATE NOT NULL,
			start_surah_id INTEGER,
			start_surah_name VARCHAR(100),
			start_ayah INTEGER,
			end_surah_id INTEGER,
			end_surah_name VARCHAR(100),
			end_ayah INTEGER,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS amaliah_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			points INTEGER DEFAULT 1,
			icon VARCHAR(50),
			is_active BOOLEAN DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS daily_amaliah (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			amaliah_type_id INTEGER NOT NULL,
			date DATE NOT NULL,
			notes TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (amaliah_type_id) REFERENCES amaliah_types(id)
		)`,
		`INSERT OR IGNORE INTO amaliah_types (name, description, points, icon) VALUES 
			('Sedekah', 'Bersedekah kepada orang yang membutuhkan', 10, 'heart'),
			('Dzikir Pagi', 'Dzikir pagi setelah subuh', 5, 'sun'),
			('Dzikir Petang', 'Dzikir petang setelah maghrib', 5, 'moon'),
			('Sholat Dhuha', 'Melaksanakan sholat dhuha', 7, 'sunrise'),
			('Sholat Tahajud', 'Melaksanakan sholat tahajud', 10, 'star'),
			('Baca Al-Quran', 'Membaca Al-Quran minimal 1 halaman', 5, 'book'),
			('Istighfar', 'Beristighfar minimal 100x', 3, 'refresh'),
			('Sholawat', 'Bersholawat minimal 100x', 5, 'message'),
			('Bantu Orang Tua', 'Membantu orang tua di rumah', 5, 'home'),
			('Tahfidz', 'Menghafal Al-Quran', 15, 'book-open')
		`,
	}

	classMigrations := []string{
		`CREATE TABLE IF NOT EXISTS classes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(50) UNIQUE NOT NULL,
			level VARCHAR(20),
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`INSERT OR IGNORE INTO classes (name, level, description) VALUES 
			('X-RPL', 'X', 'Rekayasa Perangkat Lunak'),
			('XI-RPL', 'XI', 'Rekayasa Perangkat Lunak'),
			('XII-RPL', 'XII', 'Rekayasa Perangkat Lunak')
		`,
	}
	migrations = append(migrations, classMigrations...)

	badgeMigrations := []string{
		`CREATE TABLE IF NOT EXISTS badges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			icon VARCHAR(50),
			criteria_type VARCHAR(50),
			criteria_value INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS user_badges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			badge_id INTEGER NOT NULL,
			earned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (badge_id) REFERENCES badges(id),
			UNIQUE(user_id, badge_id)
		)`,
		`INSERT OR IGNORE INTO badges (name, description, icon, criteria_type, criteria_value) VALUES 
			('Awal Langkah', 'Menyelesaikan shalat 5 waktu pertama kali', 'star', 'prayer_count', 1),
			('Istiqomah 7 Hari', 'Shalat 5 waktu berturut-turut selama 7 hari', 'fire', 'prayer_streak', 7),
			('Istiqomah 30 Hari', 'Shalat 5 waktu berturut-turut selama 30 hari', 'award', 'prayer_streak', 30),
			('Pembaca Al-Quran', 'Mulai membaca Al-Quran', 'book-open', 'quran_readings', 1),
			('Khatam 1 Juz', 'Menyelesaikan 1 Juz Al-Quran', 'book', 'quran_juz', 1),
			('Khatam Al-Quran', 'Menyelesaikan 30 Juz Al-Quran', 'check-circle', 'quran_khatam', 30),
			('Dermawan', 'Mendapatkan 100 poin amaliah', 'heart', 'amaliah_points', 100),
			('Ahli Ibadah', 'Mendapatkan 500 poin amaliah', 'sun', 'amaliah_points', 500),
			('Sang Juara', 'Mendapatkan 1000 poin amaliah', 'trophy', 'amaliah_points', 1000)
		`,
	}

	migrations = append(migrations, badgeMigrations...)

	// Performance Indexes
	indexMigrations := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`,
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE INDEX IF NOT EXISTS idx_prayers_user_date ON prayers(user_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_fastings_user_date ON fastings(user_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_quran_user_date ON quran_readings(user_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_amaliah_user_date ON daily_amaliah(user_id, date)`,
		`CREATE INDEX IF NOT EXISTS idx_user_badges_user ON user_badges(user_id)`,
	}
	migrations = append(migrations, indexMigrations...)

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			log.Printf("Migration %d failed: %v", i+1, err)
			return err
		}
		log.Printf("Migration %d executed successfully", i+1)
	}

	if err := addColumnIfNotExists(db, "users", "provinsi", "VARCHAR(100) DEFAULT ''"); err != nil {
		log.Printf("Note: %v", err)
	}
	if err := addColumnIfNotExists(db, "users", "kabkota", "VARCHAR(100) DEFAULT ''"); err != nil {
		log.Printf("Note: %v", err)
	}

	if err := seedAdminUser(db); err != nil {
		log.Printf("Failed to seed admin user: %v", err)
	}

	return nil
}

func addColumnIfNotExists(db *sql.DB, table, column, definition string) error {
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?", table, column).Scan(&exists)
	if err != nil {
		return err
	}

	if exists == 0 {
		_, err = db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, definition))
		if err != nil {
			return fmt.Errorf("failed to add column %s: %v", column, err)
		}
		log.Printf("Added column %s to table %s", column, table)
	}
	return nil
}

func seedAdminUser(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password_hash, full_name, class, role, points, avatar, bio, theme) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, "admin", "admin@ramadhan.com", string(hashedPassword), "Administrator", "", "admin", 0, "default", "", "emerald")
	if err != nil {
		return err
	}

	log.Println("Default admin user created: admin / admin123")
	return nil
}
