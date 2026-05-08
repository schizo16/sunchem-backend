# Sunchem Backend – Go/Gin REST API

Backend API cho **Sunchem Co., LTD** — hệ thống quản lý sản phẩm, blog, media và analytics.

---

## 1. Cài đặt & Chạy

### Yêu cầu
- **Go** >= 1.25
- **Docker** (tùy chọn)

### Chạy trực tiếp (dev)

```bash
cd sunchem-backend
go mod tidy
go run main.go
# API chạy tại http://localhost:8080
```

### Chạy bằng Docker

```bash
docker compose up -d --build
```

Chi tiết xem comment trong `docker-compose.yml`.

---

## 2. Cấu trúc dự án

```
sunchem-backend/
├── main.go              # Entry point
├── cmd/app/             # App bootstrap + seed data
├── internal/
│   ├── common/          # DB, middleware, migrations
│   └── modules/         # Domain modules (products, blog, auth, ...)
├── migrations/          # SQL migration files
├── Dockerfile
├── docker-compose.yml   # Orchestrate ca backend + frontend
├── go.mod
└── go.sum
```

## 3. API Endpoints

| Method | Path | Mô tả |
|--------|------|-------|
| GET | `/api/products` | Danh sách sản phẩm |
| GET | `/api/products/:slug` | Chi tiết sản phẩm |
| GET | `/api/blog` | Danh sách bài viết |
| GET | `/api/blog/:slug` | Chi tiết bài viết |
| POST | `/api/auth/login` | Đăng nhập |
| POST | `/api/contact` | Gửi liên hệ |

---

## 4. Env Variables

| Biến | Mặc định | Mô tả |
|------|----------|-------|
| `APP_ENV` | `dev` | Môi trường (dev/prod) |
| `DB_TYPE` | `sqlite` | Loại database |
| `DB_DSN` | `/app/data/sunchem.db` | Đường dẫn DB |
| `SERVER_PORT` | `8080` | Port server |
| `JWT_SECRET` | `...` | Secret key JWT |
| `UPLOAD_DIR` | `/app/uploads` | Thư mục upload |

---

## 5. Deploy với Docker Compose

Cấu trúc thư mục yêu cầu khi deploy:

```
project/
├── sunchem-backend/     ← Repo này (chứa docker-compose.yml)
│   └── docker-compose.yml
└── sunchem-frontend/    ← Clone repo frontend
    └── Dockerfile.prod
```

### Các lệnh Docker

```bash
# Build & chạy lần đầu
docker compose up -d --build

# Chạy lại bình thường
docker compose up -d

# Xóa hết (bao gồm volume DB) để seed lại
docker compose down -v
docker compose up -d --build

# Xem log
docker compose logs -f backend
docker compose logs -f frontend

# Dừng
docker compose down
```

### Ports sau khi chạy

| Service | Port | URL |
|---------|------|-----|
| Frontend (nginx) | 80 | http://localhost |
| Backend (Gin) | 8080 | http://localhost:8080/api |

---

## 6. Seed Data

Dữ liệu mẫu tự động được tạo khi chạy lần đầu (DB trống). Gồm:
- 17 sản phẩm hóa chất (SUNPERSE, TiO2, CABOT, màu paste...)
- 4 bài viết blog

Để seed lại: `docker compose down -v` rồi `docker compose up -d --build`.
