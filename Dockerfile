# Stage 1
FROM golang:1.21 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o medods cmd/medods/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/medods /app/medods

EXPOSE 8080

CMD ["./medods"]