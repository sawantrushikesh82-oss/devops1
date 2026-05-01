# Secure Managed File Transfer (MFT) Platform

Production-oriented monorepo implementing:
- API service (`cmd/api`)
- External SFTP service (`cmd/ext-sftp`)
- Internal SFTP service (`cmd/int-sftp`)
- React operations UI (`web/`)

## Current status

This repository is scaffolded and wired end-to-end, but **full build/test in this execution environment is blocked** by outbound module download restrictions (`proxy.golang.org` returns `403 Forbidden`).

You can still run successfully in a normal developer machine/network.

## Run locally with Docker Compose

### 1) Generate keys

Linux/macOS:

```bash
mkdir -p configs
openssl genpkey -algorithm RSA -out configs/jwt_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in configs/jwt_private.pem -out configs/jwt_public.pem
ssh-keygen -t rsa -b 4096 -f configs/sftp_host_key -N ""
```

Windows PowerShell:

```powershell
# if OpenSSL is unavailable, install Git for Windows/OpenSSL first, then run
mkdir configs -ErrorAction SilentlyContinue
openssl genpkey -algorithm RSA -out configs/jwt_private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in configs/jwt_private.pem -out configs/jwt_public.pem
ssh-keygen -t rsa -b 4096 -f configs/sftp_host_key -N ""
```

### 2) Start services

```bash
docker compose up --build
```

### 3) Access endpoints

- API: http://localhost:8080
- API health: http://localhost:8080/health
- Web: http://localhost:5173
- Prometheus: http://localhost:9090
- External SFTP: `localhost:2222`
- Internal SFTP: `localhost:2223`

## Run services directly (without Docker)

Start PostgreSQL first, then:

```bash
go run ./cmd/ext-sftp
go run ./cmd/int-sftp
go run ./cmd/api
```

Frontend:

```bash
cd web
npm install
npm run dev
```

## Known follow-up work

1. Add unit/integration tests for SFTP VFS permission enforcement.
2. Add API integration tests for JWT login/refresh and role-based route protection.
3. Add sync worker e2e tests using temp dirs and fixture files.
4. Add CI pipeline stages for lint/test/build with caching.
