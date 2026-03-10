# URL Shortener

A simple URL shortener service built with Go, Chi router, and SQLite.

## Project Structure

```
cmd/url-shortener/       — application entry point
config/                  — YAML configuration files
internal/
  config/                — config loader
  storage/sqlite/        — SQLite storage layer
  http-server/
    handlers/url/
      save/              — POST /save
      redirect/          — GET /get/{alias}
      delete/            — DELETE /delete/{alias}
    middleware/logger/    — custom request logger
  lib/
    api/response/        — shared JSON response helpers
    random/              — random alias generator
    logger/              — logger utilities
```

## Setup

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
   ```
   CONFIG_PATH=./config/local.yaml
   ```

5. **Run the server**
   ```bash
   go run cmd/url-shortener/main.go
   ```

## Configuration

Edit `config/local.yaml`:

```yaml
env: "local"
storage_path: "./storage/storage.db"
http_server:
  address: "localhost:8082"
  timeout: 4s
  idle_timeout: 30s
```

## API

### Save URL
```
POST /save
Content-Type: application/json

{"url": "https://google.com", "alias": "google"}
```
`alias` is optional — a random one is generated if omitted.

### Redirect
```
GET /get/{alias}
```
Redirects (302) to the original URL.

### Delete URL
```
DELETE /delete/{alias}
```
Deletes the URL associated with the given alias.

## Tech Stack

- [Chi](https://github.com/go-chi/chi) — HTTP router
- [SQLite](https://github.com/mattn/go-sqlite3) — storage
- [cleanenv](https://github.com/ilyakaznacheev/cleanenv) — config parsing
- [godotenv](https://github.com/joho/godotenv) — .env file loading
- [validator](https://github.com/go-playground/validator) — request validation