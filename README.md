# Go Boilerplate - API Application

Aplikasi Go Boilerplate dengan clean architecture, JWT authentication, dan middleware keamanan lengkap.

## Fitur

- JWT Authentication (access token; secret key wajib minimal 32 karakter)
- Account lockout: 5 percobaan login gagal = lock 15 menit (per-IP di middleware, dan per-akun via `AccountLimiter`)
- Rate limiting: 100 req/menit global (per-IP), 10 req/15menit khusus endpoint auth
- Per-account login/registration throttle, dengan Redis (multi-instance) atau in-memory fallback (single-instance)
- Body size limit per request (`MAX_BODY_BYTES`, default 16MB; endpoint auth dibatasi 1MB)
- Security headers (CSP, HSTS, X-Frame-Options, Referrer-Policy, X-Content-Type-Options, dll)
- CSRF token & session cookie management (untuk jalur berbasis cookie; jalur utama tetap JWT via header)
- Host header validation & SSL redirect
- Endpoint `/metrics` (Prometheus) terproteksi: HTTP Basic Auth atau IP/CIDR allowlist, default hanya loopback
- Database Migrations & Seeders (CLI)
- Logging (rotasi file via lumberjack) & Telemetry (OpenTelemetry tracing + Prometheus metrics)
- File Upload dengan validasi ekstensi + MIME sniffing + ukuran (Local storage & AWS S3)

> Catatan: anotasi Swagger sudah ada di beberapa handler, tapi generator/serving Swagger UI
> (`swaggo`) belum di-wire ke project ini — folder `docs/` masih kosong.

## Struktur Folder

```
backend/
├── cmd/              # CLI commands (migrate, seeder)
├── pkg/              # Reusable packages (logger, telemetry)
├── internal/         # Internal application code
│   ├── config/       # Configuration (db, jwt, cors, security, redis)
│   │   └── storage/  # Local & S3 file storage
│   ├── db/           # Database connection & migrations
│   ├── middleware/   # Auth, logging, rate limiting, body limit, security headers
│   ├── modules/      # Feature modules (auth, user)
│   ├── routes/       # Route definitions
│   └── util/         # Utilities (validators, helpers, uploaders)
├── database/
│   ├── migrations/   # SQL migration files
│   └── seeders/      # Database seeders
├── docs/             # (placeholder) Swagger output, belum digenerate
├── storage/          # SQLite database (jika DB_CONNECTION=sqlite)
└── uploads/          # Uploaded files (local storage)
```

## Setup & Running

### Prerequisites

- Go 1.27+
- Database (PostgreSQL/MySQL/SQLite/MongoDB)
- Redis (opsional, untuk rate limiting lintas instance)
- Environment variables (lihat `.env.example` dan `.env.database.example`)

### Local Development

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   cp .env.database.example .env.database
   # Edit .env dan .env.database sesuai konfigurasi kamu
   # JWT_SECRET_KEY wajib diisi (min. 32 karakter), contoh: openssl rand -hex 32
   ```

4. **Run migrations**
   ```bash
   go run main.go migrate
   ```

5. **Run seeders (optional)**
   ```bash
   go run main.go db:seed
   ```

6. **Run the server**
   ```bash
   go run main.go
   ```
   Server berjalan di `http://localhost:8080`

### Docker Deployment

1. **Setup environment**
   ```bash
   cp .env.example .env
   cp .env.database.example .env.database
   # REDIS_PASSWORD WAJIB diisi - container redis akan gagal start tanpa ini
   ```

2. **Run dengan docker-compose (local)**
   ```bash
   docker-compose up -d
   ```

3. **Run dengan docker-compose (production)**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```
   Di production, port database & Redis tidak dipublish ke host — hanya reachable
   lewat jaringan internal Docker (`app-network`).

### VPS Deployment

1. **Setup VPS**
   ```bash
   apt update && apt upgrade -y
   apt install -y postgresql postgresql-contrib
   ```

2. **Setup application**
   ```bash
   mkdir -p /opt/go-boilerplate
   cd /opt/go-boilerplate
   # Copy project files
   cp .env.example .env
   cp .env.database.example .env.database
   # Edit dengan value production
   mkdir -p storage/sqlite uploads
   ```

3. **Run migrations & start**
   ```bash
   go run main.go migrate
   go run main.go
   ```

## CLI Commands

```bash
# Run migrations
go run main.go migrate

# Rollback last migration
go run main.go migrate:rollback

# Show migration status
go run main.go migrate:status

# Fresh migrate (drop all tables)
go run main.go migrate:fresh

# Create new migration
go run main.go make:migration migration_name

# Run seeders
go run main.go db:seed

# Run specific seeder
go run main.go db:seed --func UsersSeeder

# Create new seeder
go run main.go make:seeder seeder_name
```

## API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register user baru | No |
| POST | `/api/v1/auth/login` | Login dengan email/password | No |
| GET | `/api/v1/me` | Ambil identitas user dari JWT | Ya (Bearer token) |
| GET | `/metrics` | Prometheus metrics | Basic Auth / IP allowlist |

## Security Features

- JWT authentication (HS256, algoritma dipaksa, secret key wajib ≥32 karakter)
- Timing-safe login: dummy bcrypt compare saat email tidak ditemukan (anti user enumeration)
- Rate limiting global & khusus auth, plus lockout per-IP dan per-akun (in-memory / Redis)
- Body size limit untuk mencegah payload raksasa
- Security headers (CSP, HSTS, XSS filter, nosniff, X-Frame-Options)
- CSRF token generation untuk jalur berbasis cookie
- SQL injection prevention (GORM parameterized query di semua modul)
- Password hashing dengan bcrypt (cost 12)
- File upload divalidasi ekstensi + MIME sniffing isi file + limit ukuran, nama file di-randomize
- CORS configuration (`AllowCredentials: false` by default)
- Host header validation & SSL redirect
- Endpoint `/metrics` terproteksi Basic Auth atau IP/CIDR allowlist (default: hanya loopback)
- `SetTrustedProxies(nil)` — mencegah `X-Forwarded-For` dipalsukan untuk membypass rate limit

## Docker Compose

### Local Development

```bash
# Start services (app + database + redis)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Production

```bash
# Siapkan .env dan .env.database dengan value production
# (REDIS_PASSWORD wajib diisi)
docker-compose -f docker-compose.prod.yml up -d
```

## Environment Variables Reference

| Variable | Description | Required |
|----------|-------------|----------|
| `JWT_SECRET_KEY` | Secret key untuk signing JWT (min. 32 karakter) | Yes |
| `ACCESS_TOKEN_LIFETIME` | Durasi access token (contoh: `15m`) | No |
| `REFRESH_TOKEN_LIFETIME` | Durasi refresh token (contoh: `168h`), disiapkan untuk fitur mendatang | No |
| `DB_CONNECTION` | Tipe database (`postgres`/`mysql`/`sqlite`/`mongodb`) | Yes |
| `DB_HOST` / `DB_PORT` / `DB_USER` / `DB_PASSWORD` / `DB_NAME` | Kredensial database SQL/MongoDB | Yes (kecuali sqlite) |
| `DB_FILE` | Path file database (khusus `sqlite`) | Yes (jika sqlite) |
| `DB_SSLMODE` | Mode SSL PostgreSQL | No |
| `S3_BUCKET_NAME` | Nama bucket S3 untuk upload file | No |
| `ALLOWED_HOSTS` | Host yang diizinkan (comma-separated) | No (default `localhost`) |
| `ALLOWED_ORIGINS` | Origin CORS yang diizinkan (comma-separated) | No |
| `SECURE_SSL_REDIRECT` | Redirect HTTP ke HTTPS | No |
| `CSRF_COOKIE_SECURE` / `SESSION_COOKIE_SECURE` | Flag `Secure` pada cookie CSRF/session | No |
| `SECURE_BROWSER_XSS_FILTER` / `SECURE_CONTENT_TYPE_NOSNIFF` | Toggle header `X-XSS-Protection` / `X-Content-Type-Options` | No |
| `MAX_BODY_BYTES` | Batas ukuran body request (bytes), default 16MB | No |
| `METRICS_USER` / `METRICS_PASSWORD` | Basic Auth untuk endpoint `/metrics` | No |
| `METRICS_ALLOWLIST` | CIDR yang boleh akses `/metrics` tanpa Basic Auth | No |
| `REDIS_HOST` | Host Redis; kosongkan untuk in-memory fallback | No |
| `REDIS_PORT` / `REDIS_PASSWORD` / `REDIS_DB` | Konfigurasi koneksi Redis | Ya jika `REDIS_HOST` diisi / pakai docker-compose |
| `GIN_MODE` | Mode Gin (`debug`/`release`/`test`) | No |
| `APP_PORT` | Port HTTP server | No (default `8080`) |

## License

MIT
