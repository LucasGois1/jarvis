FROM golang:latest as builder

# Set environment variables
# ENV PATH="/root/.cargo/bin:${PATH}"
# ENV USER=root

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux

# Insall protoc packages
RUN apt update && apt install -y protobuf-compiler && apt install -y protoc-gen-go
RUN apt-get install -y gcc-aarch64-linux-gnu
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Install Rust
# RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y \ 
#     && rustup default stable \ 
#     && rustup target add x86_64-unknown-linux-gnu

# Install dependencies
WORKDIR /go/src
RUN ln -sf /bin/bash /bin/sh
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Build tiktoken-cffi
# RUN cd /go/pkg/mod/github.com/j178/tiktoken-go@v0.2.1/tiktoken-cffi && \
#     cargo build --release

# Build
# tiktoken needs CGO_ENABLED=1
RUN go build -ldflags="-w -s -extldflags '-static'" -o chatserver

FROM alpine

RUN apk update \
    && apk upgrade \
    && add gcc build-base \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true

WORKDIR /go/src
# COPY --from=builder /etc/passwd /etc/passwd
# COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/src/chatserver .

# USER root:root

# RUN chmod +x chatserver

EXPOSE 8080
EXPOSE 50051

ENTRYPOINT ["/chatserver"]
# CMD ["tail", "-f", "/dev/null"]