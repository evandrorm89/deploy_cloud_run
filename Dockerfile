FROM golang:1.21 as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

FROM alpine:latest as certs
RUN apk --no-cache add ca-certificates

FROM scratch
WORKDIR /app
COPY --from=build /app/cloudrun .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["./cloudrun"]
