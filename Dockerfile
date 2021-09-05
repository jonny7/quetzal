FROM golang:1.17-alpine AS builder
ENV CGO_ENABLED=0

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/gitlab.com/jonny7/quetzal

COPY . .

RUN go mod download
RUN go build -o main ./app

EXPOSE 8010

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/gitlab.com/jonny7/quetzal/main ./main

CMD ["./main"]
