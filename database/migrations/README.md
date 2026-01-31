# Panduan Migrasi Database

Dokumen ini menjelaskan cara menjalankan migrasi MySQL untuk tabel auth (users dan email verification token). Informasi ditulis dengan bahasa sederhana untuk developer yang belum familiar.

## 1. Persiapan
- Pastikan MySQL sudah berjalan di `127.0.0.1:3306`.
- Ada database bernama `onboarding`.
- Username/password: `root` / `root`.
- Instal tool [golang-migrate](https://github.com/golang-migrate/migrate) dulu (misal: `brew install golang-migrate` atau `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`).

## 2. Struktur Migrasi
Folder `database/migrations/` berisi file versi:
- `0001_create_users_table.up.sql` → membuat tabel `users`.
- `0001_create_users_table.down.sql` → menghapus tabel `users`.
- `0002_create_email_verification_tokens_table.up.sql` → membuat tabel token.
- `0002_create_email_verification_tokens_table.down.sql` → menghapus tabel token.

Setiap file `.up.sql` digunakan saat migrasi naik (menambahkan), sedangkan `.down.sql` untuk rollback (kembali satu langkah).

## 3. Cara Menjalankan Migrasi
Gunakan command berikut:

```bash
migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/onboarding?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations up
```

Perintah ini membuat tabel `users` dan `email_verification_tokens` sesuai urutan versi.

## 4. Cara Rollback (Turun satu langkah)
Kalau mau membatalkan migrasi terakhir:

```bash
migrate -database "mysql://root:root@tcp(127.0.0.1:3306)/onboarding?charset=utf8mb4&parseTime=True&loc=Local" -path database/migrations down 1
```

Perintah ini menjalankan file `.down.sql` versi terakhir (misal menghapus `email_verification_tokens` jika sebelumnya baru dibuat).

## 5. Catatan Tambahan
- Semua tabel menggunakan `CHAR(36)` untuk UUID, `utf8mb4`, dan `InnoDB` agar aman untuk production.
- Migration ini hanya untuk auth (users + verification). Tabel onboarding akan ditambahkan nanti.
- Jangan pakai `AutoMigrate` karena kita pakai file SQL versioned agar kontrol lebih ketat.
