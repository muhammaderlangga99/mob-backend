# API Contract for Merchant Onboarding

## Overview
This document outlines the API contract for the Merchant Onboarding process, detailing the endpoints, request payloads, and responses.

---

# 🔐 AUTH FLOW API CONTRACTS

## Auth Flow Overview

```
1️⃣ User SIGN UP (Register)
   ↓
2️⃣ Backend kirim EMAIL VERIFICATION (link)
   ↓
3️⃣ User klik link → Email terverifikasi
   ↓
4️⃣ User bisa SIGN IN
   ↓
5️⃣ Backend return ACCESS TOKEN
   ↓
6️⃣ User dapat mengakses dashboard dengan token
```

---

## 1️⃣ ENDPOINT: REGISTER (SIGN UP)

### POST /api/auth/register

**Description:** User mendaftar sebagai merchant baru

**Request Payload:**
```json
{
  "full_name": "Budi Santoso",
  "business_name": "Toko Jaya Sejahtera",
  "email": "budi@email.com",
  "phone_number": "08123456789",
  "password": "securePassword123",
  "confirm_password": "securePassword123",
  "merchant_sales": "2"
}
```

**Response (Success - 201):**
```json
{
  "success": true,
  "message": "Registration successful. Please verify your email.",
  "data": {
    "user_id": "uuid",
    "email": "budi@email.com",
    "status": "PENDING_EMAIL_VERIFICATION"
  }
}
  ---

  ## ⚙️ Status Flow & Actions

  ### Prinsip Utama
  - Status tidak boleh diubah manual oleh frontend.
  - Status hanya berubah melalui aksi (action) yang ditetapkan.
  - Semua endpoint penyimpanan data (business-entity, payment-feature, terms) **TIDAK** mengubah status.
  - Status awal: `DRAFT`.

  ### Hanya Satu Endpoint untuk Mengajukan Persetujuan
  - Endpoint yang berwenang mengubah status adalah:
    - `POST /api/onboarding/{merchant_onboarding_id}/request-approval`
  - Endpoint lain **tidak** boleh mengubah status.

  ## 9️⃣ ENDPOINT: REQUEST APPROVAL (UBAH STATUS)

  ### POST /api/onboarding/{merchant_onboarding_id}/request-approval

  Description: Mengajukan proses onboarding untuk approval. Endpoint ini akan mengubah status onboarding dari `DRAFT` menjadi `SUBMITTED_FOR_APPROVAL` jika validasi terpenuhi.

  Path Parameters:
  - `merchant_onboarding_id` (required)

  Request Payload (optional, dapat berisi catatan pengaju):
  ```json
  {
    "requested_by": {
      "user_id": "user-uuid-12345",
      "email": "budi@email.com",
      "full_name": "Budi Santoso"
    },
    "note": "Permintaan approval selesai - mohon review",
    "requested_at": "2026-01-29T12:00:00Z"
  }
  ```

  Response (Success - 200):
  ```json
  {
    "success": true,
    "message": "Onboarding submitted for approval",
    "data": {
      "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
      "previous_status": "DRAFT",
      "status": "SUBMITTED_FOR_APPROVAL",
      "requested_at": "2026-01-29T12:00:00Z",
      "requested_by": {
        "user_id": "user-uuid-12345",
        "email": "budi@email.com",
        "full_name": "Budi Santoso"
      },
      "action_audit_id": "action-uuid-req-123"
    }
  }
  ```

  Response (Validation Error - 400):
  ```json
  {
    "success": false,
    "message": "Onboarding cannot be submitted: missing required acceptance or incomplete data",
    "code": "VALIDATION_FAILED",
    "errors": [
      "Terms not accepted",
      "Payment setup incomplete"
    ]
  }
  ```

  Response (Not Allowed - 409):
  ```json
  {
    "success": false,
    "message": "Onboarding is not in DRAFT state",
    "code": "INVALID_STATUS_TRANSITION",
    "current_status": "SUBMITTED_FOR_APPROVAL"
  }
  ```

  Notes / Server-side Rules:
  - Sebelum mengubah status ke `SUBMITTED_FOR_APPROVAL`, backend harus memvalidasi bahwa semua required sections telah terpenuhi (termasuk terms accepted, minimal payment setup, wajib field business-entity).
  - Jika validasi gagal, kembalikan 400 dengan daftar kesalahan.
  - Jika onboarding bukan `DRAFT`, tolak permintaan dengan 409 `INVALID_STATUS_TRANSITION`.
  - Log action dengan `action_audit_id`, `requested_by`, dan `requested_at` untuk audit.

  ## GET CURRENT STATUS (UNTUK FRONTEND)

  ### GET /api/onboarding/{merchant_onboarding_id}/status

  Description: Mengambil status current onboarding beserta history singkat perubahan status (audit) — digunakan frontend untuk men-disable tombol submit ketika belum bisa diajukan.

  Response (Success - 200):
  ```json
  {
    "success": true,
    "data": {
      "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
      "current_status": "DRAFT",
      "allowed_actions": ["save_business_entity","save_payment_setup","accept_terms","request_approval"],
      "status_history": [
        {
          "from": "DRAFT",
          "to": "SUBMITTED_FOR_APPROVAL",
          "action": "request_approval",
          "by": {
            "user_id": "user-uuid-12345",
            "email": "budi@email.com"
          },
          "at": "2026-01-29T12:00:00Z",
        }
      ]
    }
  }
  ```

  Notes for Frontend:
  - Frontend must not attempt to change status directly. Gunakan `POST /request-approval` untuk submit.
  - Frontend should call `GET /status` to decide whether to enable the "Request Approval" button.

  ## Enforcement Summary
  - Do NOT create a generic `update-status` endpoint.
  - All data-saving endpoints MUST return data and keep status unchanged.
  - Only `POST /request-approval` may change status to `SUBMITTED_FOR_APPROVAL` (and future server-reviewed transitions like APPROVED/REJECTED must be handled by backend workflows).

  ---
- User status awal: `PENDING_EMAIL_VERIFICATION`
- Email verification link dikirim ke email user
- Token verification berlaku 24 jam (TTL)
- Field `referral_code` opsional

---

## 2️⃣ ENDPOINT: EMAIL VERIFICATION

### GET /api/auth/verify-email?token=xxxx

**Description:** User memverifikasi email melalui link yang dikirim

**Query Parameters:**
- `token` (required): Verification token dari email link

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Email verified successfully"
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "message": "Invalid or expired verification link"
}
```

**Notes:**
- Token bersifat sekali pakai (one-time use)
- Setelah berhasil, user status berubah menjadi: `ACTIVE`
- User bisa melakukan login setelah email terverifikasi
- Tidak menggunakan OTP, hanya link verification

---

## 3️⃣ ENDPOINT: RESEND VERIFICATION EMAIL (OPTIONAL)

### POST /api/auth/resend-verification

**Description:** Mengirim ulang email verifikasi jika user belum menerima atau link expired

**Request Payload:**
```json
{
  "email": "budi@email.com"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Verification email resent. Please check your inbox."
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "message": "Email not found or already verified"
}
```

---

## 4️⃣ ENDPOINT: LOGIN (SIGN IN)

### POST /api/auth/login

**Description:** User login dengan email dan password

**Request Payload:**
```json
{
  "email": "budi@email.com",
  "password": "securePassword123"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "user-uuid-12345",
      "full_name": "Budi Santoso",
      "business_name": "Toko Jaya Sejahtera",
      "email": "budi@email.com",
      "phone_number": "08123456789",
      "email_verified": true,
      "status": "ACTIVE"
    }
  }
}
```

**Response (Email belum verifikasi - 403):**
```json
{
  "success": false,
  "message": "Please verify your email before signing in",
  "code": "EMAIL_NOT_VERIFIED"
}
```

**Response (Invalid credential - 401):**
```json
{
  "success": false,
  "message": "Invalid email or password",
  "code": "INVALID_CREDENTIALS"
}
```

**Response (User not found - 404):**
```json
{
  "success": false,
  "message": "User not found",
  "code": "USER_NOT_FOUND"
}
```

---

## 5️⃣ ENDPOINT: GET CURRENT USER (ME)

### GET /api/auth/me

**Description:** Mendapatkan data user yang sedang login (untuk preload di dashboard)

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (Success - 200):**
```json
{
  "success": true,
  "data": {
    "id": "user-uuid-12345",
    "full_name": "Budi Santoso",
    "business_name": "Toko Jaya Sejahtera",
    "email": "budi@email.com",
    "phone_number": "08123456789",
    "email_verified": true,
    "status": "ACTIVE"
  }
}
```

**Response (Unauthorized - 401):**
```json
{
  "success": false,
  "message": "Unauthorized. Invalid or expired token",
  "code": "UNAUTHORIZED"
}
```

**Notes:**
- Endpoint ini digunakan untuk preload user data di dashboard
- Field yang dikembalikan: `full_name`, `business_name`, `email`, `phone_number`
- Diperlukan untuk menampilkan di dashboard setelah login

---

## 6️⃣ ENDPOINT: LOGOUT

### POST /api/auth/logout

**Description:** User logout dan invalidate access token

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## 7️⃣ ENDPOINT: REFRESH TOKEN (OPTIONAL)

### POST /api/auth/refresh-token

**Description:** Refresh access token menggunakan refresh token

**Request Payload:**
```json
{
  "refresh_token": "refresh_token_value"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "data": {
    "access_token": "new_jwt_access_token_here",
    "token_type": "Bearer",
    "expires_in": 3600
  }
}
```

---

### Token Storage
- **Option 1:** HttpOnly Cookie (Recommended untuk security)
  - Stored di: `Set-Cookie: access_token=...; HttpOnly; Secure; SameSite=Strict`
  - Frontend tidak bisa akses langsung (prevent XSS)
  - Automatically sent dengan setiap request
  
- **Option 2:** Authorization Header (Untuk API consumption)
  - Stored di: localStorage atau sessionStorage
  - Sent di header: `Authorization: Bearer <token>`
  - Frontend perlu handle manual

### Token Expiration & Refresh
- **Access Token:** 2 hour (3600 seconds)
- **Refresh Token:** 7 days (opsional)
- Backend validate token setiap request
- If token expired → return 401 → frontend redirect ke login

---

## 📊 DASHBOARD DATA MAPPING

Setelah user berhasil login, dashboard perlu mengambil data berikut dari endpoint `/api/auth/me`:

| Field | Source | Kegunaan |
|-------|--------|----------|
| `full_name` | User profile | Menampilkan greeting "Halo, Budi" |
| `business_name` | User profile | Header/title dashboard merchant |
| `email` | User profile | Untuk konfirmasi & reset password |
| `phone_number` | User profile | Untuk contact verification |

---

# 📋 ONBOARDING API CONTRACTS

## 1️⃣ Endpoint: Business Entity – Merchant

### POST /api/onboarding/business-entity/merchant

**Request Payload:**
```json
{
  "merchant_onboarding_id": "uuid",
  "business_type": "company | individual",
  "brand_name": "Toko Jaya",
  "legal_name": "PT Toko Jaya Sejahtera",
  "business_category_id": "mcc_id",
  "company_type": "PT",
  "established_year": 2020,
  "employee_count": 10,
  "business_mode": "offline",
  "business_ownership_status": "owned",
  "address": {
    "street": "...",
    "rt": 1,
    "rw": 2,
    "province_id": "ID-JB",
    "city_id": "ID-JB-BDG",
    "district_id": "ID-JB-BDG-01",
    "subdistrict_id": "ID-JB-BDG-01-01",
    "postal_code": 40123
  },
  "status": "DRAFT"
}
```

**Response:**
```json
{
  "success": true,
  "merchant_onboarding_id": "uuid"
}
```

## 2️⃣ Endpoint: Business Entity – Owner

### POST /api/onboarding/business-entity/owner

**Request Payload:**
```json
{
  "merchant_onboarding_id": "uuid",
  "owner_name": "Budi",
  "birth_date": "1990-01-01",
  "birth_place": "Bandung",
  "nationality": "WNI | WNA",
  "nik": "3273xxxxxxxx",
  "passport_number": null,
  "kitas_number": null,
  "address_ktp": { ... },
  "address_domicile": { ... },
  "is_domicile_same_as_ktp": true,
  "status": "DRAFT"
}
```

**Response:**
```json
{
  "success": true
}
```

## 3️⃣ Endpoint: Business Entity – PIC Admin

### POST /api/onboarding/business-entity/pic-admin

**Request Payload:**
```json
{
  "merchant_onboarding_id": "uuid",
  "pic_name": "Andi",
  "pic_phone": "08123xxxx",
  "pic_email": "andi@email.com",
  "status": "DRAFT"
}
```

## 4️⃣ Endpoint: Business Entity – Settlement

### POST /api/onboarding/business-entity/settlement

**Request Payload:**
```json
{
  "merchant_onboarding_id": "uuid",
  "bank_id": "BCA",
  "account_number": "1234567890",
  "account_name": "PT Toko Jaya Sejahtera",
  "status": "DRAFT"
}
```

## 5️⃣ File Upload

### POST /api/files/upload
**Content-Type:** multipart/form-data

**FormData:**
- file (image/pdf)
- category (e.g: "owner_ktp", "npwp", "business_license")

**Response:**
```json
{
  "file_id": "file_uuid",
  "file_url": "https://cdn.xxx/file.pdf",
  "mime_type": "application/pdf"
}
```

### GET /api/files/{file_id}

**Response:**
(binary stream or signed URL)

## 6️⃣ Get Data Saat Edit Merchant

### GET /api/onboarding/{merchant_onboarding_id}

**Response:**
```json
{
  "merchant": { ... },
  "owner": { ... },
  "pic_admin": { ... },
  "settlement": { ... },
  "documents": [ ... ],
  "status": "DRAFT"
}
```

---

# 💳 PAYMENT FEATURE API CONTRACTS

## Payment Feature Overview

Payment feature adalah fitur pembayaran yang diaktifkan merchant berdasarkan device yang digunakan. Setiap device memiliki payment features yang dapat dipilih.

### Device & Payment Feature Mapping

| Device | Payment Features |
|--------|------------------|
| **EDC** | CREDIT_DEBIT_CARD_PRESENT, QRIS, VA, BNPL |
| **SoftPOS** | CREDIT_DEBIT_CARD_PRESENT |
| **Soundbox QR Static** | QRIS |
| **QR Static (Sticker)** | QRIS |
| **Payment Link** (device=null) | QRIS, VA, CNP |

### Payment Feature Description

- **CREDIT_DEBIT_CARD_PRESENT**: Pembayaran kartu kredit/debit dengan physical card present
- **QRIS**: Quick Response Code Indonesian Standard - pembayaran via scan QR
- **VA**: Virtual Account - pembayaran transfer bank ke nomor virtual
- **BNPL**: Buy Now Pay Later - cicilan tanpa bunga
- **CNP**: Card Not Present - pembayaran kartu tanpa physical card

---

## 1️⃣ ENDPOINT: SAVE PAYMENT SETUP (DRAFT)

### POST /api/onboarding/payment-setup

**Description:** Menyimpan konfigurasi payment feature dalam status DRAFT

**Request Payload:**
```json
{
  "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
  "payment_setups": [
    {
      "device": "edc",
      "payment_features": [
        "CREDIT_DEBIT_CARD_PRESENT",
        "QRIS",
        "VA"
      ] 
    },
    {
      "device": "softpos",
      "payment_features": [
        "CREDIT_DEBIT_CARD_PRESENT"
      ]
    }
  ],
  "status": "DRAFT"
}
```

**Response (Success - 200/201):**
```json
{
  "success": true,
  "message": "Payment setup saved successfully",
  "data": {
    "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
    "payment_setups": [
      {
        "id": "payment-setup-uuid-1",
        "device": "edc",
        "payment_features": [
          "CREDIT_DEBIT_CARD_PRESENT",
          "QRIS",
          "VA"
        ],
        "status": "DRAFT"
      },
      {
        "id": "payment-setup-uuid-2",
        "device": "softpos",
        "payment_features": [
          "CREDIT_DEBIT_CARD_PRESENT"
        ],
        "status": "DRAFT"
      }
    ],
    "created_at": "2026-01-29T10:30:00Z",
    "updated_at": "2026-01-29T10:30:00Z"
  }
}
```

**Response (Error - 400):**
```json
{
  "success": false,
  "message": "Invalid payment setup configuration",
  "errors": [
    {
      "device": "edc",
      "error": "payment_features contains invalid feature: INVALID_FEATURE. Allowed: CREDIT_DEBIT_CARD_PRESENT, QRIS, VA, BNPL"
    }
  ]
}
```

**Response (Error - 404):**
```json
{
  "success": false,
  "message": "Merchant onboarding not found"
}
```

**Validation Rules:**
- `merchant_onboarding_id` wajib ada
- `payment_setups` minimum 1 setup
- Setiap setup harus memiliki `device` dan `payment_features`
- Payment features yang dipilih HARUS sesuai dengan device mapping
- Device boleh null untuk Payment Link
- Tidak boleh ada duplicate device dalam satu request
- Payment features tidak boleh kosong

---

## 2️⃣ ENDPOINT: GET PAYMENT SETUP (EDIT)

### GET /api/onboarding/{merchant_onboarding_id}/payment-setup

**Description:** Mengambil konfigurasi payment feature yang tersimpan (untuk edit)

**Path Parameters:**
- `merchant_onboarding_id` (required): UUID merchant onboarding

**Response (Success - 200):**
```json
{
  "success": true,
  "data": {
    "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
    "payment_setups": [
      {
        "id": "payment-setup-uuid-1",
        "device": "edc",
        "payment_features": [
          "CREDIT_DEBIT_CARD_PRESENT",
          "QRIS",
          "VA"
        ],
        "status": "DRAFT"
      },
      {
        "id": "payment-setup-uuid-2",
        "device": "softpos",
        "payment_features": [
          "CREDIT_DEBIT_CARD_PRESENT"
        ],
        "status": "DRAFT"
      },
      {
        "id": "payment-setup-uuid-3",
        "device": null,
        "payment_features": [
          "QRIS",
          "VA",
          "CNP"
        ],
        "status": "DRAFT"
      }
    ],
    "created_at": "2026-01-29T10:30:00Z",
    "updated_at": "2026-01-29T10:45:00Z"
  }
}
```

**Response (Not Found - 404):**
```json
{
  "success": false,
  "message": "Payment setup not found for this merchant"
}
```

---

## 3️⃣ EXAMPLE SCENARIOS

### Scenario 1: EDC Only

```json
{
  "merchant_onboarding_id": "uuid",
  "payment_setups": [
    {
      "device": "edc",
      "payment_features": [
        "CREDIT_DEBIT_CARD_PRESENT",
        "QRIS",
        "VA",
        "BNPL"
      ]
    }
  ],
  "status": "DRAFT"
}
```

### Scenario 2: SoftPOS + Payment Link

```json
{
  "merchant_onboarding_id": "uuid",
  "payment_setups": [
    {
      "device": "softpos",
      "payment_features": [
        "CREDIT_DEBIT_CARD_PRESENT"
      ]
    },
    {
      "device": null,
      "payment_features": [
        "QRIS",
        "VA"
      ]
    }
  ],
  "status": "DRAFT"
}
```

### Scenario 3: Multi Device Setup

```json
{
  "merchant_onboarding_id": "uuid",
  "payment_setups": [
    {
      "device": "edc",
      "payment_features": [
        "CREDIT_DEBIT_CARD_PRESENT",
        "QRIS"
      ]
    },
    {
      "device": "soundbox",
      "payment_features": [
        "QRIS"
      ]
    },
    {
      "device": "sticker",
      "payment_features": [
        "QRIS"
      ]
    }
  ],
  "status": "DRAFT"
}
```

---

# **Terms & Conditions API CONTRACT**

## Terms Overview

- Terms harus dibaca dan dichecklist oleh user sebelum melanjutkan proses onboarding.
- Acceptance adalah persetujuan hukum dan harus tercatat sebagai legal record yang auditable.
- Terms bersifat versioned; setiap perubahan membuat versi baru.

## 1️⃣ ENDPOINT: ACCEPT TERMS

### POST /api/onboarding/terms/accept

Description: User menyetujui Terms & Conditions yang disediakan backend. Acceptance disimpan sebagai legal record.

Request Payload:
```json
{
  "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
  "accepted_at": "2026-01-29T11:00:00Z",
  "accepted_by": {
    "user_id": "user-uuid-12345",
    "email": "budi@email.com",
    "full_name": "Budi Santoso"
  }
}
```

Response (Success - 201):
```json
{
  "success": true,
  "message": "Terms accepted",
  "data": {
    "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
    "accepted_at": "2026-01-29T11:00:00Z",
    "accepted_by": {
      "user_id": "user-uuid-12345",
      "email": "budi@email.com",
      "full_name": "Budi Santoso"
    },
    "created_at": "2026-01-29T11:00:00Z"
  }
}
```

Response (Error - 400):
```json
{
  "success": false,
  "message": "Invalid payload or terms_version not found"
}
```

Notes:
- `merchant_onboarding_id` wajib.
- `accepted_at` harus berupa ISO8601 timestamp.
- Acceptance adalah immutable legal record (backend should append history, do not overwrite previous acceptance records).

## 2️⃣ ENDPOINT: GET TERMS STATUS (UNTUK EDIT)

### GET /api/onboarding/{merchant_onboarding_id}/terms

Description: Mengambil status acceptance dan detail terms terbaru/yang diterima untuk keperluan edit atau review.

Path Parameters:
- `merchant_onboarding_id` (required)

Response (Success - 200):
```json
{
  "success": true,
  "data": {
    "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
    "accepted": {
      "is_accepted": true,
      "terms_version": "v1.2.0",
      "accepted_at": "2026-01-29T11:00:00Z",
      "accepted_by": {
        "user_id": "user-uuid-12345",
        "email": "budi@email.com",
        "full_name": "Budi Santoso"
      },
    }
  }
}
```

Response (Not Accepted - 200):
```json
{
  "success": true,
  "data": {
    "merchant_onboarding_id": "550e8400-e29b-41d4-a716-446655440000",
    "accepted": {
      "is_accepted": false,
      "accepted_at": null
    }
  }
}
```

Response (Not Found - 404):
```json
{
  "success": false,
  "message": "Merchant onboarding not found"
}
```

## Validation & Business Rules

- Onboarding **TIDAK BISA** di-submit final (see `/api/onboarding/{merchant_onboarding_id}/submit`) jika tidak ada acceptance record untuk current terms version. Backend harus menolak submit dengan error:
```json
{
  "success": false,
  "message": "Terms not accepted",
  "code": "TERMS_NOT_ACCEPTED"
}
```
- Terms harus versioned. Backend menyimpan daftar versi terms (dokumen) beserta URL atau content-hash.
- Acceptance harus auditable: simpan `accepted_at`, `accepted_by` (user_id, email), `user_agent`.

---

## 7️⃣ Master Data Endpoint

### MCC
### GET /api/master/mcc

**Response:**
```json
[
  { "id": "5812", "name": "Restaurant" }
]
```

### Wilayah
### GET /api/master/provinces
### GET /api/master/cities?province_id=ID-JB
### GET /api/master/districts?city_id=ID-JB-BDG
### GET /api/master/subdistricts?district_id=ID-JB-BDG-01

**Response:**
```json
[
  { "id": "...", "name": "..." }
]
```

### Bank
### GET /api/master/banks

**Response:**
```json
[
  { "id": "BCA", "name": "Bank Central Asia" }
]
```

## 8️⃣ Submit Final Onboarding

### POST /api/onboarding/{merchant_onboarding_id}/submit

**Response:**
```json
{
  "success": true,
  "status": "SUBMITTED"
}
```

