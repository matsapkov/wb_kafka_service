# build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/server

# run stage
FROM alpine:3.20
WORKDIR /app

RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder /bin/app /app/app

EXPOSE 8081
ENTRYPOINT ["/app/app"]