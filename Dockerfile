ARG GO_VERSION=1.12.4
FROM golang:${GO_VERSION} AS builder

ENV GO111MODULE=on

WORKDIR /app

# cache go modules
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# build the daemon
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build github.com/xmattstrongx/supermarket

# copy the prebuilt daemon to tiny scratch image
FROM scratch
COPY --from=builder /app/supermarket /bin/
CMD ["/bin/supermarket","daemon"]