# Go Boilerplate - API Application

Aplikasi Go Boilerplate dengan clean architecture, JWT authentication, dan middleware keamanan lengkap.

## Fitur

- JWT Authentication (Access Token + Refresh Token)
- Role-based access control
- Rate limiting (100 req/min global, 10 req/15min untuk auth)
- Account lockout (5 failed attempts = 15 min lock)
- Security headers (CSP, HSTS, Referrer-Policy, dll)
- CSRF Token & Session Management
- Database Migrations & Seeders
- Swagger API Documentation
- Logging & Telemetry (OpenTelemetry + Prometheus)
- File Upload (Local & AWS S3)

## Struktur Folder

```
myapp/
├── cmd/              # CLI commands (migrate, seeder)
├── pkg/              # Reusable packages (logger, telemetry)
├── internal/         # Internal application code
│   ├── config/       # Configuration
│   ├── db/           # Database connection & migrations
│   ├── middleware/   # Auth, logging, rate limiting, security
│   ├── modules/      # Feature modules (auth, user, product)
│   ├── routes/       # Route definitions
│   └── util/         # Utilities (validators, helpers, uploaders)
├── database/
│   ├── migrations/   # SQL migration files
│   └── seeders/      # Database seeders
├── docs/             # Swagger documentation
├── storage/          # SQLite database
├── uploads/          # Uploaded files
└── web/              # Static files (HTML, CSS, JS)
```

## Setup & Running

### Prerequisites

- Go 1.22+
- Database (PostgreSQL/MySQL/SQLite)
- Environment variables

### Local Development

1. **Clone repository**
   ```bash
   git clone <repository-url>
   cd myapp
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp .env.example .env
   cp .env.database.example .env.database
   # Edit .env and .env.database with your configuration
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
   Server will run on `http://localhost:8080`

7. **Access Swagger UI**
   Open `http://localhost:8080/docs/index.html`

### Docker Deployment

1. **Build image**
   ```bash
   docker build -t go-boilerplate .
   ```

2. **Run with docker-compose (local)**
   ```bash
   docker-compose up -d
   ```

3. **Run with docker-compose (production)**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

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
   # Edit with your values
   mkdir -p storage/sqlite uploads
   ```

3. **Run migrations**
   ```bash
   go run main.go migrate
   ```

4. **Build & run**
   ```bash
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
| POST | `/v1/auth/register` | Register new user | No |
| POST | `/v1/auth/login` | Login with email/password | No |
| GET | `/metrics` | Prometheus metrics | No |
| GET | `/docs/*` | Swagger UI | No |

## Security Features

- JWT authentication with configurable token lifetime
- Rate limiting to prevent abuse
- Account lockout after failed login attempts
- Security headers (CSP, HSTS, XSS filter)
- CSRF token generation
- SQL injection prevention
- Password hashing with bcrypt
- CORS configuration
- Host validation

## Docker Compose

### Local Development

```bash
# Start services (app + database)
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
# Create .env and .env.database files with production values
# Then start production
docker-compose -f docker-compose.prod.yml up -d
```

## Environment Variables Reference

| Variable | Description | Required |
|----------|-------------|----------|
| `JWT_SECRET_KEY` | Secret key for JWT signing | Yes |
| `ACCESS_TOKEN_LIFETIME` | Access token duration (e.g., 15m) | No |
| `REFRESH_TOKEN_LIFETIME` | Refresh token duration (e.g., 168h) | No |
| `DB_CONNECTION` | Database type (postgres/mysql/sqlite/mongodb) | Yes |
| `DB_HOST` | Database host | Yes |
| `DB_PORT` | Database port | Yes |
| `DB_USER` | Database username | Yes |
| `DB_PASSWORD` | Database password | Yes |
| `DB_NAME` | Database name | Yes |
| `DB_SSLMODE` | PostgreSQL SSL mode | No |
| `ALLOWED_HOSTS` | Allowed hosts | No |
| `SECURE_SSL_REDIRECT` | Redirect HTTP to HTTPS | No |

## License

MIT
