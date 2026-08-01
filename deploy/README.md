# Deployment Guide

This guide covers building images, running services locally with
docker-compose, and deploying to a VPS with rootless Podman Quadlets.

## Local development with docker-compose

```sh
docker compose up --build
```

This starts three services on loopback:

- PostgreSQL on `127.0.0.1:5432`
- API on `127.0.0.1:8081`
- Web (Caddy) on `127.0.0.1:8082`

The API health check is at `http://127.0.0.1:8081/healthz`.
The web client is at `http://127.0.0.1:8082`.

## Building images for production

Build from the repo root:

```sh
podman build -t dos-freightflow-api:latest -f deploy/images/Dockerfile.api .
podman build -t dos-freightflow-web:latest -f deploy/images/Dockerfile.web .
```

The API image is a static Go binary on Alpine, running as a non-root user.
The web image serves the built React bundle via Caddy.

## VPS deployment with Podman Quadlets

### 1. Prepare the host

```sh
sudo useradd -m -s /bin/bash dos
sudo loginctl enable-linger dos
sudo apt-get install -y podman uidmap slirp4netns fuse-overlayfs
```

`enable-linger` keeps services running after logout — without it, all
containers stop when the user session ends.

### 2. Create secret files

As the `dos` user:

```sh
mkdir -p ~/.config/dos-freightflow
cat > ~/.config/dos-freightflow/api.env <<'EOF'
DATABASE_URL=postgres://dos:SECRET@127.0.0.1:5433/dosfreightflow?sslmode=require
SESSION_SECRET=GENERATE_WITH_openssl_rand_hex_32
OTP_SECRET=GENERATE_WITH_openssl_rand_hex_32
EOF

cat > ~/.config/dos-freightflow/postgres.env <<'EOF'
POSTGRES_PASSWORD=GENERATE_A_STRONG_PASSWORD
EOF

chmod 600 ~/.config/dos-freightflow/api.env
chmod 600 ~/.config/dos-freightflow/postgres.env
```

### 3. Install Quadlet units

```sh
mkdir -p ~/.config/containers/systemd
cp deploy/quadlet/*.container ~/.config/containers/systemd/
systemctl --user daemon-reload
systemctl --user enable --now dos-freightflow-postgres.service
systemctl --user enable --now dos-freightflow-api.service
systemctl --user enable --now dos-freightflow-web.service
```

### 4. Run migrations

```sh
podman exec dos-freightflow-api /app/migrate up
```

### 5. Verify

```sh
curl http://127.0.0.1:8081/healthz
curl http://127.0.0.1:8082/
```

### 6. Pin image digests

After all checks pass, pin architecture-specific digests:

```sh
podman pull dos-freightflow-api:latest
DIGEST=$(podman inspect --format '{{.Digest}}' dos-freightflow-api:latest)
# Update the Quadlet Image= line to dos-freightflow-api@sha256:$DIGEST
systemctl --user daemon-reload
systemctl --user restart dos-freightflow-api.service
```

## Rollback

To roll back to a previous version:

```sh
podman pull dos-freightflow-api:PREVIOUS_TAG
# Update the Quadlet Image= line
systemctl --user daemon-reload
systemctl --user restart dos-freightflow-api.service
```

Database migrations roll back with:

```sh
podman exec dos-freightflow-api /app/migrate down
```

## Security checklist

- All containers run as non-root users
- `DropCapability=ALL`, `ReadOnly=true`, `Tmpfs=/tmp` on all units
- Secrets in `chmod 600` env files, never committed
- Caddy sets HSTS, nosniff, DENY frame, no-referrer
- API binds to loopback only — no direct public access
- Health checks on all services
- Image digests pinned after acceptance checks pass