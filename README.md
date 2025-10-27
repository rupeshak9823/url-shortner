# URL Shortener Service

A backend service that allows users to shorten URLs, redirect to original links, and view basic analytics.

---

## Features

- Shorten long URLs with Base62 encoding  
- Idempotent shortening (same long URL → same short code)  
- Redirect to original URL  
- Persistent storage (PostgreSQL)  
- Caching with Redis  
- Basic rate limiting  
- Dockerized setup  

---

## Tech Stack

- Go 1.24  
- PostgreSQL 15  
- Redis 7  
- GORM v2  
- Goose (migrations)  
- Docker + Docker Compose  

---

## Architecture Overview
Client → Handler → Service → Repository → (Postgres + Redis)
- Handlers expose HTTP APIs  
- Services contain business logic  
- Repository handles data access and caching  
- Redis used for fast lookups  
- Postgres stores persistent mappings and analytics

---

## Setup

### 1. Prerequisites
- Docker and Docker Compose installed  
- Port `8080` available  

### 2. Run the Service
docker-compose up --build

---

## API Documentation

### 1. Shorten URL
- POST /url
  ```
  curl --location 'localhost:8080/url' \
  --header 'Content-Type: application/json' \
  --data '{
      "longURL": "https://chatgpt.com/c/68fe34e0-edb4-8324-b359-b5232bd95adf"
  }'
  ```
- Response
  ```
  {
    "data": {
        "id": "0ce3c8bf-ae5e-4f67-b499-7cb0d94acbb7",
        "shortCode": "WbGbHefo",
        "longURL": "https://chatgpt.com/c/68fe34e0-edb4-8324-b359-b5232bd95adf",
        "createdAt": "2025-10-26T19:20:44.85756Z",
        "clickedCount": 19,
        "lastAccessedAt": "2025-10-27T12:36:47.78859Z"
    }
  }
  ```
### 2. Redirect
- GET /url/:shortCode
  ```
    curl --location --request GET 'localhost:8080/url/WbGbHefo'
  ```
- Response
  ```
    307 Temporary Redirect → redirects to original URL.
  ```

### 3. Admin List URLs
- GET /admin/urls?page=1&perPage=2
  ```
    curl --location 'localhost:8080/admin/urls?page=1&perPage=2'
  ```
- Response
  ```
    {
      "data": [
          {
              "id": "0ce3c8bf-ae5e-4f67-b499-7cb0d94acbb7",
              "shortCode": "WbGbHefo",
              "longURL": "https://chatgpt.com/c/68fe34e0-edb4-8324-b359-b5232bd95adf",
              "createdAt": "2025-10-26T19:20:44.85756Z",
              "clickedCount": 19,
              "lastAccessedAt": "2025-10-27T12:36:47.78859Z"
          },
          {
              "id": "aacf8604-d3fc-40b6-8129-7ec677f836a8",
              "shortCode": "dBJ0AIPM",
              "longURL": "https://chatgpt.com/c/68fe34e0-edb4-8324-b359-b5232bd95ade",
              "createdAt": "2025-10-26T16:11:38.450912Z",
              "clickedCount": 15,
              "lastAccessedAt": "2025-10-26T19:20:23.969096Z"
          }
      ],
      "total": 2
    }
  ```
---

## Design Decisions and Trade-offs

| Aspect | Decision | Reasoning / Trade-off |
|--------|-----------|------------------------|
| **URL Encoding** | Used Base62 (0–9, a–z, A–Z) | Generates short, readable, and URL-safe codes without special characters. |
| **Caching** | Redis used for hot lookups | Greatly reduces DB load and improves redirect speed. Trade-off: potential cache inconsistency during updates. |
| **Database** | PostgreSQL | Reliable relational database with strong indexing and JSON support. Trade-off: slightly heavier for small-scale deployments. |
| **ORM** | GORM | Simplifies model management and migrations. Trade-off: some loss of fine-grained SQL control. |
| **Rate Limiting** | In-memory IP-based limiter (token bucket) | Prevents abuse. Trade-off: rate limit resets on instance restart (can be improved with Redis). |
| **Architecture** | Layered: Handler → Service → Repository | Encourages clean separation of concerns and testability. |
| **Error Handling** | Centralized HTTP error formatter | Ensures consistent API responses across endpoints. |

---




