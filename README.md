# URL Shortener

A high-performance URL shortener service built with Go, Chi router, and SQLite. Features built-in Basic Authentication for protected endpoints and a full CI/CD deployment pipeline via GitHub Actions.

## Features

- **Standard Endpoints**: Create, Redirect, and Delete short URLs.
- **RESTful JSON API**: Validated inputs, typed JSON responses.
- **Basic Auth**: Private endpoints (Save, Delete) are protected using Basic Authentication.
- **SQLite Database**: Lightweight, file-based persistence.
- **Structured Logging**: Leveraging Go's `slog` package.
- **CI/CD Pipeline**: Automated deployments to an Ubuntu VM (Oracle Cloud) using GitHub Actions, SSH/rsync, and `systemd`.

## Project Structure

```text
cmd/url-shortener/       — application entry point
config/                  — YAML configuration files (local, dev, prod)
deployment/              — systemd service file
internal/
  config/                — config loader (YAML + environment vars)
  storage/sqlite/        — SQLite storage layer
  http-server/
    handlers/url/        — API handlers (save, redirect, delete)
    middleware/logger/   — custom request logger
  lib/
    api/response/        — shared JSON response helpers
    random/              — random alias generator
    logger/              — logger utilities
.github/workflows/       — GitHub Actions CI/CD workflows
```

## Setup & Local Development

1. **Clone the repository**

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Create the storage directory**
   ```bash
   mkdir -p storage
   ```

4. **Configure environment**
   Create a `.env` file in the project root:
   ```env
   CONFIG_PATH=./config/local.yaml
   ```

5. **Run the server**
   ```bash
   go run cmd/url-shortener/main.go
   ```

## Configuration

Configuration is managed via `cleanenv`, combining YAML files and environment variables.

Example `config/local.yaml`:
```yaml
env: "local"
storage_path: "./storage/storage.db"
http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 30s
  user: "admin"
  password: "your_local_password"
```

*Note: In production (`config/prod.yaml`), passwords can be supplied via the `HTTP_SERVER_PASSWORD` environment variable to keep secrets out of source control.*

## API

### 1. Redirect to Original URL (Public)
```http
GET /get/{alias}
```
Redirects (302) to the original URL.

### 2. Save URL (Protected)
```http
POST /url/save
Authorization: Basic <base64(user:pass)>
Content-Type: application/json

{
  "url": "https://google.com", 
  "alias": "google-test"
}
```
*`alias` is optional — a random alphabetic alias is generated if omitted.*

### 3. Delete URL (Protected)
```http
DELETE /url/delete/{alias}
Authorization: Basic <base64(user:pass)>
```
Deletes the URL associated with the given alias.

## Deployment (CI/CD)

The project includes an automated deployment workflow `.github/workflows/deploy.yaml` aimed at an Ubuntu-based VPS (like Oracle Cloud Free Tier). 

**Deployment Process:**
1. Triggers on `workflow_dispatch` (Manual via GitHub Actions UI) requiring a specific Tag (e.g., `v1.0.0`).
2. Builds the Go binary targeting `GOOS=linux GOARCH=amd64`.
3. Uses `rsync` over SSH to securely copy the repository and binary to the target VM.
4. Generates a local `.env` file on the server holding production secrets (`HTTP_SERVER_PASSWORD`).
5. Restarts the Go server using a `systemd` daemon (`url-shortener.service`).

**GitHub Secrets Required for CI/CD:**
- `DEPLOY_SSH_KEY`: The private SSH Key for connecting to the VPS.
- `AUTH_PASS`: The production Basic Auth password.

**Server Networking:**
Ensure port `8082` is open in both your Cloud Provider's ingress rules (VCN) and the VM's internal firewall (e.g., `sudo iptables -I INPUT 1 -p tcp --dport 8082 -j ACCEPT`).

## Tech Stack

- **[Chi](https://github.com/go-chi/chi)** — HTTP router
- **[SQLite](https://github.com/mattn/go-sqlite3)** — DB Storage
- **[cleanenv](https://github.com/ilyakaznacheev/cleanenv)** — Configuration Parsing
- **[godotenv](https://github.com/joho/godotenv)** — Local `.env` integrations
- **[validator](https://github.com/go-playground/validator)** — Struct field validation
