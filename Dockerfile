FROM golang:1.20.3

# Set environment variables
ENV PATH="/root/.cargo/bin:${PATH}"
ENV USER=root

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y \ 
    && rustup default stable \ 
    && rustup target add x86_64-unknown-linux-gnu

WORKDIR /go/src
RUN ln -sf /bin/bash /bin/sh
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

RUN cd /go/pkg/mod/github.com/j178/tiktoken-go@v0.2.1/tiktoken-cffi && \
    cargo build --release

RUN apt update && apt install -y protobuf-compiler && apt install -y protoc-gen-go
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# # Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o /chatserver

# # Optional:
# # To bind to a TCP port, runtime parameters must be supplied to the docker command.
# # But we can document in the Dockerfile what ports
# # the application is going to listen on by default.
# # https://docs.docker.com/engine/reference/builder/#expose
# EXPOSE 8080

# # Run
# CMD ["/chatserver"]

CMD [ "tail", "-f", "/dev/null" ]