FROM golang:1.16 AS builder

WORKDIR /opt
COPY main.go .
COPY go.mod .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app main.go

FROM alpine:latest
WORKDIR /root
COPY --from=builder /opt/app .
CMD [ "/root/app" ]
