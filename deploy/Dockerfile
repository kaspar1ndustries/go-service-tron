FROM golang:1.19.3 as builder
WORKDIR /go/src/go-tron-tiny-wallet/
COPY ./src ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/go-tron-tiny-wallet/app .
CMD ["./app"]