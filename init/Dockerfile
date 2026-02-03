FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o icstz-server ./cmd/icstz-server

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/icstz-server .
EXPOSE 24501
CMD ["./icstz-server", "24501"]
