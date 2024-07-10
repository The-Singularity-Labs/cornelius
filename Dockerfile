# Stage 1: Build dependencies (cache Go mod)
FROM golang:bookworm AS compiler
WORKDIR /cornelius

# Copy go.mod to enable caching
COPY go.mod ./
COPY go.sum ./

# Install Go dependencies (cached based on go.mod)
RUN go mod download

COPY ./ ./

RUN go build -o cornelius .

# Stage 2: Build application
FROM node:22-bookworm-slim AS installer
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        ca-certificates \
        git
RUN npm install -g ardrive-cli@2.0.4

# Stage 3: Final slim image
FROM installer  AS cornelius
WORKDIR /scratch

# Copy binary and certificates
COPY --from=compiler /cornelius/cornelius /cornelius
ENTRYPOINT ["/cornelius"]