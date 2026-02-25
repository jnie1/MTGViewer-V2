#Building stage
FROM golang:1.23-bookworm AS builder

WORKDIR /mtgviewer-v2

#Get dependencies
COPY go.mod go.sum ./
RUN go mod download

#Get source code
COPY . .

#Compile code
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o mtgviewer-v2-backend .

#Runtime stage
FROM debian:bookworm-slim

#Grab certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /mtgviewer-v2

#Bring the binaries over
COPY --from=builder /mtgviewer-v2/mtgviewer-v2-backend .

#Port this application will be listening
EXPOSE 8080

#Run executable when container starts
ENTRYPOINT [ "./mtgviewer-v2-backend" ]