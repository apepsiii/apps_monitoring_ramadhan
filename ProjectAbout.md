# Project About - Rencana Pengembangan Amaliah Ramadhan App

## Visi dan Misi

### Visi
Menjadikan aplikasi monitoring amaliah Ramadhan yang terbaik untuk siswa, membantu mereka membangun kebiasaan ibadah yang konsisten dan meningkatkan spiritualitas selama bulan Ramadhan.

### Misi
1. Memudahkan siswa dalam tracking ibadah harian
2. Memberikan gambaran progress ibadah yang jelas
3. Mendorong kompetisi sehat antar siswa
4. Membantu admin/guru dalam monitoring amaliah siswa

## Fase Pengembangan

### Fase 1: Foundation (Minggu 1-2)
**Tujuan**: Menyiapkan fondasi aplikasi yang solid

#### Sprint 1: Setup & Konfigurasi
- [ ] Inisialisasi project Go dengan Echo framework
- [ ] Setup struktur folder sesuai standar
- [ ] Konfigurasi database (SQLite untuk dev, PostgreSQL untuk production)
- [ ] Setup Tailwind CSS dan konfigurasi build
- [ ] Setup environment variables
- [ ] Membuat base layout dan component HTML

#### Sprint 2: Authentication System
- [ ] Database migration untuk tabel users
- [ ] Implementasi register untuk siswa
- [ ] Implementasi login dengan JWT
- [ ] Middleware authentication
- [ ] Role-based access control (Admin vs User)
- [ ] Halaman login dan register dengan desain mobile-first
- [ ] Logout functionality
- [ ] Password hashing dengan bcrypt

### Fase 2: Core Features (Minggu 3-5)
**Tujuan**: Membangun fitur inti aplikasi

#### Sprint 3: Shalat Monitoring
- [ ] Database migration tabel prayers
- [ ] API endpoints untuk CRUD data shalat
- [ ] Form input shalat dengan 5 waktu
- [ ] Status shalat: jamaah, sendiri, belum, tidak
- [ ] Riwayat shalat harian
- [ ] Statistik shalat mingguan
- [ ] Reminder/notifikasi shalat (optional)

#### Sprint 4: Fasting Status
- [ ] Database migration tabel fastings
- [ ] API endpoints untuk tracking puasa
- [ ] Check-in puasa harian
- [ ] Status puasa: puasa, tidak (dengan alasan)
- [ ] Statistik puasa selama Ramadhan
- [ ] Progress bar puasa
- [ ] Visualisasi hari berpuasa

#### Sprint 5: Quran Reading Tracker
- [ ] Database migration tabel quran_readings
- [ ] API endpoints untuk tracking bacaan
- [ ] Input juz/halaman yang dibaca
- [ ] Progress khatam Quran
- [ ] Target harian bacaan
- [ ] Riwayat bacaan

### Fase 3: Engagement Features (Minggu 6-7)
**Tujuan**: Menambahkan fitur yang meningkatkan engagement user

#### Sprint 6: Daily Amaliah
- [ ] Database migration tabel amaliah_types dan daily_amaliah
- [ ] Master data jenis amaliah (sedekah, dzikir, dll)
- [ ] Form check-in amaliah harian
- [ ] Point/reward system untuk setiap amaliah
- [ ] Total poin user
- [ ] Riwayat amaliah

#### Sprint 7: Gamification & Leaderboard
- [ ] Ranking siswa berdasarkan poin
- [ ] Leaderboard mingguan dan keseluruhan
- [ ] Badge/achievement system
- [ ] Streak tracking (consistency)
- [ ] Notifikasi milestone

### Fase 4: Admin Dashboard (Minggu 8-9)
**Tujuan**: Membangun dashboard untuk admin/guru

#### Sprint 8: Admin Management
- [ ] Dashboard overview admin
- [ ] CRUD data siswa
- [ ] Import data siswa dari Excel/CSV
- [ ] Manajemen kelas/rombel
- [ ] Reset password siswa

#### Sprint 9: Reporting & Analytics
- [ ] Laporan amaliah per siswa
- [ ] Laporan aggregat per kelas
- [ ] Statistik shalat, puasa, quran
- [ ] Export laporan (PDF/Excel)
- [ ] Grafik dan visualisasi data
- [ ] Filter dan search data

### Fase 5: Polish & Launch (Minggu 10-11)
**Tujuan**: Mempersiapkan aplikasi untuk production

#### Sprint 10: UI/UX Enhancement
- [ ] Animasi dan transisi
- [ ] Loading states
- [ ] Error handling dan pesan error yang user-friendly
- [ ] Dark mode (optional)
- [ ] Offline support (PWA) (optional)
- [ ] Responsiveness testing di berbagai device

#### Sprint 11: Testing & Deployment
- [ ] Unit testing
- [ ] Integration testing
- [ ] Security audit
- [ ] Performance optimization
- [ ] Deployment ke server
- [ ] SSL configuration
- [ ] Backup strategy
- [ ] User manual dan dokumentasi

## Spesifikasi Teknis Detail

### Database Schema

#### Tabel Users
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    class VARCHAR(50),
    role ENUM('admin', 'user') DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabel Prayers
```sql
CREATE TABLE prayers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    date DATE NOT NULL,
    subuh ENUM('jamaah', 'sendiri', 'belum', 'tidak') DEFAULT 'belum',
    dzuhur ENUM('jamaah', 'sendiri', 'belum', 'tidak') DEFAULT 'belum',
    ashar ENUM('jamaah', 'sendiri', 'belum', 'tidak') DEFAULT 'belum',
    maghrib ENUM('jamaah', 'sendiri', 'belum', 'tidak') DEFAULT 'belum',
    isya ENUM('jamaah', 'sendiri', 'belum', 'tidak') DEFAULT 'belum',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### Tabel Fastings
```sql
CREATE TABLE fastings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    date DATE NOT NULL,
    status ENUM('puasa', 'tidak') DEFAULT 'puasa',
    reason VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### Tabel Quran Readings
```sql
CREATE TABLE quran_readings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    date DATE NOT NULL,
    juz_start INTEGER,
    juz_end INTEGER,
    pages INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### Tabel Amaliah Types
```sql
CREATE TABLE amaliah_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    points INTEGER DEFAULT 1,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Tabel Daily Amaliah
```sql
CREATE TABLE daily_amaliah (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    amaliah_type_id INTEGER NOT NULL,
    date DATE NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (amaliah_type_id) REFERENCES amaliah_types(id)
);
```

### Design System (Monnile Style)

#### Warna
- Primary: #6366F1 (Indigo)
- Secondary: #8B5CF6 (Purple)
- Success: #10B981 (Emerald)
- Warning: #F59E0B (Amber)
- Danger: #EF4444 (Red)
- Background: #F9FAFB (Gray 50)
- Surface: #FFFFFF (White)
- Text Primary: #111827 (Gray 900)
- Text Secondary: #6B7280 (Gray 500)

#### Typography
- Font Family: Inter / system-ui
- Heading: Bold (700)
- Body: Regular (400)
- Small: 12px, Regular
- Base: 14px, Regular
- Large: 16px, Medium
- XL: 18px, Semibold
- 2XL: 20px, Bold

#### Spacing
- Base unit: 4px
- xs: 4px
- sm: 8px
- md: 16px
- lg: 24px
- xl: 32px
- 2xl: 48px

#### Border Radius
- sm: 4px
- md: 8px
- lg: 12px
- xl: 16px
- full: 9999px

#### Shadows (Soft)
- sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05)
- md: 0 4px 6px -1px rgba(0, 0, 0, 0.1)
- lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1)

### Component Library

#### Mobile Components
1. **Bottom Navigation**
   - 4-5 menu items
   - Active state dengan ikon filled
   - Smooth transition

2. **Cards**
   - White background
   - Soft shadow
   - 12px border radius
   - 16px padding

3. **Buttons**
   - Primary: Indigo background, white text
   - Secondary: White background, indigo border
   - Ghost: Transparent with hover effect
   - Full width on mobile
   - 40px height minimum (touch friendly)

4. **Inputs**
   - Bottom border style atau boxed
   - Clear label
   - Error state with red border
   - Icon support

5. **Progress Indicators**
   - Circular progress untuk Quran
   - Linear progress untuk puasa
   - Streak counter dengan fire icon

6. **Lists**
   - Clean list dengan divider
   - Avatar/Icon di sebelah kiri
   - Arrow/Action di sebelah kanan

## Fitur Lanjutan (Future Development)

### Fase 6: Advanced Features (Post Launch)
- [ ] Push notification untuk reminder
- [ ] Social features (share progress)
- [ ] Integration dengan Google Calendar
- [ ] Multi-language support (Indonesia, English, Arabic)
- [ ] Dark mode toggle
- [ ] Widget untuk home screen
- [ ] Voice input untuk dzikir counter
- [ ] Location-based features

### Fase 7: Scale & Optimization
- [ ] Caching layer (Redis)
- [ ] CDN untuk static assets
- [ ] Database indexing dan optimization
- [ ] Load balancing
- [ ] Auto-scaling

## Metrics dan KPI

### Technical Metrics
- Page Load Time: < 2 detik
- API Response Time: < 500ms
- Uptime: > 99.5%
- Mobile Responsiveness Score: > 90

### Business Metrics
- Daily Active Users (DAU)
- Retention Rate (Day 1, Day 7, Day 30)
- Feature adoption rate
- User satisfaction score
- Average session duration

## Risk Assessment

### Risks
1. **Timeline Risk**: Ramadhan sudah dekat (Feb 2026)
2. **Technical Risk**: Tim developer mungkin terbatas
3. **User Adoption**: Siswa mungkin tidak familiar dengan teknologi
4. **Server Load**: Banyak user login bersamaan

### Mitigation
1. Prioritaskan core features (MVP approach)
2. Gunakan template/framework yang sudah ada
3. Buat user guide yang simple dan mudah dipahami
4. Load testing sebelum launch

## Budget Estimation (Simplified)

### Development
- Server/VPS: $20-50/bulan
- Domain: $10-15/tahun
- SSL Certificate: Free (Let's Encrypt)
- Development Tools: Free/Open Source

### Timeline Summary
- Total Development Time: 10-11 minggu
- Target Launch: 10 Februari 2026 (1 minggu sebelum Ramadhan)
- Start Development: 1 Desember 2025

## Success Criteria

Aplikasi dianggap sukses jika:
1. 80% siswa aktif menggunakan aplikasi setiap hari
2. Tidak ada bug critical yang mengganggu penggunaan
3. Aplikasi dapat handle load saat ramai pengguna
4. Feedback positif dari siswa dan admin
5. Meningkatnya konsistensi ibadah siswa

## Kesimpulan

Project ini bertujuan untuk membangun aplikasi monitoring amaliah Ramadhan yang user-friendly, mobile-first, dan dapat membantu siswa dalam meningkatkan kualitas ibadah mereka. Dengan timeline 10-11 minggu dan approach MVP, kita dapat meluncurkan aplikasi tepat waktu untuk Ramadhan 2026.

**Next Steps:**
1. Finalisasi requirement dengan stakeholder
2. Setup development environment
3. Mulai Fase 1: Foundation
4. Daily/weekly standup untuk tracking progress
