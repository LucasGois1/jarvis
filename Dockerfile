
FROM golang:1.20-alpine as builder

RUN apk add bash ca-certificates git gcc g++ libc-dev

# Install dependencies
WORKDIR /go/src
RUN ln -sf /bin/bash /bin/sh
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main

FROM alpine

RUN apk add --no-cache --upgrade ca-certificates openssl

WORKDIR /go/src
# Get certificates
COPY --from=builder /go/src/main .

EXPOSE 8080
EXPOSE 50051

ENTRYPOINT ["./main"]