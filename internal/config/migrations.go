package config

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`DROP TABLE IF EXISTS users`,
		`CREATE TABLE users (
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
		`DROP TABLE IF EXISTS quran_readings`,
		`CREATE TABLE quran_readings (
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

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			log.Printf("Migration %d failed: %v", i+1, err)
			return err
		}
		log.Printf("Migration %d executed successfully", i+1)
	}

	return nil
}
