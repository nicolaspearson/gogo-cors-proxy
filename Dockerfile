FROM golang:latest AS builder

RUN mkdir -p /go/src/proxy
COPY . /go/src/proxy/
RUN CGO_ENABLED=0 GOOS=linux go build /go/src/proxy/proxy.go

FROM scratch
COPY --from=builder /go/proxy /usr/local/bin/
