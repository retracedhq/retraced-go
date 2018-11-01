FROM golang:1.11

RUN go get github.com/satori/go.uuid
RUN go get github.com/stretchr/testify/assert

WORKDIR /go/src/github.com/retracedhq/retraced-go
