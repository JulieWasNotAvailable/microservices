FROM golang:1.23.3-alpine

COPY config /root/.aws/config
COPY credentials /root/.aws/credentials

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

CMD ["go", "run", "main.go"]
