# apps/backend/Dockerfile

# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# 1. Копіюємо описи залежностей
COPY go.mod go.sum ./

# 2. Завантажуємо їх (це закешується докером, якщо go.mod не змінився)
RUN go mod download

# 3. Копіюємо код і збираємо
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: Run (без змін)
FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/server .
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]