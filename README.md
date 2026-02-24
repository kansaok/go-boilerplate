myapp/
│
├── cmd/ # Berisi entry point aplikasi
│ └── root.go
│ └── migrate.go
│ └── seed.go
│
├── pkg/ # Berisi package reusable yang bisa digunakan ulang oleh project lain
│ |── logger/ # Custom logger
│ | └── logger.go
│ └── telemetry/ # OpenTelemetry dan Prometheus setup
│ ├── otel.go # Inisialisasi OpenTelemetry
│ └── prometheus.go # Setup middleware Prometheus
│
├── internal/ # Berisi kode yang hanya digunakan secara internal
│ ├── config/ # Konfigurasi aplikasi
│ │ └── config.go
│ │ └── storage
│ │ └── local_storage.go
│ │ └── s3_storage.go
│ ├── db/ # Koneksi database dan migrasi
│ │ └── db.go
│ │ └── migrate.go
│ │ └── seeder.go
│ │ └── queries/
│ │ └── schema/
│ │ └── migrations/
│ │ └── seeders/
│ │ ├── user_seeder.go
│ │ └── product_seeder.go
| |
│ ├── middleware/ # Middleware (Autentikasi, Logging, dll)
│ │ └── auth.go
│ │ └── logging.go
│ ├── service/ # Business logic utama (terpisah dari HTTP handling)
│ │ ├── auth.service.go
│ │ ├── user.service.go
│ │ └── product.service.go
│ ├── repository/ # Data access layer (untuk query DB)
│ │ ├── auth.repository.go
│ │ ├── user.repository.go
│ │ └── product.repository.go
│ |── controller/ # Controller untuk mengelola logika request
│ | ├── auth.controller.go # Kontroler untuk auth seperti register, login, refresh token, logout
│ | ├── user.controller.go # Kontroler untuk user
│ | └── product.controller.go # Kontroler untuk produk
│ ├── model/ # Struct dan model untuk database dan response/request
│ │ ├── user.model.go
│ │ └── product.model.go
│ ├── routes/ # Routing terpisah
│ │ └── routes.go # gabungan route
│ │ └── auth.route.go # Definisi rute auth
│ │ └── user.route.go # Definisi rute user
│ │ └── product.route.go # nDefinisi rute product
│ └── util/
│ └── helper.go # Utility/helper functions
│ └── response.go
│ └── validation.go
|
|── storage/ # untuk menyimpan database sqlite
| └── sqlite/
|
|── uploads/ # untuk menyimpan data upload
|
├── web/ # Berisi static file (HTML, CSS, JS jika dibutuhkan)
│ └── assets/
│ └── css/
│ └── js/
│
│── main.go # File utama untuk memulai aplikasi
├── .env # Environment variables (optional, jika menggunakan)
├── .env.example # Environment variables (optional, jika menggunakan)
├── go.mod # Go module file
├── go.sum # Dependensi dari module
└── README.md
└── CHANGELOG.md
└── sqlc.yaml
