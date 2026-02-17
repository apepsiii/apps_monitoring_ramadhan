# Amaliah Ramadhan Monitoring App

Aplikasi monitoring amaliah Ramadhan untuk siswa dengan tampilan mobile-first yang memudahkan tracking ibadah harian selama bulan suci Ramadhan.

## 🚀 Fitur

### Untuk Siswa (User)
- ✅ **Shalat Monitoring** - Tracking 5 waktu shalat harian
- ✅ **Status Puasa** - Check-in puasa dengan berbagai status
- ✅ **Bacaan Al-Quran** - Tracking progress khatam Quran
- ✅ **Amaliah Harian** - Tracking kebaikan dan poin reward
- ✅ **Leaderboard** - Kompetisi sehat antar siswa

### Untuk Admin
- ✅ **Dashboard Admin** - Overview statistik
- ✅ **Manajemen Siswa** - CRUD data siswa
- ✅ **Laporan** - Laporan amaliah per siswa
- ✅ **Statistik** - Grafik dan analisis data

## 🛠️ Tech Stack

- **Backend**: Go (Golang) dengan Echo Framework
- **Database**: SQLite (Development) / PostgreSQL (Production)
- **Frontend**: HTML Template + Tailwind CSS
- **Authentication**: JWT (JSON Web Token)
- **Styling**: Tailwind CSS (Mobile-First Design)

## 📋 Prerequisites

- Go 1.21 atau lebih tinggi
- Node.js (untuk Tailwind CSS)
- SQLite3

## 🚀 Getting Started

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/amaliah-monitoring.git
cd amaliah-monitoring
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod tidy

# Install Node.js dependencies
npm install
```

### 3. Setup Environment Variables

```bash
cp .env.example .env
# Edit .env file sesuai kebutuhan
```

### 4. Build CSS

```bash
# Development (watch mode)
npm run dev

# Production
npm run build:css
```

### 5. Run Application

```bash
go run cmd/main.go
```

Aplikasi akan berjalan di `http://localhost:8080`

## 📁 Project Structure

```
/
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

## 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Nama aplikasi | Amaliah Ramadhan App |
| `APP_ENV` | Environment | development |
| `APP_PORT` | Port server | 8080 |
| `DB_DRIVER` | Database driver | sqlite |
| `DB_NAME` | Database name/path | ./amaliah.db |
| `JWT_SECRET` | Secret key JWT | default-secret |

## 🎨 UI/UX Design

Aplikasi menggunakan design system **Monnile-style** dengan karakteristik:
- Card-based layout
- Bottom navigation
- Soft shadows dan rounded corners
- Pastel color palette
- Mobile-first responsive design

## 🔐 Authentication

- **User**: Login dengan username dan password
- **Admin**: Role-based access control
- **Session**: JWT token dengan cookie

## 📝 API Endpoints

### Authentication
- `GET /login` - Halaman login
- `POST /login` - Proses login
- `GET /register` - Halaman register
- `POST /register` - Proses register
- `GET /logout` - Logout

### User Routes
- `GET /user/dashboard` - Dashboard user
- `GET /user/prayers` - Tracking shalat
- `POST /user/prayers` - Simpan data shalat
- `GET /user/fasting` - Status puasa
- `POST /user/fasting` - Simpan status puasa
- `GET /user/quran` - Bacaan Quran
- `POST /user/quran` - Simpan bacaan
- `GET /user/amaliah` - Amaliah harian
- `POST /user/amaliah` - Simpan amaliah

### Admin Routes
- `GET /admin/dashboard` - Dashboard admin
- `GET /admin/users` - Manajemen siswa
- `POST /admin/users` - Tambah siswa
- `GET /admin/reports` - Laporan
- `GET /admin/statistics` - Statistik

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run specific package
go test ./internal/...
```

## 📦 Deployment

### Build untuk Production

```bash
# Build binary
CGO_ENABLED=1 go build -o bin/amaliah-app cmd/main.go

# Build CSS
npm run build:css

# Run
./bin/amaliah-app
```

### Docker (Opsional)

```bash
# Build image
docker build -t amaliah-app .

# Run container
docker run -p 8080:8080 amaliah-app
```

## 🤝 Contributing

1. Fork repository
2. Buat branch feature (`git checkout -b feature/amazing-feature`)
3. Commit perubahan (`git commit -m 'Add amazing feature'`)
4. Push ke branch (`git push origin feature/amazing-feature`)
5. Buat Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 👥 Team

- Backend Developer
- Frontend Developer
- UI/UX Designer

## 📞 Support

Jika ada pertanyaan atau masalah, silakan buat issue di repository ini.

---

**Ramadhan 1447 H / 2026 M**

*Dibuat dengan ❤️ untuk memudahkan ibadah Ramadhan*
