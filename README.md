<p align="center">
  <img src="https://blog.golang.org/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" width="500" />
</p>

# Merchant Onboarding Backend

## Teknologi dan Struktur
- **Gin**: router HTTP utama untuk endpoint RESTful.
- **GORM & MySQL**: model domain (auth + onboarding) diterjemahkan ke tabel `users` dan `email_verification_tokens` serta didukung migrasi SQL versioned (`database/migrations`).
- **golang-migrate**: jalankan skrip up/down untuk menjaga versi schema.
- **Swagger (swaggo)**: dokumentasi `/swagger/index.html` dengan BearerAuth untuk endpoint auth.
- **SMTP Email Service**: kirim link verifikasi (env `SMTP_*` plus `APP_URL`).

<div align="center">
  <img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" width="180" />
  <img src="https://miro.medium.com/v2/resize:fit:1400/format:webp/1*XBvxUxqycRC8B8KGCuzJVw.png" width="220" />
  <img src="https://dev.mysql.com/common/logos/logo-mysql-170x115.png" width="220" />
</div>

## Flow Auth (sesuai API contract)
1. **Register `/api/auth/register`**
   - Payload: `full_name`, `business_name`, `email`, `phone_number`, `password`, `confirm_password`.
   - Simpan user dengan status `PENDING_EMAIL_VERIFICATION`, buat token 24 jam, kirim email via SMTP, log link.
2. **Verify Email `/api/auth/verify-email?token=xxx`**
   - Validasi token belum expired/used.
   - Ubah status user → `ACTIVE`, `email_verified = true`, tandai token used.
3. **Login `/api/auth/login`**
   - Pastikan user `ACTIVE` & `email_verified`.
   - Kembalikan JWT Bearer (`access_token`, `expires_in`, `token_type`, profile).
4. **Get Me `/api/auth/me`**
   - Middleware JWT memvalidasi token, handler ambil user dari repo, response profil sesuai kontrak.

## Merchant Onboarding (level overview)
- Data disimpan per langkah: business entity (merchant/owner/PIC/settlement), payment setup, terms, dokumen.
- Status awal `DRAFT`, hanya `POST /api/onboarding/{id}/request-approval` (menjaga aturan) yang boleh mengubah status ke `SUBMITTED_FOR_APPROVAL`.
- Frontend wajib memanggil endpoint status sebelum submit; backend menolak jika belum lengkap.

## Setup dan DX
- Copy `.env.example` → `.env`, isi:
  ```env
  JWT_SECRET=...
  DB_DSN=root:root@tcp(127.0.0.1:3306)/onboarding?charset=utf8mb4&parseTime=True&loc=Local
  APP_URL=http://localhost:8080
  SMTP_HOST=smtp.gmail.com
  SMTP_PORT=587
  SMTP_USERNAME=...
  SMTP_PASSWORD=...
  SMTP_FROM=...
  ```
- Jalankan migrations:
  ```bash
  migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/onboarding?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations up
  ```
- Start server: `go run cmd/api/main.go`.
- Swagger: `http://localhost:8080/swagger/index.html` (BearerAuth login).

## Catatan terakhir
- Semua error ditangani dengan kode HTTP sesuai kontrak (400/401/403/409/500).
- Response selalu `{success,message,data}` atau objek error serupa.
- SMTP service hanya di-usecase, handler tetap tipis; config lengkap di env.
