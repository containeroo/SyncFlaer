FROM golang:1.18-alpine as builder

RUN mkdir -p /go/src/github.com/containeroo/syncflaer
WORKDIR /go/src/github.com/containeroo/syncflaer

RUN apk add --no-cache git

ADD . /go/src/github.com/containeroo/syncflaer
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -installsuffix nocgo -o /syncflaer github.com/containeroo/syncflaer/cmd/syncflaer


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /syncflaer ./
ENTRYPOINT ["./syncflaer"]
