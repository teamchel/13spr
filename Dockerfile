FROM golang:1.25-alpine

WORKDIR /app

# Устанавливаем git
RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]