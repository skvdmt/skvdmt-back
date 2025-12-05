FROM golang:1.25.3-alpine AS builder
WORKDIR /usr/src/skvdmt-back
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /usr/local/bin/skvdmt-back ./cmd/main.go

FROM builder AS testing
RUN go test -v ./...

FROM alpine:3.22.2
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
WORKDIR /usr/local/bin
COPY --from=builder /usr/local/bin/skvdmt-back ./skvdmt-back
RUN mkdir -p /var/log/skvdmt-back
RUN mkdir -p /etc/skvdmt-back
COPY ./config/* /etc/skvdmt-back
EXPOSE 8000
CMD ["skvdmt-back"]
