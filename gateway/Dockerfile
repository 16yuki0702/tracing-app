# builder
FROM golang:latest as builder
MAINTAINER 16yuki0702

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway .

# application
FROM alpine:latest

WORKDIR /root/
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]
