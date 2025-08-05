# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Run stage
FROM alpine:latest

# Install bash (untuk wait-for-it.sh) dan tzdata (timezone database)
RUN apk add --no-cache bash tzdata

WORKDIR /app

COPY --from=builder /app/main .
COPY wait-for-it.sh .

RUN chmod +x wait-for-it.sh

EXPOSE 8080

CMD ["./wait-for-it.sh", "db:5432", "--", "./main"]
