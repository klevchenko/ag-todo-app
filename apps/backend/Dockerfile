# Stage 1: Build (Збираємо бінарник)
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
# COPY go.sum ./ # поки немає залежностей, це не треба
COPY main.go .
# CGO_ENABLED=0 критично для distroless
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: Run (Запускаємо)
FROM gcr.io/distroless/static-debian12
WORKDIR /
COPY --from=builder /app/server .
EXPOSE 8080
# Запускаємо від імені користувача nonroot (ID 65532) для безпеки
USER nonroot:nonroot
ENTRYPOINT ["/server"]