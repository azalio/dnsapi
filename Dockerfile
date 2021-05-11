FROM golang:1.16.0-alpine3.13 AS builder

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/github.com/azalio/dnsapi
COPY . .

# RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w' -mod vendor -o /go/bin/dnsapi

FROM scratch
# FROM golang:1.16.0-alpine3.13
COPY --from=builder /go/bin/dnsapi /dnsapi
ENTRYPOINT ["/dnsapi"]
