FROM golang:latest as builder

# Install dependencies
WORKDIR /go/src
RUN ln -sf /bin/bash /bin/sh
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o main

FROM scratch

WORKDIR /go/src
COPY --from=builder /go/src/main .

EXPOSE 8080
EXPOSE 50051

CMD ["./main"]