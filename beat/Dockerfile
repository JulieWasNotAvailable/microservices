FROM golang:1.23.3-alpine

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

CMD ["go", "run", "main.go"]
