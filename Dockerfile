# Подготовка.
FROM golang:alpine AS preper
ARG NAME
WORKDIR /usr/src/${NAME}
COPY . .
RUN mkdir /usr/local/share/doc
COPY ./swagger.yaml /usr/local/share/doc/swagger.yaml
COPY ./config /etc
RUN go mod download

# Тестирование.
FROM preper AS testing
ARG DB_PASSWORD
RUN set -o pipefail
RUN go test --tags=unit -v ./...
RUN go test --tags=integration -v ./...
RUN go test --tags=e2e -v ./...

# Сборка.
FROM preper AS building
RUN go build -v -o /usr/local/bin/${NAME} ./cmd/main.go

# Релиз.
FROM alpine AS release
EXPOSE 8000
ARG NAME
# Настройки.
RUN apk add tzdata
RUN ln -s /usr/share/zoneinfo/Europe/Moscow /etc/localtime
# Копирование конфигурации.
COPY ./config /etc
# Копирование исполняемого файла.
COPY --from=building /usr/local/bin/${NAME} /usr/local/bin/${NAME}
# Копирование документации.
RUN mkdir /usr/local/share/doc
COPY ./swagger.yaml /usr/local/share/doc/swagger.yaml
# Копирование точки входа.
COPY ./docker-entrypoint.sh /usr/local/bin
RUN echo "exec ${NAME}" >> /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENTRYPOINT [ "docker-entrypoint.sh" ]
