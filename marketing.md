# Sharing Project: Apps Monitoring Ramadhan

Dokumen ini berisi rencana konten Threads yang lebih santai untuk menceritakan proses pembuatan aplikasi ini. Tujuannya *sharing* karya, edukasi teknis (homeserver), dan ajakan berbagi (charity).

## 🎯 Target Pembaca
Teman-teman developer, guru, panitia Ramadhan, pegiat homelab/self-hosting.

## ✨ Poin Penting (Narrative)
*   **Tech**: Jalan di STB bekas (Armbian) + Cloudflare Tunnel. Low cost high impact.
*   **Purpose**: Digitalisasi amaliah Ramadhan siswa.
*   **Charity**: 100% donasi buat anak yatim & lansia. Gratis buat yang udah donasi di tempat lain.

---

## 🧵 Draft Konten Threads (Storytelling)

**Gaya bahasa**: Santai, teknis tapi ringan, transparan.

### Day 1: The "Why" (Masalah Klasik)
**Topik**: Buku kegiatan fisik yang sering ilang/rusak.
**Isi**:
> "Tiap Ramadhan, drama buku kegiatan ilang atau kecuci di saku celana tuh pasti ada. 😅
> Iseng mikir, 'Hari gini masa masih manual?'. Akhirnya weekend kemarin coba ngoding dikit buat digitalkin proses ini.
> Project iseng yang semoga berfaedah. Stay tuned, bakal gue spill progress-nya! 🛠️ #CodingRamadhan"

### Day 2: The "Tech Stack" (Homeserver Low Budget)
**Topik**: Pamer setup hemat pakai STB Bekas.
**Isi**:
> "Banyak yang tanya, 'Bang, servernya pake apa? Mahal gak?'
>
> Jawabannya: **Nggak sama sekali.**
>
> Ini cuma jalan di **STB Bekas Indihome** (B860H) yang di-install **Armbian**.
> Servernya taruh di rumah, akses publik pakai **Cloudflare Tunnel** (dapet SSL gratis, aman).
>
> Modal minim, manfaat maksimal. Siapa bilang digitaslisasi sekolah harus mahal? 😎
> *(Lampirkan foto STB yang lagi nyala led-nya)* #SelfHosted #Armbian #LowBudgetTech"

### Day 3: The "MVP" (Tampilan Aplikasi)
**Topik**: Demo UI yang simpel & ringan.
**Isi**:
> "Lanjut codingan kemarin. Ini penampakan dashboard-nya! ✨
>
> 1. Siswa login.
> 2. Ceklis amaliah (sholat, tadarus).
> 3. Realtime ke rekap guru.
>
> Ringan banget karena cuma serve HTML/Go biasa dari STB tadi. Gak perlu spek dewa.
> Gimana menurut kalian tampilannya? 😁"
*(Lampirkan screenshot UI)*

### Day 4: Gamification (Fastabiqul Khairat)
**Topik**: Fitur Leaderboard.
**Isi**:
> "Biar siswa gak bosen, gue tambahin bumbu **Gamification**. 🏆
>
> Ada poin & leaderboard. Konsepnya *Fastabiqul Khairat* (berlomba dalam kebaikan).
> Jadi siswa semangat ngejar target ibadah, bukan karena takut dihukum, tapi karena seru balapan (dalam hal positif) sama temen sekelasnya. 🕌"

### Day 5: For Teachers (Otomatisasi Laporan)
**Topik**: Membantu administrasi guru.
**Isi**:
> "Guru PAI biasanya paling pusing pasca-Lebaran: Ngoreksi tumpukan buku kegiatan.
>
> Di sini, gue bikinin fitur *Auto Report*. Sekali klik, rekap nilai satu bulan keluar (PDF/Excel).
> Biar Pak/Bu Guru bisa fokus bimbing siswa, urusan admin biar sistem yang kerjain. 🫡"

### Day 6: The Charity Model (Konsep Donasi)
**Topik**: Penjelasan kenapa ini gratis/donasi.
**Isi**:
> "Karena servernya cuma pake STB bekas di rumah (listrik irits), cost operasionalnya hampir nol.
>
> Jadi, gue putuskan aplikasi ini **GRATIS** buat sekolah/masjid manapun.
>
> Tapi, gue buka donasi sukarela.
> **100% hasil donasi akan disalurkan buat santunan Anak Yatim & Janda Jompo** di sekitar tempat tinggal gue tanggal 12 Maret nanti.
>
> Jadi developernya dapet pahala jariyah, kalian dapet aplikasi + pahala sedekah. Win-win solution kan? 🤝"

### Day 7: Launch & Call to Action
**Topik**: Rilis & Opsi Donasi Eksternal.
**Isi**:
> "Alhamdulillah, Apps Monitoring Ramadhan **READY TO USE**! 🎉
>
> Buat sekolah yang mau pakai, silakan DM. Gue bantu setup-kan, **GRATIS**.
>
> Syaratnya?
> 1. Cukup doain author-nya.
> 2. Kalau ada rezeki lebih, boleh ikut donasi buat santunan nanti.
> 3. ATAU, kalau kalian udah/mau donasi ke anak yatim di lingkungan kalian sendiri, tinggal **kirim dokumentasi/fotonya** ke gue. Itu udah cukup sebagai 'bayaran' jasa setup-nya.
>
> Yuk jadikan teknologi jalan buat berbagi! 🌙 #RamadhanTech #BerbagiItuIndah"

---

## 🤝 Detail Skema Donasi & Setup

Ini penjelasan detail buat yang nanya-nanya di DM nanti:

### 1. Kenapa Gratis/Murah Banget?
Karena infrastrukturnya mandiri (*self-hosted*).
*   **Hardware**: STB B860H (Bekas, Murah).
*   **OS**: Armbian (Open Source, Gratis).
*   **Network**: Cloudflare Tunnel (Gratis, Secure, dapet SSL).
*   **Listrik**: Kecil banget (5V adaptator).

### 2. Penyaluran Donasi
*   **Tujuan**: Santunan Anak Yatim & Janda Jompo.
*   **Lokasi**: Warga sekitar sekolah/rumah (Local Community Support).
*   **Waktu**: 12 Maret 2026.
*   **Transparansi**: Dokumentasi penyaluran akan di-share di Threads nanti.

### 3. Opsi "Bayar" Pakai Dokumentasi
Kami menyadari kebaikan itu universal. Jika sekolah/pengguna ingin berdonasi sendiri di lingkungannya:
*   Silakan lakukan santunan di tempat masing-masing.
*   Foto/videokan kegiatannya.
*   Kirim ke kami sebagai bukti.
*   Kami akan bantu setup aplikasi sampai jalan **GRATIS** (sebagai bentuk dukungan sesama penebar kebaikan).
