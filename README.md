# Bank Sampah Enterprise Ecosystem

Arsitektur aplikasi pengelolaan bank sampah modern berbasis kombinasi:
1. **Pola Utama: Modular Monolith Architecture (`services/core-go`)**
   - **Single Deployable Unit**: Seluruh sistem backend dikompilasi menjadi satu binary executable mandiri (`core_server.exe`), RAM < 50 MB pada beban standar.
   - **Domain Encapsulation**: Kode diorganisasikan per domain bisnis terpisah pada direktori `app/modules/` (`auth`, `user`, `pickup`, `weighbridge`, `tracking`, `wallet`, `reward`, `csr`).
2. **Pola Internal Modul: Clean Layered Architecture (5-Tier)**
   - `route.go` & `handler.go`: Presentation Layer (Routing Echo, Binding, JSON Response).
   - `service.go`: Business Logic Layer (Kalkulasi saldo, validasi aturan, callback async).
   - `repository.go`: Data Access Layer (GORM DB Queries, Pagination, Transaksi).
   - `model.go` & `redis.go`: Entity Layer (GORM Struct, Schema Mapping) & Redis Caching/Geospatial.
3. **AI Microservice (`services/ai-fastapi`)**:
   - Klasifikasi visual sampah (ONNX YOLO model inference).
   - Estimasi bobot, nilai ekonomi (Rupiah), dan jejak karbon (CO2 reduction).
   - Optimasi rute penjemputan sampah kurir dengan **Google OR-Tools Capacitated Vehicle Routing Problem (CVRP)**.
4. **IoT Weighbridge (`firmware/weighbridge-esp32`)**:
   - Pembacaan sensor timbangan digital (HX711) dan auto-submit tiket timbangan ke server Go.
5. **Apps & Contracts**:
   - `apps/web-admin`: Dashboard Next.js untuk Admin, Operator Timbangan, dan Laporan ESG Mitra CSR.
   - `apps/mobile`: Aplikasi Flutter untuk Nasabah dan Kurir.
   - `docs/api-contract/openapi.yaml`: OpenAPI 3.0 specification.

---

## Struktur Direktori

```text
bank-sampah/
├── docker-compose.yml          # Postgres 16, Redis 7, RabbitMQ, 2 Services
├── Makefile                    # Target: up, down, build, run-go, run-ai, test
├── .env.example
├── .github/workflows/          # CI Pipeline dengan path-filter per service
│
├── bin/
│   └── core_server.exe         # Single compiled executable binary
│
├── services/
│   ├── core-go/                # Go Modular Monolith (Echo + GORM + Redis)
│   │   ├── cmd/server/main.go
│   │   ├── app/modules/
│   │   │   ├── auth/           # route.go handler.go service.go repository.go model.go redis.go
│   │   │   ├── user/           # route.go handler.go service.go repository.go model.go
│   │   │   ├── pickup/         # route.go handler.go service.go repository.go model.go
│   │   │   ├── weighbridge/    # route.go handler.go service.go repository.go model.go (Auto-credit wallet)
│   │   │   ├── tracking/       # route.go handler.go service.go repository.go model.go redis.go (GEOADD/SEARCH)
│   │   │   ├── wallet/         # route.go handler.go service.go repository.go model.go (Ledger, Saldo, Redeem)
│   │   │   ├── reward/         # route.go handler.go service.go repository.go model.go (Points exchange)
│   │   │   └── csr/            # route.go handler.go service.go repository.go model.go (ESG metrics)
│   │   ├── common/             # Middleware (JWT, RBAC, Audit), response helper, error types
│   │   ├── config/             # DB & Redis connection pool
│   │   ├── database/migrations/
│   │   ├── go.mod
│   │   └── Dockerfile
│   │
│   └── ai-fastapi/             # AI & Optimization Microservice (FastAPI)
│       ├── cmd/api/main.py
│       ├── app/modules/
│       │   ├── classify/       # handler.py schema.py service.py (ONNX visual classifier)
│       │   ├── estimate/       # handler.py schema.py service.py (Density & price estimator)
│       │   └── routing/        # handler.py schema.py service.py (Google OR-Tools CVRP)
│       ├── common/             # Audit middleware, error handler
│       ├── config/settings.py
│       ├── requirements.txt
│       └── Dockerfile
│
├── apps/
│   ├── mobile/                 # Flutter mobile app
│   └── web-admin/              # Next.js admin & CSR portal
│
├── firmware/weighbridge-esp32/ # ESP32 Arduino code for HX711 scale
└── docs/api-contract/          # OpenAPI 3.0 specification
```

---

## Cara Menjalankan

### 1. Menggunakan Docker Compose (Rekomendasi)
```bash
# Salin konfigurasi environment
cp .env.example .env

# Jalankan seluruh stack (Database, Cache, Queue, dan 2 Service)
docker-compose up -d --build
```

### 2. Menjalankan Core Go Secara Lokal
```bash
cd services/core-go
go run ./cmd/server/main.go
# Atau kompilasi menjadi single binary:
go build -o ../../bin/core_server.exe ./cmd/server/main.go
```

### 3. Menjalankan AI FastAPI Secara Lokal
```bash
cd services/ai-fastapi
pip install -r requirements.txt
uvicorn cmd.api.main:app --reload --port 8000
```
