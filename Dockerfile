FROM golang:1.24-alpine AS build

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/migrate ./cmd/migrate

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /out/api /app/api
COPY --from=build /out/migrate /app/migrate
COPY migrations /app/migrations

ENV SERVICE_NAME=hirify-go-test
ENV HTTP_ADDR=:8080
ENV MIGRATIONS_DIR=/app/migrations

EXPOSE 8080

CMD ["/app/api"]

