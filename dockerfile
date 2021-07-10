FROM golang:1.16.5-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh 

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]