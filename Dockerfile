FROM golang:1.15-alpine as builder

RUN mkdir -p /go/src/syncflaer
WORKDIR /go/src/syncflaer

RUN apk add --no-cache git

ADD . /go/src/syncflaer
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -installsuffix nocgo -o /syncflaer syncflaer/cmd


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /syncflaer ./
ENTRYPOINT ["./syncflaer"]
