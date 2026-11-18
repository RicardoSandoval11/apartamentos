FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/apartments-service ./main.go

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /
COPY --from=builder /app/apartments-service /apartments-service

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT [ "/apartments-service" ]