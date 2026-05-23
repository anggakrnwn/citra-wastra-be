
# Citra Wastra Backend
Backend service for Citra Wastra, a batik education & classification platform built with Go.

## Features
- Authentication with JWT
- Batik classification
- Media upload
- Learning & quiz system
- Gamification (XP & rewards)
- Redis caching
- Background worker

## Tech Stack
- Go
- PostgreSQL
- Redis
- Docker

## Project Structure
```
config/        # app configuration
dto/           # request/response DTO
handler/       # HTTP handlers
middleware/    # auth & middleware
models/        # database models
repository/    # database access
routes/        # API routes
service/       # business logic
utils/         # helpers
worker/        # background jobs
```

## Run Locally

Clone repository:
```bash
git clone https://github.com/anggakrnwn/citra-wastra-be.git
cd citra-wastra-be
```

Setup environment:
```bash
cp .env.example .env
```

Run dependencies:
```bash
docker compose up -d
```

Run application:
```bash
go run main.go
```

## Test
```bash
go test ./...
```
