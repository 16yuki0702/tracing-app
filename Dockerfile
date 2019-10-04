# builder
FROM golang:latest as builder
MAINTAINER 16yuki0702

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# application
FROM alpine:latest

WORKDIR /root/
RUN apk --no-cache add ca-certificates curl
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]
