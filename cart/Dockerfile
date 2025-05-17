FROM golang:1.23.3-alpine

COPY .aws/config /root/.aws/config
COPY .aws/credentials /root/.aws/credentials

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

CMD ["go", "run", "main.go"]
