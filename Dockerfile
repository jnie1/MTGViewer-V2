# Frontend build stage (Vue)
FROM node:20-bullseye AS frontend-builder

WORKDIR /mtgviewer-v2

# Copy the repo so `file:` dependencies and monorepo references resolve during npm install
COPY . .

# Install frontend dependencies and build
RUN cd frontend && npm ci && npm run build

# Backend build stage (Go)
FROM golang:1.25-bookworm AS backend-builder

WORKDIR /mtgviewer-v2/backend

# Get Go module dependencies
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

# Copy the backend source into the module directory
COPY ./backend/. ./

# Bring built frontend into the backend build context so the backend can serve ./dist
COPY --from=frontend-builder /mtgviewer-v2/frontend/dist /mtgviewer-v2/dist

# Compile Go binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /mtgviewer-v2/mtgviewer-v2-backend .

# Runtime stage
FROM debian:bookworm-slim

# Grab certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

WORKDIR /mtgviewer-v2

# Copy the backend binary and frontend dist into the runtime image
COPY --from=backend-builder /mtgviewer-v2/mtgviewer-v2-backend .
COPY --from=backend-builder /mtgviewer-v2/dist ./dist
COPY backend/.env .

# Port this application will be listening
EXPOSE 8080

# Run executable when container starts
ENTRYPOINT [ "./mtgviewer-v2-backend" ]