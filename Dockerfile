FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY . .

WORKDIR /app/cmd/api
RUN go mod download
RUN go build -o /api

FROM alpine:latest

RUN apk --no-cache add ca-certificates py-pip openssh
RUN pip install --upgrade pip && pip install ansible

COPY --from=builder /api /api

# Запускаем сервисы
CMD ["/api"]
