FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/MOOVE-Network/location_service
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o location_service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/MOOVE-Network/location_service/location_service .
CMD ["./location_service"]
