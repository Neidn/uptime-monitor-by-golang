FROM golang:1.20.6-alpine3.8

WORKDIR /app

COPY . .

RUN go build -o build ./...

CMD ["./build/monitor"]
