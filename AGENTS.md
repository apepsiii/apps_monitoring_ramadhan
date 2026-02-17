# Project Lookup - Amaliah Ramadhan Monitoring App

## Overview
Aplikasi monitoring amaliah Ramadhan untuk siswa dengan tampilan mobile-first yang memudahkan tracking ibadah harian selama bulan suci Ramadhan.

## Tech Stack

### Backend
- **Language**: Go (Golang)
- **Framework**: Echo (web framework)
- **Database**: SQLite / PostgreSQL
- **Authentication**: JWT (JSON Web Token)
- **Password Hashing**: bcrypt

### Frontend
- **Template Engine**: Go HTML Template
- **Styling**: Tailwind CSS (mobile-first design)
- **JavaScript**: Vanilla JS dengan HTMX untuk interaktivitas
- **UI Reference**: Monnile-style mobile components

### Features

#### Core Features
1. **Shalat Monitoring**
   - Tracking 5 waktu shalat (Subuh, Dzuhur, Ashar, Maghrib, Isya)
   - Status: Sudah/Sholat Berjamaah/Sholat Sendiri/Belum/Tidak Sholat
   - Riwayat harian dan mingguan

2. **Status Puasa**
   - Check-in puasa harian
   - Status: Puasa/Tidak Puasa (sakit/perjalanan/haid/dll)
   - Statistik puasa selama Ramadhan

3. **Bacaan Quran**
   - Tracking juz yang dibaca
   - Progress bacaan harian
   - Target khatam Quran

4. **Amaliah Kebaikan Harian**
   - Tracking amaliah positif (sedekah, dzikir, dll)
   - Point/reward system
   - Leaderboard siswa

5. **User Management**
   - Login/Register siswa
   - Role-based access (Admin & User)
   - Profile management

#### Admin Features
- Dashboard admin
- Manajemen data siswa (CRUD)
- Laporan dan statistik amaliah
- Export data

## Project Structure

```
/ (root)
├── cmd/
│   └── main.go                 # Entry point aplikasi
├── internal/
│   ├── config/                 # Konfigurasi aplikasi
│   ├── handlers/               # HTTP handlers (controllers)
│   ├── middleware/             # Echo middleware
│   ├── models/                 # Database models
│   ├── repository/             # Database queries
│   ├── services/               # Business logic
│   └── utils/                  # Helper functions
├── web/
│   ├── static/                 # CSS, JS, Images
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── templates/              # HTML templates
│       ├── layouts/
│       ├── partials/
│       ├── auth/
│       ├── admin/
│       └── user/
├── migrations/                 # Database migrations
├── tests/                      # Unit tests
├── docs/                       # Dokumentasi
├── .env.example               # Environment variables template
├── go.mod
├── go.sum
└── README.md
```

## Database Schema (Overview)

### Tables
1. `users` - Data user (siswa & admin)
2. `prayers` - Data shalat harian
3. `fastings` - Data puasa harian
4. `quran_readings` - Data bacaan Quran
5. `daily_amaliah` - Data amaliah kebaikan
6. `amaliah_types` - Master data jenis amaliah

## API Endpoints Structure

### Authentication
- POST /api/auth/login
- POST /api/auth/register
- POST /api/auth/logout
- GET /api/auth/me

### User Routes
- GET /api/user/dashboard
- GET /api/user/prayers
- POST /api/user/prayers
- GET /api/user/fasting
- POST /api/user/fasting
- GET /api/user/quran
- POST /api/user/quran
- GET /api/user/amaliah
- POST /api/user/amaliah

### Admin Routes
- GET /api/admin/dashboard
- GET /api/admin/users
- POST /api/admin/users
- PUT /api/admin/users/:id
- DELETE /api/admin/users/:id
- GET /api/admin/reports
- GET /api/admin/statistics

## UI/UX Design Principles

### Mobile-First Approach
- Design untuk mobile terlebih dahulu
- Responsive breakpoints
- Touch-friendly interface

### Monnile-Style Elements
- Card-based layout
- Bottom navigation
- Floating action buttons
- Clean typography
- Soft shadows
- Rounded corners
- Pastel color palette

## Development Environment

### Prerequisites
- Go 1.21+
- Node.js (untuk Tailwind CSS)
- SQLite / PostgreSQL

### Running Locally
```bash
# Install dependencies
go mod tidy

# Run migrations
go run cmd/migrate/main.go

# Start development server
go run cmd/main.go

# Run Tailwind CSS compiler (parallel)
npm run dev
```

## Deployment
- Platform: Railway / Render / VPS
- Environment: Production
- SSL: Let's Encrypt
- Reverse Proxy: Nginx

## Team Roles
- Backend Developer
- Frontend Developer
- UI/UX Designer

## Timeline
- Week 1: Setup & Database Design
- Week 2: Authentication & User Management
- Week 3: Core Features (Shalat, Puasa, Quran)
- Week 4: Amaliah & Admin Dashboard
- Week 5: UI Polish & Testing
- Week 6: Deployment & Documentation

## Notes
- Ramadhan 2026 dimulai sekitar 17 Februari 2026
- Target launch: 1 minggu sebelum Ramadhan
- Priority: Core features (shalat, puasa, quran) harus stabil terlebih dahulu
