FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add git tzdata make \
    && rm -rf /var/cache/apk/*

WORKDIR /go/mimic

# Download go modules
COPY go.mod .
COPY go.sum .
RUN GO111MODULE=on GOPROXY=https://proxy.golang.org go mod download

COPY . .

RUN make build

# Create a minimal container to run a Golang static binary
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/mimic/mimic .

ENTRYPOINT ["/mimic"]
EXPOSE 80