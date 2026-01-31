# Ringkasan Analisis Project MOB Backend

Dokumen ini merangkum struktur, alur, dan pola desain berdasarkan kode yang ada dan `api-contract.md`. Fokus: membantu developer baru memahami project dalam sekali baca.

## A. Struktur Project

- `cmd/`
  - Entry point aplikasi. `cmd/api/main.go` melakukan bootstrap server Gin, konfigurasi dependency, dan registrasi route.
- `internal/`
  - Berisi implementasi aplikasi yang tidak diekspor keluar module (idiomatic Go).
  - `internal/domain/` → model/domain entity (contoh: user auth, token verifikasi).
  - `internal/dto/` → struktur request/response sesuai kontrak API.
  - `internal/handler/` → HTTP handler Gin; hanya parsing request, validasi ringan, panggil usecase, lalu kirim response.
  - `internal/usecase/` → business logic (auth flow).
  - `internal/repository/` → abstraksi persistence (interface) dan implementasi in-memory.
  - `internal/middleware/` → middleware lintas endpoint (JWT auth).

**Kenapa struktur ini dipakai (Gin backend):**
- Memisahkan concern agar handler tetap tipis.
- Business logic bisa dites dan diganti storage tanpa mengubah handler.
- Mudah dikembangkan ke onboarding flow berikutnya tanpa mengacaukan auth.

## B. Flow Pattern Aplikasi (Request → Response)

**Gambaran umum:**
```
Client → Router (Gin) → Handler → Usecase → Repository → (kembali) → Handler → Response JSON
```

Peran utama:
- `main.go`: bootstrap server, baca env (`JWT_SECRET`), inisialisasi repo/usecase/handler, register route.
- Router (Gin): mapping endpoint ke handler.
- Handler: parsing JSON, validasi awal (binding), mapping error ke HTTP status + response format.
- Usecase: aturan bisnis (auth flow).
- Repository: akses data (in-memory map saat ini).

Penjelasan singkat untuk FE:
- Request masuk ke endpoint.
- Handler membaca data JSON dan memanggil usecase.
- Usecase menjalankan aturan bisnis (cek password, status user, token, dll).
- Repository menyimpan/ambil data.
- Handler mengirim response format `{ success, message, data }`.

## C. End-to-End Flow (Auth & Onboarding)

Referensi utama: `api-contract.md` + implementasi auth di kode.

### 1) REGISTER
**Contract:** `POST /api/auth/register`

Flow di code:
- Handler `Register` di `internal/handler/auth/handler.go`.
- Usecase `Register` di `internal/usecase/auth/service.go`:
  - Validasi password == confirm_password.
  - Cek email sudah terdaftar.
  - Hash password (bcrypt).
  - Simpan user status `PENDING_EMAIL_VERIFICATION`.
  - Generate token verifikasi (UUID, TTL 24 jam) dan simpan.
  - Log link verifikasi ke console (simulasi email).

### 2) EMAIL VERIFICATION
**Contract:** `GET /api/auth/verify-email?token=...`

Flow di code:
- Handler `VerifyEmail` → usecase `VerifyEmail`.
- Usecase:
  - Cari token.
  - Validasi belum dipakai dan belum expired.
  - Update user: status `ACTIVE`, `EmailVerified=true`.
  - Token jadi one-time (marked used).

### 3) LOGIN
**Contract:** `POST /api/auth/login`

Flow di code:
- Handler `Login` → usecase `Login`.
- Usecase:
  - Cari user berdasarkan email.
  - Validasi password (bcrypt compare).
  - Tolak jika status belum `ACTIVE` atau `EmailVerified=false`.
  - Generate JWT access token (HS256).

### 4) GET CURRENT USER (ME)
**Contract:** `GET /api/auth/me`

Flow di code:
- Middleware `AuthMiddleware` validasi JWT, inject `userID` ke context.
- Handler `Me` → usecase `GetMe` → repository.

### 5) Kapan & bagaimana onboarding merchant dimulai
- **Belum diimplementasikan.**
- Onboarding flow ada di `api-contract.md`, namun belum ada handler/usecase/repository untuk endpoint onboarding.

### 6) Status DRAFT dipertahankan
- **Belum diimplementasikan.**
- Kontrak menegaskan status `DRAFT` tidak diubah oleh endpoint penyimpanan data dan hanya berubah lewat `request-approval`.

### 7) Data disimpan per step
- **Belum diimplementasikan.**
- Kontrak menyiapkan endpoint terpisah untuk save per step (business entity, payment setup, terms). Implementasi belum ada.

### 8) Proses SUBMIT / REQUEST APPROVAL
- **Belum diimplementasikan.**
- Detail flow hanya ada di `api-contract.md`.

## D. Contoh Pembuatan Feature (Real Case)

**Contoh: Register** (feature yang sudah ada)

- **Endpoint mulai:** `POST /api/auth/register`
- **File terlibat:**
  - `cmd/api/main.go` → register route.
  - `internal/handler/auth/handler.go` → `Register` handler.
  - `internal/usecase/auth/service.go` → `Register` business logic.
  - `internal/repository/auth/memory_user_repo.go` → simpan user.
  - `internal/repository/auth/memory_token_repo.go` → simpan token.
  - `internal/domain/auth/user.go` + `internal/domain/auth/email_verification.go` → entity.
  - `internal/dto/auth.go` → request/response DTO.

**Urutan kerja (router → handler → usecase → repo):**
1. Router mapping `/api/auth/register` di `main.go`.
2. Handler membaca JSON ke `RegisterRequest`.
3. Usecase jalankan validasi + hash + simpan user + simpan token.
4. Repository in-memory menyimpan data.
5. Handler kirim response 201 dengan payload sesuai kontrak.

**Data masuk:** `full_name`, `business_name`, `email`, `phone_number`, `password`, `confirm_password`, `merchant_sales`, `referral_code (opsional)`.

**Data keluar:** `user_id`, `email`, `status` dengan message sukses.

**Kenapa pattern ini dipakai:**
- Memisahkan HTTP concern dari aturan bisnis.
- Usecase tetap bersih dan mudah dites.
- Repository bisa diganti database tanpa ubah handler.

## E. Prinsip Desain yang Dipakai

- **Status tidak diubah langsung oleh frontend**
  - Berdasarkan kontrak, status hanya berubah lewat aksi tertentu agar proses onboarding konsisten dan bisa diaudit.

- **Hanya satu endpoint boleh mengubah status**
  - Kontrak menetapkan `POST /request-approval` sebagai satu-satunya perubahan status; mencegah status loncat/inkonsisten.

- **Draft dipisah per step**
  - Kontrak memisahkan data per section (business entity, payment, terms) agar tiap langkah bisa disimpan bertahap tanpa mengubah status.

- **Handler tidak berisi logic berat**
  - Handler fokus parsing/response; bisnis di usecase agar maintainable dan scalable.

## Catatan Status Implementasi

- **Auth flow:** sudah ada (register, verify email, login, me).
- **Onboarding flow:** **Belum diimplementasikan** (semua endpoint onboarding di kontrak masih placeholder di sisi code).
