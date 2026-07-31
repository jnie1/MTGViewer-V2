# Frontend build stage (Vue)
FROM node:20-bookworm AS frontend-builder

WORKDIR /mtgviewer-v2/frontend

ARG VITE_API_URL=" "

ENV VITE_API_URL=$VITE_API_URL

# Copy frontend package files first to leverage build cache
COPY frontend/package*.json ./
RUN npm ci

# Copy the frontend source and build
COPY frontend ./
RUN npm run build

# Backend build stage (Go)
FROM golang:1.25-bookworm AS backend-builder

WORKDIR /mtgviewer-v2/backend

# Get Go module dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source code
COPY backend ./

# Compile Go binary
RUN CGO_ENABLED=0 go build -tags=prod -ldflags="-s -w" -o /mtgviewer-v2/mtgviewer-v2-backend .

FROM scratch
# Runtime stage
FROM debian:bookworm-slim

# Grab certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

WORKDIR /mtgviewer-v2

# Copy the backend binary and frontend dist into the runtime image
COPY --from=backend-builder /mtgviewer-v2/mtgviewer-v2-backend .
COPY --from=frontend-builder /mtgviewer-v2/frontend/dist ./dist

# Port this application will be listening
EXPOSE 8080

# Run executable when container starts
ENTRYPOINT [ "./mtgviewer-v2-backend" ]